package user

type User struct {
	Username string `json:"username"`
	password string
	UserID   string `json:"id"`
}

func NewUser(Username string) *User {
	return &User{Username: Username}
}

func (u *User) WithID(Id string) *User {
	u.UserID = Id
	return u
}

func (u *User) WithPassword(Password string) *User {
	u.password = Password
	return u
}

func (u *User) GetPassword() string {
	return u.password
}
