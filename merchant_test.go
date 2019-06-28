package gateway

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

const (
	merchantKey    = "5c8a9491dca25af694004d5e1711b217"
	merchantSecret = "64012120f9fb7daaa9f6ae48a159584d"
	apiBase        = "https://dev-gateway.fox.one"
)

func TestMerchantCreateMember(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	c := NewMerchantClient(apiBase).WithSession(merchantKey, merchantSecret)

	ctx := context.Background()
	output, err := c.CreateMember(ctx)
	if !assert.Nil(t, err) {
		return
	}

	s := output.Session
	m := output.Member

	assert.NotEmpty(t, s.Key)
	assert.NotEmpty(t, s.Secret)
	assert.True(t, s.ExpiredAt > s.CreatedAt)

	memberID := m.ID
	assert.False(t, m.IsPinSet)

	m, s, err = c.LoginMember(ctx, memberID, time.Hour)
	if assert.Nil(t, err) {
		assert.NotEmpty(t, s.Key)
		assert.NotEmpty(t, s.Secret)
		assert.True(t, s.ExpiredAt > s.CreatedAt)
	}
}

func TestMerchantService(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	c := NewMerchantClient(apiBase).WithSession(merchantKey, merchantSecret)

	ctx := context.Background()
	memberID := "73a563c6c3884b1fb88bf0093dbd04a3"
	wallets, err := c.MemberWallets(ctx, memberID, "")
	if assert.Nil(t, err) && assert.NotEmpty(t, wallets) {
		for _, w := range wallets {
			assert.NotEmpty(t, w.WalletID)
			assert.Equal(t, memberID, w.MemberID)
		}
	}

}

func TestClearMemberSession(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	ctx := context.Background()
	c := NewMerchantClient(apiBase).WithSession(merchantKey, merchantSecret)
	err := c.ClearUserSessions(ctx, "73a563c6c3884b1fb88bf0093dbd04a3")
	assert.Nil(t, err)
}

func TestJsonMarshalBytes(t *testing.T) {
	var form struct {
		Body []byte `json:"body"`
	}

	form.Body = []byte("{\"name\":\"yiplee\"}")
	data, _ := json.Marshal(form)
	assert.Empty(t, string(data))

	form.Body = nil
	json.Unmarshal(data, &form)
	assert.Empty(t, string(form.Body))
}
