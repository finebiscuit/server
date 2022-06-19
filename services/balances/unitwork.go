package balances

import (
	"context"

	"github.com/finebiscuit/server/services/balances/balance"
	"github.com/finebiscuit/server/services/balances/kind"
)

type TxFn func(ctx context.Context, fn func(ctx context.Context, uow UnitOfWork) error) error

type UnitOfWork interface {
	Kinds() kind.Repository
	Balances() balance.Repository
}