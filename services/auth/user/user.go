package user

import (
	"github.com/finebiscuit/server/model/buid"
	"golang.org/x/crypto/bcrypt"
)

type ID struct {
	buid.BUID
}

func ParseID(s string) (ID, error) {
	id, err := buid.Parse(s)
	return ID{id}, err
}

type User struct {
	ID             ID
	Email          string
	HashedPassword string
}

func NewFromEmailAndPassword(email, password string) (*User, error) {
	id, err := buid.New()
	if err != nil {
		return nil, err
	}

	pwHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u := &User{
		ID:             ID{id},
		Email:          email,
		HashedPassword: string(pwHash),
	}
	return u, nil
}

func (u *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password))
}
