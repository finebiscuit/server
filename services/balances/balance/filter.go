package balance

import "github.com/finebiscuit/server/services/balances/kind"

type Filter struct {
	IDs     []ID
	KindIDs []kind.ID
}
