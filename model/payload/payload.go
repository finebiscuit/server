package payload

type Payload struct {
	Scheme  Scheme
	Version string
	Blob    []byte
}

func New(scheme Scheme, version string, blob []byte) (Payload, error) {
	if !scheme.Valid() {
		return Payload{}, ErrInvalidScheme
	}
	if version == "" {
		return Payload{}, ErrEmptyVersion
	}
	p := Payload{
		Scheme:  scheme,
		Version: version,
		Blob:    blob,
	}
	return p, nil
}

func Must(p Payload, err error) Payload {
	if err != nil {
		panic(err)
	}
	return p
}
