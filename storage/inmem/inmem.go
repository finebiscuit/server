package inmem

import (
	"context"
	"sync"

	"github.com/finebiscuit/server/services/auth"
	"github.com/finebiscuit/server/services/auth/session"
	"github.com/finebiscuit/server/services/auth/user"
	"github.com/finebiscuit/server/services/auth/workspace"
	"github.com/finebiscuit/server/services/balances"
	"github.com/finebiscuit/server/services/balances/balance"
)

type InMem struct {
	mu sync.Mutex
	DB *Database
}

func New() *InMem {
	return &InMem{
		DB: &Database{
			Balances:   make(map[balance.ID]*StorageBalance),
			Users:      make(map[user.ID]*user.User),
			Workspaces: make(map[workspace.ID]*workspace.Workspace),
			Sessions:   make(map[session.ID]*session.Session),
		},
	}
}

func (s *InMem) AuthTxFn() auth.TxFn {
	return func(ctx context.Context, fn func(ctx context.Context, uow auth.UnitOfWork) error) error {
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
	Balances   map[balance.ID]*StorageBalance
	Users      map[user.ID]*user.User
	Workspaces map[workspace.ID]*workspace.Workspace
	Sessions   map[session.ID]*session.Session
}
