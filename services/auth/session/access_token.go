package session

import (
	"fmt"
	"time"

	"github.com/finebiscuit/server/services/auth/user"
	"github.com/finebiscuit/server/services/auth/workspace"
)

type AccessToken struct {
	SessionID   ID
	UserID      user.ID
	WorkspaceID workspace.ID
	ExpiresAt   time.Time
}

func (s *Session) GenerateAccessToken(wsID workspace.ID, expiresAt time.Time) (*AccessToken, error) {
	if s.ExpiresAt.After(time.Now()) {
		return nil, fmt.Errorf("session expired")
	}
	tok := &AccessToken{
		SessionID:   s.ID,
		UserID:      s.UserID,
		WorkspaceID: wsID,
		ExpiresAt:   expiresAt,
	}
	return tok, nil
}
