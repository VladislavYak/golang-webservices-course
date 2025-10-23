package user

type User struct {
	Username string `json:"username"`
	// Password string
	UserID string `json:"id"`
}

func NewUser(Username string) *User {
	return &User{Username: Username}
}

// yakovlev: это надо бы убрать
func (u *User) WithID(Id string) *User {
	u.UserID = Id
	return u
}
