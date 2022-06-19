package balances_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/finebiscuit/server/services/balances"
	"github.com/finebiscuit/server/services/balances/balance"
	"github.com/finebiscuit/server/storage/inmem"
)

func TestService_ListBalances(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		db := inmem.New()
		svc := balances.NewService(db.BalancesTxFn())

		result, err := svc.ListBalances(context.Background(), balance.Filter{})
		require.NoError(t, err)

		expected := []*balance.WithEntry{}
		assert.Equal(t, expected, result)
	})

	t.Run("Success", func(t *testing.T) {
		db := inmem.New()
		svc := balances.NewService(db.BalancesTxFn())

		b := balance.Must(balance.New("", ""))
		e := balance.MustEntry(balance.NewEntry())
		db.DB.Balances[b.ID] = &inmem.StorageBalance{
			Balance:      *b,
			CurrentEntry: *e,
		}

		result, err := svc.ListBalances(context.Background(), balance.Filter{})
		require.NoError(t, err)

		expected := []*balance.WithEntry{
			{
				Balance: *b,
				Entry:   *e,
			},
		}
		assert.Equal(t, expected, result)
	})
}

func TestService_CreateBalance(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db := inmem.New()
		svc := balances.NewService(db.BalancesTxFn())

		b := balance.Must(balance.New("", ""))
		e := balance.MustEntry(balance.NewEntry())

		err := svc.CreateBalance(context.Background(), b, e)
		require.NoError(t, err)

		expected := &inmem.StorageBalance{
			Balance: *b,
			CurrentEntry: *e,
		}
		assert.Equal(t, expected, db.DB.Balances[b.ID])
	})
}
