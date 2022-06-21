package balance

import (
	"context"

	"github.com/finebiscuit/server/model/date"
)

type Repository interface {
	Get(ctx context.Context, id ID) (*WithEntry, error)
	List(ctx context.Context, filter Filter) ([]*WithEntry, error)
	Create(ctx context.Context, b *Balance) error
	Update(ctx context.Context, b *Balance) error

	GetEntry(ctx context.Context, balanceID ID, entryYMD date.Date) (*Entry, error)
	CreateEntry(ctx context.Context, balanceID ID, e *Entry) error
	UpdateEntry(ctx context.Context, balanceID ID, e *Entry) error
}
