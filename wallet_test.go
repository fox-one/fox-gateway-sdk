package gateway

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

const (
	merchantKey    = "5c8a9491dca25af694004d5e1711b217"
	merchantSecret = "64012120f9fb7daaa9f6ae48a159584d"
	apiBase        = "https://dev-gateway.fox.one"
	memberID       = "111"
)

func TestGetAssets(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	merchantClient := NewMerchantClient(apiBase).WithSession(merchantKey, merchantSecret)
	memberSvc := merchantClient.MemberService("payment", memberID)

	assets, err := memberSvc.GetAssets(ctx)

	assert.NotEmpty(t, assets.Key)
}
