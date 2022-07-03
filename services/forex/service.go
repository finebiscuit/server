package forex

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

type Service interface {
	GetRate(ctx context.Context, base, target string, date time.Time) (decimal.Decimal, error)
	ListRates(ctx context.Context, base string, targets []string, date time.Time) (map[string]decimal.Decimal, error)
}

func NewDummyService() Service {
	return &serviceImpl{}
}

type serviceImpl struct{}

func (serviceImpl) GetRate(_ context.Context, _, _ string, _ time.Time) (decimal.Decimal, error) {
	return decimal.NewFromInt(1), nil
}

func (i serviceImpl) ListRates(_ context.Context, _ string, targets []string, _ time.Time) (map[string]decimal.Decimal, error) {
	m := make(map[string]decimal.Decimal)
	for _, t := range targets {
		m[t] = decimal.NewFromInt(1)
	}
	return m, nil
}
