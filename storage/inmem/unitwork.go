package inmem

import (
	"context"

	"github.com/finebiscuit/server/services/auth/session"
	"github.com/finebiscuit/server/services/auth/user"
	"github.com/finebiscuit/server/services/auth/workspace"
	"github.com/finebiscuit/server/services/balances/balance"
)

type unitOfWork struct {
	db *Database
}

func (s *InMem) newUnitOfWork(ctx context.Context) *unitOfWork {
	return &unitOfWork{
		db: s.DB, // TODO: this needs to be deep-copied
	}
}

func (uow *unitOfWork) Balances() balance.Repository {
	return &accountingBalancesRepo{uow: uow}
}

func (uow *unitOfWork) Users() user.Repository {
	return &usersRepo{uow: uow}
}

func (uow *unitOfWork) Workspaces() workspace.Repository {
	return &workspacesRepo{uow: uow}
}

func (uow *unitOfWork) Sessions() session.Repository {
	return &sessionsRepo{uow: uow}
}
