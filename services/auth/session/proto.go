package session

import authv1 "github.com/finebiscuit/proto/biscuit/auth/v1"

func NewFromProto(proto *authv1.Session) (*Session, error) {
	id, err := ParseID(proto.GetId())
	if err != nil {
		return nil, err
	}

	s := &Session{
		ID:        id,
		PlainCode: proto.Code,
		CreatedAt: proto.CreatedAt.AsTime(),
		ExpiresAt: proto.ExpiresAt.AsTime(),
	}
	return s, nil
}
