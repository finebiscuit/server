package workspace

import authv1 "github.com/finebiscuit/proto/biscuit/auth/v1"

func NewFromProto(proto *authv1.Workspace) (*Workspace, error) {
	id, err := ParseID(proto.GetId())
	if err != nil {
		return nil, err
	}
	return &Workspace{
		ID:          id,
		DisplayName: proto.GetDisplayName(),
	}, nil
}
