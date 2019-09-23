package gateway

import (
	"context"
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/fox-one/httpclient"
	jsoniter "github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"
)

// MerchantClient merchant client
type MerchantClient struct {
	*Client
	key    string
	secret string
}

// Merchant merchant
func (c *Client) Merchant() *MerchantClient {
	return &MerchantClient{
		Client: c.Group("merchant"),
	}
}

// NewMerchantClient new client
func NewMerchantClient(apiBase string) *MerchantClient {
	return NewClient(apiBase).Merchant()
}

// WithSession session
func (m *MerchantClient) WithSession(key, secret string) *MerchantClient {
	return &MerchantClient{
		Client: m.Client,
		key:    key,
		secret: secret,
	}
}

// Member member
func (m *MerchantClient) Member(id string) *MemberClient {
	return m.Client.Member().WithAuth(&merchantMemberAuth{
		merchantKey:    m.key,
		merchantSecret: m.secret,
		memberID:       id,
	})
}

// auth
type merchantAuth struct {
	key       string
	secret    string
	memberID  string
	memberPin string
	expire    time.Duration
}

func (m *merchantAuth) token(method, uri string, body []byte) string {
	claims := jwt.MapClaims{
		"exp":  time.Now().Add(m.expire).Unix(),
		"sign": signRequest(method, uri, string(body)),
		"key":  m.key,
	}

	if m.memberID != "" {
		claims["mem"] = m.memberID

		if m.memberPin != "" {
			payload := map[string]interface{}{
				"t": time.Now().Unix(),
				"n": newNonce(),
				"p": m.memberPin,
			}

			data, _ := jsoniter.Marshal(payload)
			pinToken, err := rsaEncrypt(data)

			if err != nil {
				log.Panic(err)
			}

			claims["pin"] = pinToken
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(m.secret))
	if err != nil {
		log.Panicln("sign merchent token", err)
	}

	return t
}

func (m *merchantAuth) Auth(req *httpclient.Request, method, uri string, body []byte) {
	req.AddToken(m.token(method, uri, body))
}

// Presign pre-sign
func (m *MerchantClient) Presign(expire time.Duration) *merchantAuth {
	return &merchantAuth{
		key:    m.key,
		secret: m.secret,
		expire: expire,
	}
}

// member

// CreateMemberOutput create member
type CreateMemberOutput struct {
	Member  *MemberView         `json:"member,omitempty"`
	Session *MemberSessionView  `json:"session,omitempty"`
	Wallets []*MemberWalletView `json:"wallets,omitempty"`
}

// CreateMember create member
func (m *MerchantClient) CreateMember(ctx context.Context) (*CreateMemberOutput, error) {
	data, err := m.POST("/member/new").Auth(m.Presign(time.Minute)).Do(ctx).Bytes()
	if err != nil {
		return nil, err
	}

	var resp struct {
		Err

		*CreateMemberOutput
	}

	if err := jsoniter.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	if resp.Code > 0 {
		return nil, resp.Err
	}

	return resp.CreateMemberOutput, nil
}

// LoginMember login member
func (m *MerchantClient) LoginMember(ctx context.Context, id string, expire time.Duration) (*MemberView, *MemberSessionView, error) {
	data, err := m.POST("/member/login").
		P("id", id).
		P("expire", int64(expire.Seconds())).
		Auth(m.Presign(time.Minute)).
		Do(ctx).Bytes()
	if err != nil {
		return nil, nil, err
	}

	var resp struct {
		Err

		Member  *MemberView        `json:"member"`
		Session *MemberSessionView `json:"session"`
	}

	if err := jsoniter.Unmarshal(data, &resp); err != nil {
		return nil, nil, err
	}

	if resp.Code > 0 {
		return nil, nil, resp.Err
	}

	return resp.Member, resp.Session, nil
}

// ClearUserSessions clear user session
func (m *MerchantClient) ClearUserSessions(ctx context.Context, memberID string, sessionKey ...string) error {
	req := m.POST("/member/logout").P("id", memberID)
	if len(sessionKey) > 0 {
		req = req.P("session_key", sessionKey[0])
	}
	data, err := req.Auth(m.Presign(time.Minute)).Do(ctx).Bytes()
	if err != nil {
		return err
	}

	var resp struct {
		Err
	}

	if err := jsoniter.Unmarshal(data, &resp); err != nil {
		return err
	}

	if resp.Code > 0 {
		return resp.Err
	}

	return nil
}

// MemberWallets member wallets
func (m *MerchantClient) MemberWallets(ctx context.Context, memberID string, service string) ([]*MemberWalletView, error) {
	result := m.GET("/member/services").
		P("member_id", memberID).
		P("service", service).
		Auth(m.Presign(time.Minute)).
		Do(ctx)

	data, err := result.Bytes()
	if err != nil {
		return nil, err
	}

	var wallets []*MemberWalletView
	if err := jsoniter.Unmarshal(data, &wallets); err == nil {
		return wallets, nil
	}

	var e Err
	if jsoniter.Unmarshal(data, &e) == nil && e.Code > 0 {
		return nil, e
	}

	_, status := result.Status()
	return nil, errors.New(status)
}

// FetchSnapshots fetch snapshots
func (m *MerchantClient) FetchSnapshots(ctx context.Context, service, assetID, cursor, order string, limit int) ([]*WalletSnapshotView, *Pagination, error) {
	result := m.GET("/member/snapshots").
		P("service", service).
		P("asset_id", assetID).
		P("cursor", cursor).
		P("order", order).
		P("limit", limit).
		Auth(m.Presign(time.Minute)).
		Do(ctx)

	data, err := result.Bytes()
	if err != nil {
		return nil, nil, err
	}

	var resp struct {
		Err
		Snapshots  []*WalletSnapshotView `json:"snapshots"`
		Pagination *Pagination           `json:"pagination"`
	}

	if err := jsoniter.Unmarshal(data, &resp); err != nil {
		return nil, nil, err
	}

	if resp.Code > 0 {
		return nil, nil, resp.Err
	}

	return resp.Snapshots, resp.Pagination, nil
}

// ReadExternalSnapshots external snapshots
func (m *MerchantClient) ReadExternalSnapshots(ctx context.Context, from, to int64, limit int) ([]*ExternalSnapshotView, error) {
	result := m.GET("/snapshots/external").
		P("from", from).
		P("to", to).
		P("limit", limit).
		Auth(m.Presign(time.Minute)).
		Do(ctx)

	data, err := result.Bytes()
	if err != nil {
		return nil, err
	}

	var resp struct {
		Err
		Snapshots []*ExternalSnapshotView `json:"snapshots"`
	}

	if err := jsoniter.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	if resp.Code > 0 {
		return nil, resp.Err
	}

	return resp.Snapshots, nil
}
