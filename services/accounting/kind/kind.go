package kind

import "github.com/finebiscuit/server/model/buid"

type ID struct {
	buid.BUID
}

type Kind struct {
	ID ID
}
