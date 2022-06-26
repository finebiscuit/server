package inmem

import (
	"context"

	"github.com/finebiscuit/server/services/auth/user"
)

type usersRepo struct {
	uow *unitOfWork
}

func (r usersRepo) Get(ctx context.Context, id user.ID) (*user.User, error) {
	u, ok := r.uow.db.Users[id]
	if !ok {
		return nil, user.ErrNotFound
	}
	return u, nil
}

func (r usersRepo) GetByLogin(ctx context.Context, login string) (*user.User, error) {
	for _, u := range r.uow.db.Users {
		if u.Email == login {
			return u, nil
		}
	}
	return nil, user.ErrInvalidCredentials
}

func (r usersRepo) Create(ctx context.Context, u *user.User) error {
	if _, ok := r.uow.db.Users[u.ID]; ok {
		return user.ErrAlreadyExists
	}
	for _, usr := range r.uow.db.Users {
		if u.Email == usr.Email {
			return user.ErrEmailAlreadyTaken
		}
	}
	r.uow.db.Users[u.ID] = u
	return nil
}

func (r usersRepo) Update(ctx context.Context, u *user.User) error {
	r.uow.db.Users[u.ID] = u
	return nil
}
