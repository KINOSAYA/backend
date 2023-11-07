package service

import (
	"authentication-service/internal/repository"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"time"
)

const (
	hmacSampleSecret = "secret"
	duration         = time.Hour * 24
)

type AuthorizationService interface {
	GenerateToken(id int, username string) (string, error)
	ParseToken(tokenString string) (int, string, error)
}

type AuthService struct {
	repo repository.DatabaseRepo
}

type MyCustomClaims struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	IssuedAt  int64  `json:"issuedAt"`
	ExpiresAt int64  `json:"expiresAt"`
	jwt.RegisteredClaims
}

func (a AuthService) GenerateToken(id int, username string) (string, error) {

	claims := MyCustomClaims{
		ID:        id,
		Username:  username,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(duration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(hmacSampleSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (a AuthService) ParseToken(tokenString string) (int, string, error) {

	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(hmacSampleSecret), nil
	})
	if err != nil {
		log.Printf("error parsing tokenString, %v\n", err)
		return 0, "", err
	}

	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		fmt.Printf("%v %v %v", claims.ID, claims.Username, claims.ExpiresAt)
		if claims.ExpiresAt < time.Now().Unix() {
			return 0, "", errors.New("token expired")
		}
		return claims.ID, claims.Username, nil
	} else {
		fmt.Printf("error in casting to MyCustomClaims, %v\n", err)
		return 0, "", err
	}
}

func NewAuthService(repo repository.DatabaseRepo) AuthorizationService {
	return &AuthService{
		repo: repo,
	}
}
