package inmem

import (
	"context"

	"github.com/finebiscuit/server/services/auth/user"
	"github.com/finebiscuit/server/services/auth/workspace"
)

type workspacesRepo struct {
	uow *unitOfWork
}

func (r workspacesRepo) Get(ctx context.Context, id workspace.ID) (*workspace.Workspace, error) {
	w, ok := r.uow.db.Workspaces[id]
	if !ok {
		return nil, workspace.ErrNotFound
	}
	return w, nil
}

func (r workspacesRepo) List(ctx context.Context, userID user.ID) ([]*workspace.Workspace, error) {
	result := make([]*workspace.Workspace, 0)
	for _, w := range r.uow.db.Workspaces {
		if w.UserID == userID {
			result = append(result, w)
		}
	}
	return result, nil
}

func (r workspacesRepo) Create(ctx context.Context, w *workspace.Workspace) error {
	r.uow.db.Workspaces[w.ID] = w
	return nil
}
