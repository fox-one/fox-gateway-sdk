package gateway

import (
	"context"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGetAssets(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	merchantClient := NewMerchantClient(apiBase).WithSession(merchantKey, merchantSecret)
	memberSvc := merchantClient.MemberService("payment", "e0814259f9c34d58b010eb674049d883")

	ctx := context.Background()
	assets, err := memberSvc.ReadAssets(ctx, 1)
	assert.Nil(t, err)
	assert.NotEmpty(t, assets)
}
