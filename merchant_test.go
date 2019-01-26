package gateway

import (
	"context"
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
	m, s, err := c.CreateMember(ctx)
	if !assert.Nil(t, err) {
		return
	}

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
	data, err := c.GET("/member/services").P("member_id", "73a563c6c3884b1fb88bf0093dbd04a3").Auth(c.Presign(time.Minute)).Do(ctx).Bytes()
	assert.Nil(t, err)
	assert.NotEmpty(t, string(data))
}

func TestClearMemberSession(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	ctx := context.Background()
	c := NewMerchantClient(apiBase).WithSession(merchantKey, merchantSecret)
	err := c.ClearUserSessions(ctx, "73a563c6c3884b1fb88bf0093dbd04a3")
	assert.Nil(t, err)
}
