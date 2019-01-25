package gateway

import (
	"time"

	"github.com/fox-one/httpclient"
)

type MemberService struct {
	client   *Client
	authFunc func(expire time.Duration) httpclient.Authenticator
}

func (m *MemberClient) Service(name string) *MemberService {
	return &MemberService{
		client: m.client.Group(name),
		authFunc: func(expire time.Duration) httpclient.Authenticator {
			return m.Presign(expire)
		},
	}
}

func (m *MemberClient) ServiceWithPin(name, pin string) *MemberService {
	return &MemberService{
		client: m.client.Group(name),
		authFunc: func(expire time.Duration) httpclient.Authenticator {
			return m.PresignWithPin(pin, expire)
		},
	}
}

func (m *MerchantClient) MemberService(name, member string) *MemberService {
	return &MemberService{
		client: m.client.Group(name),
		authFunc: func(expire time.Duration) httpclient.Authenticator {
			return m.PresignMember(member, expire)
		},
	}
}

func (m *MemberService) Presign(expire time.Duration) httpclient.Authenticator {
	return m.authFunc(expire)
}
