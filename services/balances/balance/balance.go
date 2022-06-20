package balance

import (
	"github.com/finebiscuit/server/model/buid"
	"github.com/finebiscuit/server/model/payload"
)

type ID struct {
	buid.BUID
}

func ParseID(s string) (ID, error) {
	id, err := buid.Parse(s)
	return ID{id}, err
}

type Balance struct {
	ID         ID
	TypeID     string
	CurrencyID string
	Payload    payload.Payload
}

type WithEntry struct {
	Balance
	Entry Entry
}

func New(typeID, currencyID string, p payload.Payload) (*Balance, error) {
	id, err := buid.New()
	if err != nil {
		return nil, err
	}
	b := &Balance{
		ID:         ID{id},
		TypeID:     typeID,
		CurrencyID: currencyID,
		Payload:    p,
	}
	return b, nil
}

func Must(b *Balance, err error) *Balance {
	if err != nil {
		panic(err)
	}
	return b
}
