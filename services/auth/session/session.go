package session

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/finebiscuit/server/model/buid"
	"github.com/finebiscuit/server/services/auth/user"
)

type ID struct {
	buid.BUID
}

func ParseID(s string) (ID, error) {
	id, err := buid.Parse(s)
	return ID{id}, err
}

type Session struct {
	ID         ID
	PlainCode  string
	HashedCode string
	UserID     user.ID
	CreatedAt  time.Time
	ExpiresAt  time.Time
}

func New(userID user.ID, ttl time.Duration) (*Session, error) {
	id, err := buid.New()
	if err != nil {
		return nil, err
	}

	b := make([]byte, 40)
	_, err = rand.Read(b)
	if err != nil {
		return nil, err
	}
	code := base64.URLEncoding.EncodeToString(b)
	codeHash, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	t := time.Now()

	sess := &Session{
		ID:         ID{id},
		PlainCode:  code,
		HashedCode: string(codeHash),
		UserID:     userID,
		CreatedAt:  t,
		ExpiresAt:  t.Add(ttl),
	}
	return sess, nil
}

func (s *Session) CompareCode(code string) error {
	return bcrypt.CompareHashAndPassword([]byte(s.HashedCode), []byte(code))
}
