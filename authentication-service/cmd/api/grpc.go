package main

import (
	"authentication-service/internal/auth"
	"authentication-service/internal/models"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

type AuthServer struct {
	auth.UnimplementedAuthServiceServer
	Models models.Models
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
	res := &auth.UserResponse{
		Message: "Inserted user!",
		Data: &auth.ResponseData{
			ID:    uint64(id),
			Token: "tempToken",
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

	auth.RegisterAuthServiceServer(s, &AuthServer{Models: app.Models})

	log.Printf("gRPC Server started on port %s", gRpcPort)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to listen for gRPC: %v", err)
	}
}
