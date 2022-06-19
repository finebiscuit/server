package balance

import (
	"github.com/finebiscuit/server/model/buid"
	"github.com/finebiscuit/server/services/balances/kind"
)

type ID struct {
	buid.BUID
}

type Balance struct {
	ID     ID
	KindID kind.ID
}

type WithEntry struct {
	Balance
	Entry Entry
}
