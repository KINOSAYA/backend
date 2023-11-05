package models

import "time"

func New() Models {
	return Models{
		UserEntry: User{},
	}
}

type Models struct {
	UserEntry User
}

type User struct {
	ID        int
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
