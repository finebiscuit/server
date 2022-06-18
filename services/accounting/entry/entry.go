package entry

import (
	"github.com/finebiscuit/server/services/accounting/balance"
)

type ID string

type Entry struct {
	ID        ID
	BalanceID balance.ID
}

type WithBalance struct {
	*Entry
	Balance *balance.Balance
}
