package date

type Date struct {
	year       uint16
	month, day uint8
}

func New(year, month, day int) Date {
	return Date{
		year:  uint16(year),
		month: uint8(month),
		day:   uint8(day),
	}
}
