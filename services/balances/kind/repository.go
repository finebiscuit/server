package kind

import "context"

type Repository interface {
	Get(ctx context.Context, id ID) (*Kind, error)
	List(ctx context.Context) ([]*Kind, error)
	Create(ctx context.Context, k *Kind) error
	Update(ctx context.Context, k *Kind) error
}
