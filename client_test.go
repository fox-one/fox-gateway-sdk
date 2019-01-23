package gateway

import (
	"net/url"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJoinPath(t *testing.T) {
	u := url.URL{}
	u.Path = "member/info?hahah=v"
	u.Scheme = "https"
	u.Host = "dev.fox.one"
	s := u.String() // prints http://foo/bar.html
	assert.Empty(t, s)
}

func TestPathJoin(t *testing.T) {
	a := ""
	b := "/member/info?key=value"
	assert.Empty(t, path.Join(a, b))
}

func TestParseURL(t *testing.T) {
	uri := "api.fox.one"
	u, _ := url.Parse(uri)
	assert.Empty(t, u.Scheme)
	assert.Empty(t, u.Host)
}
