package session

import "context"

type Repository interface {
	List(ctx context.Context) ([]*Session, error)
	Get(ctx context.Context, id ID) (*Session, error)
	Create(ctx context.Context, s *Session) error
}
