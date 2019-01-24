package gateway

import (
	"time"
)

type MemberService struct {
	client   *Client
	authFunc func(expire time.Duration) Authenticator
}

func (m *Member) Service(name string) *MemberService {
	return &MemberService{
		client: m.client.Group(name),
		authFunc: func(expire time.Duration) Authenticator {
			return m.Presign(expire)
		},
	}
}

func (m *Member) ServiceWithPin(name, pin string) *MemberService {
	return &MemberService{
		client: m.client.Group(name),
		authFunc: func(expire time.Duration) Authenticator {
			return m.PresignWithPin(pin, expire)
		},
	}
}

func (m *Merchant) MemberService(name, member string) *MemberService {
	return &MemberService{
		client: m.client.Group(name),
		authFunc: func(expire time.Duration) Authenticator {
			return m.PresignMember(member, expire)
		},
	}
}

func (m *MemberService) Presign(expire time.Duration) Authenticator {
	return m.authFunc(expire)
}
