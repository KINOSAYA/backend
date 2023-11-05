package main

import (
	_ "broker-service/cmd/api/docs"
	"broker-service/internal/auth"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"time"
)

// broker is a sample API endpoint that returns a JSON response.
// @Summary hit the broker
// @Description Returns a JSON response with success status.
// @ID get-sample-response
// @Produce json
// @Success 200 {object} jsonResponse
// @Router / [get]
func (app *Config) broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Successfully hit the broker",
	}
	app.writeJSON(w, http.StatusOK, payload)
}

type requestPayload struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterUser is an API endpoint that register a user and returns a JSON response.
// @Tags user
// @Summary Register a new user
// @Description Registers a new user with the specified data.
// @Accept json
// @Produce json
// @Param requestPayload body requestPayload true "User data"
// @Success 202 {object} jsonResponse "Successful registration"
// @Router /user [post]
func (app *Config) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var requestPayload requestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	//TODO rewrite hardcode authentication host!
	conn, err := grpc.Dial(fmt.Sprintf("authentication-service:%s", authGrpcPort), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer conn.Close()

	c := auth.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	userResponse, err := c.RegisterUser(ctx, &auth.UserRequest{
		UserEntry: &auth.User{
			Username: requestPayload.Username,
			Email:    requestPayload.Email,
			Password: requestPayload.Password,
		},
	})
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: userResponse.Result,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}
