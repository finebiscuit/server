package user

import "context"

type Repository interface {
	Get(ctx context.Context, id ID) (*User, error)
	GetByLogin(ctx context.Context, login string) (*User, error)
	Create(ctx context.Context, u *User) error
	Update(ctx context.Context, u *User) error
}
