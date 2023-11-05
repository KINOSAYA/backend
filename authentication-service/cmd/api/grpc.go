package main

import (
	"authentication-service/internal/auth"
	"authentication-service/internal/models"
	"context"
	"fmt"
)

type AuthServer struct {
	auth.UnimplementedAuthServiceServer
	User models.User
}

func (a AuthServer) RegisterUser(ctx context.Context, req *auth.UserRequest) (*auth.UserResponse, error) {
	input := req.GetUserEntry()

	// insert data
	user := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
	}

	id, err := app.DB.AddUser(user)
	if err != nil {
		res := &auth.UserResponse{Result: fmt.Sprintf("smth got wrong: %s", err)}
		return res, err
	}

	// return a response
	res := &auth.UserResponse{Result: fmt.Sprintf("Inserted user!%d", id)}
	return res, nil
}
