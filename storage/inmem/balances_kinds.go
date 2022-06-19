package inmem

import (
	"context"

	"github.com/finebiscuit/server/services/balances/kind"
)

type accountingKindsRepo struct {
	uow *unitOfWork
}

func (r accountingKindsRepo) Get(ctx context.Context, id kind.ID) (*kind.Kind, error) {
	k := r.uow.db.Kinds[id]
	if k == nil {
		return nil, kind.ErrNotFound
	}
	return k, nil
}

func (r accountingKindsRepo) List(ctx context.Context) ([]*kind.Kind, error) {
	result := make([]*kind.Kind, 0, len(r.uow.db.Kinds))
	for _, v := range r.uow.db.Kinds {
		result = append(result, v)
	}
	return result, nil
}

func (r accountingKindsRepo) Create(ctx context.Context, k *kind.Kind) error {
	r.uow.db.Kinds[k.ID] = k
	return nil
}

func (r accountingKindsRepo) Update(ctx context.Context, k *kind.Kind) error {
	r.uow.db.Kinds[k.ID] = k
	return nil
}
