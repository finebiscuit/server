package balances

import (
	"context"

	"github.com/finebiscuit/server/services/balances/balance"
)

type TxFn func(ctx context.Context, fn func(ctx context.Context, uow UnitOfWork) error) error

type UnitOfWork interface {
	Balances() balance.Repository
}
