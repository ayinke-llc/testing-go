package server

import (
	testinggo "ayinke-llc/gophercrunch/testing-go"
	"ayinke-llc/gophercrunch/testing-go/cmd/config"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/middleware"
	chi "github.com/go-chi/chi/v5"
	"github.com/rs/cors"
)

func New(cfg config.Config, store testinggo.Store, httpPort int) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", httpPort),
		Handler: buildRoutes(cfg, store),
	}
}

func buildRoutes(cfg config.Config, store testinggo.Store) http.Handler {

	router := chi.NewRouter()

	router.Use(middleware.AllowContentType("application/json"))
	router.Use(middleware.RequestID)
	router.Use(jsonResponse)

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status" :"up}`))
	})

	router.Post("/items", createTaskItem(store))
	router.Get("/items/{id}", fetchTaskItem(store))

	return cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	}).Handler(router)
}
