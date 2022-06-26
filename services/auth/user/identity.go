package user

type Identity struct {
	UserID    string
	UserEmail string
}

func (u *User) AsIdentity() *Identity {
	return &Identity{
		UserID:    u.ID.String(),
		UserEmail: u.Email,
	}
}
