package forex

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

type Service interface {
	GetRate(ctx context.Context, from, to string, date time.Time) (decimal.Decimal, error)
}

func NewDummyService() Service {
	return &serviceImpl{}
}

type serviceImpl struct{}

func (serviceImpl) GetRate(_ context.Context, _, _ string, _ time.Time) (decimal.Decimal, error) {
	return decimal.NewFromInt(1), nil
}
