package accounting

import (
	"context"

	"github.com/finebiscuit/server/services/accounting/balance"
	"github.com/finebiscuit/server/services/accounting/entry"
	// "github.com/finebiscuit/server/services/accounting/kind"
)

type Service interface {
	// ListKinds(ctx context.Context) ([]*kind.Kind, error)
	// CreateOrUpdateKind(ctx context.Context, k *kind.Kind) error

	ListBalances(ctx context.Context, filter balance.Filter) ([]*entry.WithBalance, error)
	CreateBalance(ctx context.Context, b *balance.Balance, e *entry.Entry) error
	// UpdateBalance(ctx context.Context, b *balance.Balance) error

	// CreateOrUpdateEntry(ctx context.Context, e *entry.Entry) error
}

func NewService(tx TxFn) Service {
	return &serviceImpl{tx: tx}
}

type serviceImpl struct {
	tx TxFn
}

var _ Service = &serviceImpl{}

func (s serviceImpl) ListBalances(ctx context.Context, filter balance.Filter) ([]*entry.WithBalance, error) {
	var result []*entry.WithBalance
	err := s.tx(ctx, func(ctx context.Context, uow UnitOfWork) error {
		bals, err := uow.Balances().List(ctx, filter)
		if err != nil {
			return err
		}

		ids := make([]balance.ID, 0, len(bals))
		for _, b := range bals {
			ids = append(ids, b.ID)
		}

		entryByBalance, err := uow.Entries().GetLatestByBalance(ctx, ids)
		if err != nil {
			return err
		}

		result = make([]*entry.WithBalance, 0, len(bals))
		for _, b := range bals {
			if e, ok := entryByBalance[b.ID]; ok {
				result = append(result, &entry.WithBalance{
					Balance: b,
					Entry:   e,
				})
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s serviceImpl) CreateBalance(ctx context.Context, b *balance.Balance, e *entry.Entry) error {
	err := s.tx(ctx, func(ctx context.Context, uow UnitOfWork) error {
		if err := uow.Balances().Create(ctx, b); err != nil {
			return err
		}

		e.BalanceID = b.ID
		if err := uow.Entries().Create(ctx, e); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
