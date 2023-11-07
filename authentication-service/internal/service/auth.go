package service

import (
	"authentication-service/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type AuthorizationService interface {
	GenerateToken(id int) (string, error)
	ParseToken(token string) (int, error)
}

type AuthService struct {
	repo repository.DatabaseRepo
}

const (
	hmacSampleSecret = "secret"
	duration         = time.Hour * 24
)

func (a AuthService) GenerateToken(id int) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":        id,
		"issuedAt":  time.Now().Unix(),
		"expiresAt": time.Now().Add(duration).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(hmacSampleSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (a AuthService) ParseToken(token string) (int, error) {
	//TODO implement me
	panic("implement me")
}

func NewAuthService(repo repository.DatabaseRepo) AuthorizationService {
	return &AuthService{
		repo: repo,
	}
}
