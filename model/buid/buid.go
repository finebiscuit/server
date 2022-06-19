package buid

import "github.com/google/uuid"

type BUID uuid.UUID

func New() (BUID, error) {
	u, err := uuid.NewRandom()
	return BUID(u), err
}

func (id BUID) String() string {
	return uuid.UUID(id).String()
}
