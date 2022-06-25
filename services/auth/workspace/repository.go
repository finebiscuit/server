package workspace

import (
	"context"

	"github.com/finebiscuit/server/services/auth/user"
)

type Repository interface {
	Get(ctx context.Context, id ID) (*Workspace, error)
	List(ctx context.Context, userID user.ID) ([]*Workspace, error)
	Create(ctx context.Context, w *Workspace) error
}
