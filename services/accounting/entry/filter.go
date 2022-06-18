package entry

import "github.com/finebiscuit/server/services/accounting/balance"

type Filter struct {
	IDs        []ID
	BalanceIDs []balance.ID
}
