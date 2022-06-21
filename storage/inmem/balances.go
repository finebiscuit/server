package inmem

import (
	"context"

	"github.com/finebiscuit/server/model/date"
	"github.com/finebiscuit/server/services/balances/balance"
)

type accountingBalancesRepo struct {
	uow *unitOfWork
}

type StorageBalance struct {
	Balance    balance.Balance
	CurrentYMD date.Date
	Entries    map[date.Date]balance.Entry
}

func (b StorageBalance) toDomain() *balance.WithEntry {
	return &balance.WithEntry{
		Balance: b.Balance,
		Entry:   b.Entries[b.CurrentYMD],
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

func (r accountingBalancesRepo) Create(ctx context.Context, b *balance.Balance) error {
	r.uow.db.Balances[b.ID] = &StorageBalance{
		Balance: *b,
		Entries: make(map[date.Date]balance.Entry),
	}
	return nil
}

func (r accountingBalancesRepo) Update(ctx context.Context, b *balance.Balance) error {
	if dbBal := r.uow.db.Balances[b.ID]; dbBal != nil {
		dbBal.Balance = *b
	}
	return nil
}

func (r accountingBalancesRepo) GetEntry(
	ctx context.Context, balanceID balance.ID, entryYMD date.Date,
) (*balance.Entry, error) {
	b := r.uow.db.Balances[balanceID]
	if b == nil {
		return nil, balance.ErrNotFound
	}
	e, ok := b.Entries[entryYMD]
	if !ok {
		return nil, balance.ErrEntryNotFound
	}
	return &e, nil
}

func (r accountingBalancesRepo) CreateEntry(ctx context.Context, balanceID balance.ID, e *balance.Entry) error {
	dbBal, ok := r.uow.db.Balances[balanceID]
	if !ok {
		return balance.ErrNotFound
	}

	if _, ok := dbBal.Entries[e.YMD]; ok {
		return balance.ErrEntryAlreadyExists
	}

	dbBal.Entries[e.YMD] = *e
	if e.YMD.After(dbBal.CurrentYMD) {
		dbBal.CurrentYMD = e.YMD
	}
	return nil
}

func (r accountingBalancesRepo) UpdateEntry(ctx context.Context, balanceID balance.ID, e *balance.Entry) error {
	b := r.uow.db.Balances[balanceID]
	if b == nil {
		return balance.ErrNotFound
	}
	if _, ok := b.Entries[e.YMD]; !ok {
		return balance.ErrEntryNotFound
	}

	b.Entries[e.YMD] = *e
	return nil
}
