package service

import (
	"authentication-service/internal/repository"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type AuthorizationService interface {
	GenerateToken(id int, username string) (string, error)
	ParseToken(tokenString string) (int, error)
}

type AuthService struct {
	repo repository.DatabaseRepo
}

type MyCustomClaims struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	ExpiresAt time.Time `json:"expiresAt"`
	jwt.RegisteredClaims
}

const (
	hmacSampleSecret = "secret"
	duration         = time.Hour * 24
)

func (a AuthService) GenerateToken(id int, username string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":        id,
		"username":  id,
		"issuedAt":  time.Now().Unix(),
		"expiresAt": time.Now().Add(duration).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(hmacSampleSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (a AuthService) ParseToken(tokenString string) (int, error) {

	type MyCustomClaims struct {
		ID       int `json:"id"`
		Username int `json:"username"`
		jwt.RegisteredClaims
	}

	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(hmacSampleSecret), nil
	})

	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		fmt.Printf("%v %v %v", claims.ID, claims.Username, claims.ExpiresAt)
		return claims.ID, nil
	} else {
		fmt.Println(err)
		return 0, err
	}
}

func NewAuthService(repo repository.DatabaseRepo) AuthorizationService {
	return &AuthService{
		repo: repo,
	}
}
