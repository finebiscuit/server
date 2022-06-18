package entry

import (
	"context"

	"github.com/finebiscuit/server/services/accounting/balance"
)

type Repository interface {
	Get(ctx context.Context, id ID) (*Entry, error)
	GetLatestByBalance(ctx context.Context, balanceIDs []balance.ID) (map[balance.ID]*Entry, error)
	Create(ctx context.Context, e *Entry) error
	Update(ctx context.Context, e *Entry) error
}
