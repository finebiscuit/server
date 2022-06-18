package balance

import (
	"context"
)

type Repository interface {
	Get(ctx context.Context, id ID) (*Balance, error)
	List(ctx context.Context, filter Filter) ([]*Balance, error)
	Create(ctx context.Context, b *Balance) error
	Update(ctx context.Context, b *Balance) error
}
