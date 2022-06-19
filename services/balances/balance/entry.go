package balance

import "github.com/finebiscuit/server/model/date"

type Entry struct {
	YMD date.Date
}

func NewEntry() (*Entry, error) {
	return NewEntryWithDate(date.Today())
}

func NewEntryWithDate(ymd date.Date) (*Entry, error) {
	e := &Entry{
		YMD: ymd,
	}
	return e, nil
}

func MustEntry(e *Entry, err error) *Entry {
	if err != nil {
		panic(err)
	}
	return e
}
