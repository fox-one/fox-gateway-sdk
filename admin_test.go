package gateway

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAdmin(t *testing.T) {
	key := "59ded60d7a3698a5ea9e611eb0b07e48-zhangyh"
	secret := "72eca0eb4bdb0629d1e8943e7a8db036"
	c := NewAdminClient(apiBase).WithSession(key, secret)
	admin, session, err := c.CreateAdmin(context.Background(), "yiplee", "123456")
	if assert.Nil(t, err) {
		assert.Empty(t, admin.Username)
		assert.Empty(t, session.Key)
	}
}
