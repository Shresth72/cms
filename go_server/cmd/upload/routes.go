package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *application) routes() http.Handler {
	router := chi.NewRouter()

	// Middlewares
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(app.securityHeadersMiddleware)

	// Cors
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://front-end-domain.com"},
		AllowedMethods:   []string{"POST", "GET"},
		AllowedHeaders:   []string{"Content-Type", "X-OAuth-Code", "Authorization"},
		ExposedHeaders:   []string{"Etag"}, // If Etag sent through Headers
		AllowCredentials: false,            // prevent CSRF
		MaxAge:           300,
	}))

	// Multipart Routes
	router.With(app.Auth).Post("/multipart/init", app.initMultipartUpload)
	router.With(app.Auth).Post("/multipart", app.uploadChunkMiddlware(app.uploadChunk))
	router.With(app.Auth).Post("/multipart/complete", app.completeMultipart)

	// OAuth Routes
	router.Get("/auth/google/callback", app.oauthGoogleCallback)

	// Kafka Routes
	router.With(app.withAccessLogs).Post("/kafka/send", app.sendMessageToKafka)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("SMOSH! SHUT UP"))
	})

	return router
}
