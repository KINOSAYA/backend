package handlers

import (
	_ "broker-service/docs"
	"broker-service/internal/auth"
	"broker-service/internal/helpers"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"strings"
	"time"
)

type Handler interface {
	Broker(w http.ResponseWriter, r *http.Request)
	RegisterUser(w http.ResponseWriter, r *http.Request)
	LoginUser(w http.ResponseWriter, r *http.Request)
}

func NewBrokerHandler(AuthGrpcPort, authHost string) Handler {
	return &brokerHandler{
		AuthGrpcPort: AuthGrpcPort,
		AuthHost:     authHost,
	}
}

type brokerHandler struct {
	AuthGrpcPort string
	AuthHost     string
}

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// Broker is a sample API endpoint that returns a JSON response.
// @Summary hit the broker
// @Description Returns a JSON response with success status.
// @ID get-sample-response
// @Produce json
// @Success 200 {object} jsonResponse
// @Router / [get]
func (app *brokerHandler) Broker(w http.ResponseWriter, r *http.Request) {
	payload := helpers.JsonResponse{
		Error:   false,
		Message: "Successfully hit the broker",
	}
	_ = helpers.WriteJSON(w, http.StatusOK, payload)
}

type requestPayload struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterUser is an API endpoint that register a user and returns a JSON response.
// @Tags Auth
// @Summary Register a new user
// @Description Registers a new user with the specified data.
// @Accept json
// @Produce json
// @Param requestPayload body requestPayload true "User data"
// @Success 202 {object} jsonResponse "Successful registration"
// @Failure 401 {object} jsonResponse "Invalid credentials"
// @Router /auth/login [post]
func (app *brokerHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var requestPayload requestPayload

	err := helpers.ReadJSON(w, r, &requestPayload)
	if err != nil {
		_ = helpers.ErrorJSON(w, err)
		return
	}

	//TODO rewrite hardcode authentication host!
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", app.AuthHost, app.AuthGrpcPort), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
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
		if strings.Contains(err.Error(), "this email has already been used") {
			_ = helpers.ErrorJSON(w, errors.New("this email has already been used"), http.StatusUnauthorized)
		} else {
			_ = helpers.ErrorJSON(w, errors.New("this username has already been used"), http.StatusUnauthorized)
		}
		return
	}

	payload := helpers.JsonResponse{
		Error:   false,
		Message: userResponse.Message,
		Data:    userResponse.Data,
	}

	_ = helpers.WriteJSON(w, http.StatusAccepted, payload)
}

// LoginUser is an API endpoint that register a user and returns a JSON response.
// @Tags Auth
// @Summary Logs a user
// @Description Logs in a user with the given data.
// @Accept json
// @Produce json
// @Param requestPayload body requestPayload true "User data"
// @Success 202 {object} jsonResponse "Successful registration"
// @Failure 401 {object} jsonResponse "Invalid credentials"
// @Router /auth/registration [post]
func (app *brokerHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}
