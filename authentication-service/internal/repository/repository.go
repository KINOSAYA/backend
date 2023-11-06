package repository

import (
	"authentication-service/internal/models"
)

type DatabaseRepo interface {
	GetUserByEmail(user models.User) (models.User, error)
	GetUserByUsername(user models.User) (models.User, error)

	AddUser(user models.User) (int, error)

	UpdateUser(user models.User) error

	DeleteUser(user models.User) error
}
