package repository

import (
	"authentication-service/internal/models"
)

type DatabaseRepo interface {
	Authenticate(email, username, password string) (int, string, error)

	AddUser(user models.User) (int, error)

	UpdateUser(user models.User) error

	DeleteUser(user models.User) error
}
