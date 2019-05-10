package gateway

import (
	"time"

	"github.com/fox-one/httpclient"
)

type MemberService struct {
	*Client
	authFunc func(expire time.Duration) httpclient.Authenticator
}

func (m *MemberClient) Service(name string) *MemberService {
	return &MemberService{
		Client: m.Client.Group(name),
		authFunc: func(expire time.Duration) httpclient.Authenticator {
			return m.Presign(expire)
		},
	}
}

func (m *MemberClient) ServiceWithPin(name, pin string) *MemberService {
	return &MemberService{
		Client: m.Group(name),
		authFunc: func(expire time.Duration) httpclient.Authenticator {
			return m.PresignWithPin(pin, expire)
		},
	}
}

func (m *MemberService) Presign(expire time.Duration) httpclient.Authenticator {
	return m.authFunc(expire)
}
