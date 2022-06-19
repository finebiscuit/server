package balances

import (
	"context"

	"github.com/finebiscuit/server/services/balances/balance"
)

type Service interface {
	ListBalances(ctx context.Context, filter balance.Filter) ([]*balance.WithEntry, error)
	CreateBalance(ctx context.Context, b *balance.Balance, e *balance.Entry) error
}

func NewService(tx TxFn) Service {
	return &serviceImpl{tx: tx}
}

type serviceImpl struct {
	tx TxFn
}

var _ Service = &serviceImpl{}

func (s serviceImpl) ListBalances(ctx context.Context, filter balance.Filter) ([]*balance.WithEntry, error) {
	var result []*balance.WithEntry
	err := s.tx(ctx, func(ctx context.Context, uow UnitOfWork) (err error) {
		result, err = uow.Balances().List(ctx, filter)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s serviceImpl) CreateBalance(ctx context.Context, b *balance.Balance, e *balance.Entry) error {
	err := s.tx(ctx, func(ctx context.Context, uow UnitOfWork) error {
		if err := uow.Balances().Create(ctx, b, e); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
