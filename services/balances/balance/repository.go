package balance

import (
	"context"
)

type Repository interface {
	Get(ctx context.Context, id ID) (*WithEntry, error)
	List(ctx context.Context, filter Filter) ([]*WithEntry, error)
	Create(ctx context.Context, b *Balance, e *Entry) error
	Update(ctx context.Context, b *Balance) error
	UpsertEntry(ctx context.Context, balanceID ID, e *Entry) error
}
