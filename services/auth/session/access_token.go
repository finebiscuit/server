package session

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/finebiscuit/server/services/auth/user"
	"github.com/finebiscuit/server/services/auth/workspace"
)

type AccessToken struct {
	SessionID   ID
	UserID      user.ID
	WorkspaceID workspace.ID
	IssuedAt    time.Time
	ExpiresAt   time.Time
}

type tokenClaims struct {
	jwt.StandardClaims
	SessionID   string `json:"session_id"`
	WorkspaceID string `json:"workspace_id"`
}

func ParseAccessToken(token string, key []byte) (*AccessToken, error) {
	parsedTok, err := jwt.ParseWithClaims(token, &tokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := parsedTok.Claims.(*tokenClaims)
	if !ok || !parsedTok.Valid {
		return nil, fmt.Errorf("invalid JWT")
	}

	sessID, err := ParseID(claims.SessionID)
	if err != nil {
		return nil, err
	}

	userID, err := user.ParseID(claims.Subject)
	if err != nil {
		return nil, err
	}

	workspaceID, err := workspace.ParseID(claims.WorkspaceID)
	if err != nil {
		return nil, err
	}

	at := &AccessToken{
		SessionID:   sessID,
		UserID:      userID,
		WorkspaceID: workspaceID,
		IssuedAt:    time.Unix(claims.IssuedAt, 0),
		ExpiresAt:   time.Unix(claims.ExpiresAt, 0),
	}
	return at, nil
}

func (s *Session) GenerateAccessToken(wsID workspace.ID, expiresAt time.Time) (*AccessToken, error) {
	issuedAt := time.Now()
	if s.ExpiresAt.Before(issuedAt) {
		return nil, fmt.Errorf("session expired")
	}
	tok := &AccessToken{
		SessionID:   s.ID,
		UserID:      s.UserID,
		WorkspaceID: wsID,
		ExpiresAt:   expiresAt,
		IssuedAt:    issuedAt,
	}
	return tok, nil
}

func (tok *AccessToken) SignedString(key []byte) (string, error) {
	j := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":          tok.UserID.String(),
		"iat":          tok.IssuedAt.Unix(),
		"exp":          tok.ExpiresAt.Unix(),
		"session_id":   tok.SessionID.String(),
		"workspace_id": tok.WorkspaceID.String(),
	})
	s, err := j.SignedString(key)
	if err != nil {
		return "", err
	}
	return s, nil
}

func (tok *AccessToken) GetTTL() time.Duration {
	return tok.ExpiresAt.Sub(time.Now())
}
