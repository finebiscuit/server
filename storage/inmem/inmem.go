package inmem

import (
	"context"
	"sync"

	"github.com/finebiscuit/server/services/balances"
	"github.com/finebiscuit/server/services/balances/balance"
	"github.com/finebiscuit/server/services/balances/kind"
)

type InMem struct {
	mu sync.Mutex
	DB *Database
}

func New() *InMem {
	return &InMem{
		DB: &Database{
			Kinds:    make(map[kind.ID]*kind.Kind),
			Balances: make(map[balance.ID]*StorageBalance),
		},
	}
}

func (s *InMem) BalancesTxFn() balances.TxFn {
	return func(ctx context.Context, fn func(ctx context.Context, uow balances.UnitOfWork) error) error {
		s.mu.Lock()
		defer s.mu.Unlock()

		uow := s.newUnitOfWork(ctx)
		if err := fn(ctx, uow); err != nil {
			return err
		}
		s.DB = uow.db
		return nil
	}
}

type Database struct {
	Kinds    map[kind.ID]*kind.Kind
	Balances map[balance.ID]*StorageBalance
}
