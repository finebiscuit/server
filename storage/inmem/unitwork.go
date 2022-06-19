package inmem

import (
	"context"

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
