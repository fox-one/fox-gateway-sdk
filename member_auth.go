package gateway

import (
	"time"

	"github.com/fox-one/httpclient"
)

type MemberAuth interface {
	Presign(exp time.Duration) httpclient.Authenticator
	PresignWithPin(pin string, exp time.Duration) httpclient.Authenticator
}

type memberSessionAuth struct {
	key    string
	secret string
}

func (auth *memberSessionAuth) Presign(exp time.Duration) httpclient.Authenticator {
	return &memberAuth{
		key:    auth.key,
		secret: auth.secret,
		expire: exp,
	}
}

func (auth *memberSessionAuth) PresignWithPin(pin string, exp time.Duration) httpclient.Authenticator {
	return &memberAuth{
		key:    auth.key,
		secret: auth.secret,
		pin:    pin,
		expire: exp,
	}
}

type merchantMemberAuth struct {
	merchantKey    string
	merchantSecret string
	memberID       string
}

func (auth *merchantMemberAuth) Presign(exp time.Duration) httpclient.Authenticator {
	return &merchantAuth{
		key:      auth.merchantKey,
		secret:   auth.merchantSecret,
		memberID: auth.memberID,
		expire:   exp,
	}
}

func (auth *merchantMemberAuth) PresignWithPin(pin string, exp time.Duration) httpclient.Authenticator {
	return &merchantAuth{
		key:       auth.merchantKey,
		secret:    auth.merchantSecret,
		memberID:  auth.memberID,
		memberPin: pin,
		expire:    exp,
	}
}
