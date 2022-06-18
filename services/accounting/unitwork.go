package accounting

import (
	"context"

	"github.com/finebiscuit/server/services/accounting/balance"
	"github.com/finebiscuit/server/services/accounting/entry"
	"github.com/finebiscuit/server/services/accounting/kind"
)

type TxFn func(ctx context.Context, fn func(ctx context.Context, uow UnitOfWork) error) error

type UnitOfWork interface {
	Kinds() kind.Repository
	Balances() balance.Repository
	Entries() entry.Repository
}
