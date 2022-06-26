package inmem

import (
	"context"

	"github.com/finebiscuit/server/services/auth/session"
)

type sessionsRepo struct {
	uow *unitOfWork
}

func (r sessionsRepo) List(ctx context.Context) ([]*session.Session, error) {
	panic("not implemented")
}

func (r sessionsRepo) Get(ctx context.Context, id session.ID) (*session.Session, error) {
	s, ok := r.uow.db.Sessions[id]
	if !ok {
		return nil, session.ErrPermissionDenied
	}
	return s, nil
}

func (r sessionsRepo) Create(ctx context.Context, s *session.Session) error {
	r.uow.db.Sessions[s.ID] = s
	return nil
}
