package buid

import "github.com/google/uuid"

type BUID uuid.UUID

func (id BUID) String() string {
	return uuid.UUID(id).String()
}
