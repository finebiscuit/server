package serverinfo

import authv1 "github.com/finebiscuit/proto/biscuit/auth/v1"

type ServerInfo struct {
	ForexSupported bool
}

func (s ServerInfo) AsProto() *authv1.ServerInfo {
	return &authv1.ServerInfo{
		ForexSupported: s.ForexSupported,
	}
}
