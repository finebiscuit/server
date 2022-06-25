package user

type Identity struct {
	UserID ID
}

func (u *User) AsIdentity() *Identity {
	return &Identity{
		UserID: u.ID,
	}
}
