package gateway

import (
	"context"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type MemberService struct {
	*Client

	authFunc func(expire time.Duration) Authenticator
}

func (m *Member) Service(name string) *MemberService {
	return &MemberService{
		Client: NewClient(m.apiBase + "/" + name),
		authFunc: func(expire time.Duration) Authenticator {
			return m.Presign(expire)
		},
	}
}

func (m *Member) ServiceWithPin(name, pin string) *MemberService {
	return &MemberService{
		Client: NewClient(m.apiBase + "/" + name),
		authFunc: func(expire time.Duration) Authenticator {
			return m.PresignWithPin(pin, expire)
		},
	}
}

func (m *Merchant) MemberService(name, member string) *MemberService {
	return &MemberService{
		Client: NewClient(m.apiBase + "/" + name),
		authFunc: func(expire time.Duration) Authenticator {
			return m.PresignMember(member, expire)
		},
	}
}

func (m *MemberService) Presign(expire time.Duration) Authenticator {
	return m.authFunc(expire)
}

// wallet

type WalletUserAssetsView []*WalletUserAssetView

func (m *MemberService) ReadAssets(ctx context.Context) (WalletUserAssetsView, error) {
	data, err := m.GET("/assets").Auth(m.Presign(time.Minute)).Do(ctx)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Err
		WalletUserAssetsView
	}

	if err := jsoniter.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	if resp.Code > 0 {
		return nil, resp.Err
	}

	return resp.WalletUserAssetsView, nil
}
