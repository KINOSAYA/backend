package routes

import (
	"broker-service/internal/handlers"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

type Router interface {
	GetRoutes() http.Handler
}

type ChiRouters struct {
	Handler handlers.Handler
	WebPort string
}

func NewChiRouter(newHandler handlers.Handler, webPort string) *ChiRouters {
	return &ChiRouters{
		Handler: newHandler,
		WebPort: webPort,
	}
}

func (chiRouter *ChiRouters) GetRoutes() http.Handler {
	mux := chi.NewRouter()

	// specify who is allowed to connect
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Access-Control-Allow"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Use(middleware.Heartbeat("/ping"))

	mux.Get("/", chiRouter.Handler.Broker)
	mux.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://localhost:%s/swagger/doc.json", chiRouter.WebPort)), //The url pointing to API definition
	))

	mux.Route("/auth", func(mux chi.Router) {
		mux.Post("/registration", chiRouter.Handler.AuthRegisterUser)
		mux.Post("/login", chiRouter.Handler.AuthLoginUser)
	})

	return mux
}
