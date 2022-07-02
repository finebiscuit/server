package balances

import (
	"context"

	"github.com/finebiscuit/server/services/balances/balance"
)

type Service interface {
	GetBalance(ctx context.Context, id balance.ID) (*balance.WithEntry, error)
	ListBalances(ctx context.Context, filter balance.Filter) ([]*balance.WithEntry, error)
	CreateBalance(ctx context.Context, b *balance.Balance, e *balance.Entry) (*balance.WithEntry, error)
	CreateEntry(ctx context.Context, id balance.ID, e *balance.Entry) error
	UpdateEntry(ctx context.Context, id balance.ID, e *balance.Entry, versionMatch string) error
}

func NewService(tx TxFn) Service {
	return &serviceImpl{tx: tx}
}

type serviceImpl struct {
	tx TxFn
}

var _ Service = &serviceImpl{}

func (s serviceImpl) GetBalance(ctx context.Context, id balance.ID) (*balance.WithEntry, error) {
	var result *balance.WithEntry
	err := s.tx(ctx, func(ctx context.Context, uow UnitOfWork) (err error) {
		result, err = uow.Balances().Get(ctx, id)
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

func (s serviceImpl) CreateBalance(
	ctx context.Context, b *balance.Balance, e *balance.Entry,
) (*balance.WithEntry, error) {
	err := s.tx(ctx, func(ctx context.Context, uow UnitOfWork) (err error) {
		if err := uow.Balances().Create(ctx, b); err != nil {
			return err
		}

		if err := uow.Balances().CreateEntry(ctx, b.ID, e); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// This must be a separate transaction.
	var result *balance.WithEntry
	err = s.tx(ctx, func(ctx context.Context, uow UnitOfWork) error {
		result, err = uow.Balances().Get(ctx, b.ID)
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

func (s serviceImpl) CreateEntry(ctx context.Context, balanceID balance.ID, e *balance.Entry) error {
	err := s.tx(ctx, func(ctx context.Context, uow UnitOfWork) error {
		if err := uow.Balances().CreateEntry(ctx, balanceID, e); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s serviceImpl) UpdateEntry(
	ctx context.Context, balanceID balance.ID, e *balance.Entry, versionMatch string,
) error {
	err := s.tx(ctx, func(ctx context.Context, uow UnitOfWork) error {
		b, err := uow.Balances().Get(ctx, balanceID)
		if err != nil {
			return err
		}

		prev, err := uow.Balances().GetEntry(ctx, b.ID, e.YMD)
		if err != nil {
			return err
		}

		if prev.Payload.Version != versionMatch {
			return balance.ErrVersionMismatch
		}

		if e.YMD == b.Entry.YMD || e.YMD.After(b.Entry.YMD) {
			e.IsCurrent = true
		}

		if err := uow.Balances().UpdateEntry(ctx, balanceID, e); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
