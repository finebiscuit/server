package balance

import (
	"github.com/finebiscuit/server/model/buid"
	"github.com/finebiscuit/server/services/accounting/kind"
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
