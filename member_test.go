package gateway

import (
	"context"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

const (
	memberKey    = "d7ff421c8912c028ab1fa854ae5d11ba-73a563c6c3884b1fb88bf0093dbd04a3"
	memberSecret = "0e61ec3fea7c2753d6a2c1f84d07621e"
)

func TestMemberAuth(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	ctx := context.Background()
	c := NewClient(apiBase).Member(memberKey, memberSecret)
	m, err := c.MemberInfo(ctx)
	if assert.Nil(t, err) {
		assert.NotEmpty(t, m.ID)
	}
}
