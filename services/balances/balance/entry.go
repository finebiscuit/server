package balance

import (
	"github.com/finebiscuit/server/model/date"
	"github.com/finebiscuit/server/model/payload"
)

type Entry struct {
	YMD     date.Date
	Payload payload.Payload
}

func NewEntry(p payload.Payload) (*Entry, error) {
	return NewEntryWithDate(date.Today(), p)
}

func NewEntryWithDate(ymd date.Date, p payload.Payload) (*Entry, error) {
	e := &Entry{
		YMD:     ymd,
		Payload: p,
	}
	return e, nil
}

func MustEntry(e *Entry, err error) *Entry {
	if err != nil {
		panic(err)
	}
	return e
}
