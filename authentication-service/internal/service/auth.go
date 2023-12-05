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
	hmacSampleSecret     = "D)]5~3;@Xcf?cm"
	accessTokenDuration  = time.Hour          // Adjust the access token duration as needed
	refreshTokenDuration = time.Hour * 24 * 7 // Adjust the refresh token duration as needed
)

type AuthorizationService interface {
	GenerateToken(id int) (string, string, error)
	ParseToken(tokenString string, isRefresh bool) (int, error)
}

type AuthService struct {
	repo repository.DatabaseRepo
}

type authClaims struct {
	ID        int   `json:"id"`
	IsRefresh bool  `json:"is-refresh"`
	IssuedAt  int64 `json:"issuedAt"`
	ExpiresAt int64 `json:"expiresAt"`
	jwt.RegisteredClaims
}

func (a AuthService) GenerateToken(id int) (string, string, error) {

	claims := authClaims{
		ID:        id,
		IssuedAt:  time.Now().Unix(),
		IsRefresh: false,
		ExpiresAt: time.Now().Add(accessTokenDuration).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	accessTokenString, err := accessToken.SignedString([]byte(hmacSampleSecret))
	if err != nil {
		return "", "", err
	}

	// Generate refresh token
	refreshTokenClaims := authClaims{
		ID:        id,
		IssuedAt:  time.Now().Unix(),
		IsRefresh: true,
		ExpiresAt: time.Now().Add(refreshTokenDuration).Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)

	// Sign and get the complete encoded refresh token as a string using the secret
	refreshTokenString, err := refreshToken.SignedString([]byte(hmacSampleSecret))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func (a AuthService) ParseToken(tokenString string, isRefresh bool) (int, error) {

	token, err := jwt.ParseWithClaims(tokenString, &authClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(hmacSampleSecret), nil
	})
	if err != nil {
		log.Printf("error parsing tokenString, %v\n", err)
		return 0, err
	}

	if claims, ok := token.Claims.(*authClaims); ok && token.Valid {
		fmt.Printf("%v %v %v", claims.ID, claims.IsRefresh, claims.ExpiresAt)
		if claims.ExpiresAt < time.Now().Unix() {
			return 0, errors.New("token expired")
		}
		if claims.IsRefresh == isRefresh {
			return claims.ID, nil
		}
		return 0, errors.New("token is not refresh")
	} else {
		fmt.Printf("error in casting to authClaims, %v\n", err)
		return 0, err
	}
}

func NewAuthService(repo repository.DatabaseRepo) AuthorizationService {
	return &AuthService{
		repo: repo,
	}
}
