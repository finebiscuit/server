package date

import (
	"fmt"
	"time"
)

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

func NewFromString(s string) (Date, error) {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return Date{}, err
	}
	return NewFromTime(t), nil
}

func Today() Date {
	return NewFromTime(time.Now())
}

func (d Date) String() string {
	return fmt.Sprintf("%4d-%2d-%2d", d.year, d.month, d.day)
}

func (d Date) After(other Date) bool {
	if d.year > other.year {
		return true
	} else if d.year < other.year {
		return false
	}

	if d.month > other.month {
		return true
	} else if d.month < other.month {
		return false
	}

	if d.day > other.day {
		return true
	} else if d.day < other.day {
		return false
	}

	return false
}
