package date

import "time"

type Date struct {
	year  int
	month time.Month
	day   int
}

func New(year int, month time.Month, day int) Date {
	return Date{
		year:  year,
		month: month,
		day:   day,
	}
}

func NewFromTime(t time.Time) Date {
	return New(t.Year(), t.Month(), t.Day())
}

func Today() Date {
	return NewFromTime(time.Now())
}
