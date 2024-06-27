package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *application) Serve() error {
  server := &http.Server{
    Addr: fmt.Sprintf(":%d", app.cfg.port),
    Handler: app.routes(),
    ErrorLog: log.New(app.logger, "", 0),
  }

  return server.ListenAndServe()
}

func (app *application) routes() http.Handler {
	router := chi.NewRouter()

	// Middlewares
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

  router.Use(cors.Handler(cors.Options{
    AllowedOrigins: []string{"http://*", "https://*"},
    AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
		AllowCredentials: false,            // prevent CSRF
		MaxAge:           300,
  }))

  router.Get("/home", app.getAllVideos)
  router.Get("/video", app.getVideo)

  return router
}
