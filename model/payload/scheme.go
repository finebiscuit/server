package payload

type Scheme uint32

const (
	SchemeNone Scheme = iota
	SchemePlainProto
	SchemeMax
)

func NewScheme(i int) (Scheme, error) {
	s := Scheme(i)
	if !s.Valid() {
		return s, ErrInvalidScheme
	}
	return s, nil
}

func (s Scheme) Valid() bool {
	return s > SchemeNone && s < SchemeMax
}
