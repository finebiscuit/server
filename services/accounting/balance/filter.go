package balance

import "github.com/finebiscuit/server/services/accounting/kind"

type Filter struct {
	IDs     []ID
	KindIDs []kind.ID
}
