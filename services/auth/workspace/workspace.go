package workspace

import (
	"fmt"

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

type Workspace struct {
	ID          ID
	UserID      user.ID
	DisplayName string
}

func New(userID user.ID, displayName string) (*Workspace, error) {
	id, err := buid.New()
	if err != nil {
		return nil, err
	}
	w := &Workspace{
		ID:          ID{id},
		UserID:      userID,
		DisplayName: displayName,
	}
	return w, nil
}

func (ws *Workspace) CompareAccessFor(userID user.ID) error {
	if ws.UserID != userID {
		return fmt.Errorf("permission denied")
	}
	return nil
}
