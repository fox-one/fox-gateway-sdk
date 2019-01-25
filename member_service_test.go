package gateway

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	log "github.com/sirupsen/logrus"
)

func TestMemberService(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	c := NewMerchantClient(apiBase).WithSession(merchantKey, merchantSecret)
	s := c.MemberService("payment", "73a563c6c3884b1fb88bf0093dbd04a3")
	data, err := s.GET("/assets").Auth(s.Presign(time.Minute)).Do(context.Background()).Bytes()
	assert.Nil(t, err)
	assert.Empty(t, string(data))
}
