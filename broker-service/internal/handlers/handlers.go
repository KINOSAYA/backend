package handlers

import (
	_ "broker-service/docs"
	"broker-service/event"
	"broker-service/internal/auth"
	"broker-service/internal/config"
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
	AuthUser(w http.ResponseWriter, r *http.Request)
	AuthLoginUser(w http.ResponseWriter, r *http.Request)
	ParseToken(w http.ResponseWriter, r *http.Request)
	Refresh(w http.ResponseWriter, r *http.Request)

	GetSlugsByCollection(w http.ResponseWriter, r *http.Request)
	GetFilmsBySlug(writer http.ResponseWriter, request *http.Request)
}

func NewBrokerHandler(authGrpcPort, authHost string, app config.Config) Handler {
	return &brokerHandler{
		AuthGrpcPort: authGrpcPort,
		AuthHost:     authHost,
		EventService: app.AmqpService,
	}
}

type brokerHandler struct {
	AuthGrpcPort string
	AuthHost     string
	EventService event.Service
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
	payload := jsonResponse{
		Error:   false,
		Message: "Successfully hit the broker",
	}
	_ = helpers.WriteJSON(w, http.StatusOK, payload)
}

type registrationPayload struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthUser is an API endpoint that register a user and returns a JSON response.
// @Tags Auth
// @Summary Register a new user
// @Description Registers a new user with the specified data.
// @Accept json
// @Produce json
// @Param registrationPayload body registrationPayload true "User data"
// @Success 202 {object} jsonResponse "Successful registration"
// @Failure 401 {object} jsonResponse "Invalid credentials"
// @Router /auth/registration [post]
func (app *brokerHandler) AuthUser(w http.ResponseWriter, r *http.Request) {
	var requestPayload registrationPayload

	err := helpers.ReadJSON(w, r, &requestPayload)
	if err != nil {
		_ = helpers.ErrorJSON(w, err)
		return
	}

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
		} else if strings.Contains(err.Error(), "this username has already been used") {
			_ = helpers.ErrorJSON(w, errors.New("this username has already been used"), http.StatusUnauthorized)
		} else {
			_ = helpers.ErrorJSON(w, err, http.StatusBadRequest)
		}
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: userResponse.Message,
		Data:    userResponse.Data,
	}

	_ = helpers.WriteJSON(w, http.StatusAccepted, payload)
}

// AuthLoginUser is an API endpoint that authenticate a user and returns a JSON response.
// @Tags Auth
// @Summary Logs a user
// @Description Logs in a user with the given data.
// @Accept json
// @Produce json
// @Param registrationPayload body registrationPayload true "User data"
// @Success 202 {object} jsonResponse "Successful registration"
// @Failure 401 {object} jsonResponse "Invalid credentials"
// @Router /auth/login [post]
func (app *brokerHandler) AuthLoginUser(w http.ResponseWriter, r *http.Request) {
	var requestPayload registrationPayload

	err := helpers.ReadJSON(w, r, &requestPayload)
	if err != nil {
		_ = helpers.ErrorJSON(w, err)
		return
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", app.AuthHost, app.AuthGrpcPort), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return
	}
	defer conn.Close()

	c := auth.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	userResponse, err := c.AuthUser(ctx, &auth.UserRequest{
		UserEntry: &auth.User{
			Username: requestPayload.Username,
			Email:    requestPayload.Email,
			Password: requestPayload.Password,
		},
	})
	if err != nil {
		_ = helpers.ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: userResponse.Message,
		Data:    userResponse.Data,
	}

	_ = helpers.WriteJSON(w, http.StatusAccepted, payload)
}

type tempTokenReqPayload struct {
	TokenString string `json:"token"`
}
type tempTokenResponsePayload struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

// ParseToken is an API endpoint that authenticate a user and returns a JSON response.
// @Tags Auth
// @Summary token parser
// @Description token validation (token is alive for 24 hours
// @Accept json
// @Produce json
// @Param registrationPayload body tempTokenReqPayload true "User data"
// @Success 202 {object} jsonResponse "Successful registration"
// @Failure 401 {object} jsonResponse "Invalid credentials"
// @Router /auth/parse-token [post]
func (app *brokerHandler) ParseToken(w http.ResponseWriter, r *http.Request) {
	var requestPayload tempTokenReqPayload

	err := helpers.ReadJSON(w, r, &requestPayload)
	if err != nil {
		_ = helpers.ErrorJSON(w, err)
		return
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", app.AuthHost, app.AuthGrpcPort), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return
	}
	defer conn.Close()

	c := auth.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	tokenResponse, err := c.CheckToken(ctx, &auth.TokenRequest{
		TokenString: requestPayload.TokenString,
	})
	if err != nil {
		_ = helpers.ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "checked Token",
		Data: tempTokenResponsePayload{
			ID:       int(tokenResponse.Id),
			Username: tokenResponse.Username,
		},
	}

	_ = helpers.WriteJSON(w, http.StatusAccepted, payload)
}

type token struct {
	Token string `json:"token"`
}

// Refresh is a handler that provides access token update
// @Tags Auth
// @Summary Refresh Access Token
// @Description Refreshes an access token using a refresh token.
// @ID refreshAccessToken
// @Accept json
// @Produce json
// @Param requestBody body token true "Refresh token structure"
// @Success 202 {object} jsonResponse "Token refresh successful"
// @Failure 400 {object} jsonResponse "Bad request"
// @Failure 401 {object} jsonResponse "Unauthorized"
// @Router /auth/refresh [post]
func (app *brokerHandler) Refresh(w http.ResponseWriter, r *http.Request) {

	requestPayload := token{}
	err := helpers.ReadJSON(w, r, &requestPayload)
	if err != nil {
		_ = helpers.ErrorJSON(w, err)
		return
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", app.AuthHost, app.AuthGrpcPort), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return
	}
	defer conn.Close()

	c := auth.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	tokenResponse, err := c.Refresh(ctx, &auth.TokenRequest{
		TokenString: requestPayload.Token,
	})
	if err != nil {
		_ = helpers.ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "updated access Token",
		Data: token{
			Token: tokenResponse.TokenString,
		},
	}

	_ = helpers.WriteJSON(w, http.StatusAccepted, payload)

}
