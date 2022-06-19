package inmem

import (
	"context"

	"github.com/finebiscuit/server/services/balances/balance"
)

type accountingBalancesRepo struct {
	uow *unitOfWork
}

type StorageBalance struct {
	Balance      balance.Balance
	CurrentEntry balance.Entry
}

func (b StorageBalance) toDomain() *balance.WithEntry {
	return &balance.WithEntry{
		Balance: b.Balance,
		Entry:   b.CurrentEntry,
	}
}

func (r accountingBalancesRepo) Get(ctx context.Context, id balance.ID) (*balance.WithEntry, error) {
	b := r.uow.db.Balances[id]
	if b == nil {
		return nil, balance.ErrNotFound
	}
	return b.toDomain(), nil
}

func (r accountingBalancesRepo) List(ctx context.Context, filter balance.Filter) ([]*balance.WithEntry, error) {
	result := make([]*balance.WithEntry, 0, len(r.uow.db.Balances))
	for _, b := range r.uow.db.Balances {
		result = append(result, b.toDomain())
	}
	return result, nil
}

func (r accountingBalancesRepo) Create(ctx context.Context, b *balance.Balance, e *balance.Entry) error {
	r.uow.db.Balances[b.ID] = &StorageBalance{
		Balance:      *b,
		CurrentEntry: *e,
	}
	return nil
}

func (r accountingBalancesRepo) Update(ctx context.Context, b *balance.Balance) error {
	if dbBal := r.uow.db.Balances[b.ID]; dbBal != nil {
		dbBal.Balance = *b
	}
	return nil
}

func (r accountingBalancesRepo) UpsertEntry(ctx context.Context, balanceID balance.ID, e *balance.Entry) error {
	dbBal, ok := r.uow.db.Balances[balanceID]
	if !ok {
		return balance.ErrNotFound
	}
	
	// TODO: instead of always overwriting the currentEntry, keep also a map/list with all entries
	dbBal.CurrentEntry = *e
	return nil
}
