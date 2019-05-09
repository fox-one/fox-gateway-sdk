package gateway

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetPublicKey(t *testing.T) {
	key, err := getPublicKey()
	assert.Nil(t, err)
	assert.NotNil(t, key)
}

func TestRsaEncrypt(t *testing.T) {
	key, err := getPublicKey()
	assert.Nil(t, err)
	assert.NotNil(t, key)

	v := map[string]interface{}{
		"t": time.Now().Unix(),
	}

	data, _ := json.Marshal(v)
	r, err := rsaEncrypt(data)
	assert.Nil(t, err)
	assert.NotEmpty(t, r)
}
