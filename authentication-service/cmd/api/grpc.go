package main

import (
	"authentication-service/internal/auth"
	"authentication-service/internal/models"
	"authentication-service/internal/service"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

type AuthServer struct {
	auth.UnimplementedAuthServiceServer
	Models  models.Models
	Service service.AuthorizationService
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
		res := &auth.UserResponse{Message: fmt.Sprintf("smth got wrong: %s", err)}
		return res, err
	}

	// return a response
	//TODO generate token
	token, err := a.Service.GenerateToken(id, user.Username)
	if err != nil {
		return nil, err
	}
	res := &auth.UserResponse{
		Message: "Inserted user!",
		Data: &auth.ResponseData{
			ID:    uint64(id),
			Token: token,
		},
	}
	return res, nil
}

func gRPCListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	s := grpc.NewServer()

	auth.RegisterAuthServiceServer(s, &AuthServer{
		Models:  app.Models,
		Service: app.Service,
	})

	log.Printf("gRPC Server started on port %s", gRpcPort)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to listen for gRPC: %v", err)
	}
}

func (a AuthServer) AuthUser(ctx context.Context, req *auth.UserRequest) (*auth.UserResponse, error) {
	input := req.GetUserEntry()

	//TODO search for user in DB
	user := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
	}

	id, err := app.DB.Authenticate(user.Email, user.Username, user.Password)
	if err != nil {
		return nil, err
	}

	//TODO generate token
	res := &auth.UserResponse{
		Message: "Got user from DB",
		Data: &auth.ResponseData{
			ID:    uint64(id),
			Token: "tempToken",
		},
	}
	return res, nil
}

func (a AuthServer) CheckToken(ctx context.Context, req *auth.TokenRequest) (*auth.TokenResponse, error) {
	tokenString := req.GetTokenString()
	log.Println("token string is", tokenString)
	id, username, err := a.Service.ParseToken(tokenString)
	if err != nil {
		return nil, err
	}

	res := &auth.TokenResponse{
		Id:       uint64(id),
		Username: username,
	}

	return res, nil

}
