package balance

import (
	"github.com/finebiscuit/server/services/accounting/kind"
)

type ID string

type Balance struct {
	ID     ID
	KindID kind.ID
}
