package main

import (
	"bytes"
	"context"
	"io"
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
		AllowedOrigins: []string{"https://front-end-domain.com"}, 
    AllowedMethods: []string{"POST", "GET"},
    AllowedHeaders: []string{"Content-Type", "X-OAuth-Code", "Authorization"},
    ExposedHeaders: []string{"Etag"}, // If Etag sent through Headers
    AllowCredentials: false, // prevent CSRF
    MaxAge: 300,
	}))

	// Multipart Routes
	router.With(app.Auth).Post("/multipart/init", app.initMultipartUpload)
	router.With(app.Auth).Post("/multipart", app.uploadChunkMiddlware(app.uploadChunk))
	router.With(app.Auth).Post("/multipart/complete", app.completeMultipart)
	// router.With(app.Auth).Post("/multipart/uploadDb", app.uploadToDb)

  // OAuth 
  router.Get("/auth/google/callback", app.oauthGoogleCallback)
  // router.Get("/auth/refreshtoken", app.getNewAccessToken)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("SMOSH! SHUT UP"))
	})

	return router
}

// Middlewares
func (app *application) Auth(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    tokenString := r.Header.Get("Authorization")
    if tokenString == "" {
      app.errorResponse(w, r, http.StatusUnauthorized, "Invalid Token")
      return
    }
 
    userInfo, err := app.VerifyJWT(tokenString)
    if err != nil {
      app.errorResponse(w, r, http.StatusUnauthorized, "Invalid Token")
      return
    }

    ctx := context.WithValue(r.Context(), "userInfo", userInfo)
    next.ServeHTTP(w, r.WithContext(ctx))
  })
}

func (app *application) uploadChunkMiddlware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(int64(maxChunkSize))
		if err != nil {
			app.serverError(w, r, err, "error parsing form data")
			return
		}

		// form key - chunk
		file, _, err := r.FormFile("chunk")
		if err != nil {
			app.notFoundError(w, r)
			return
		}
		defer file.Close()

		var fileBytes bytes.Buffer
		if _, err = io.Copy(&fileBytes, file); err != nil {
			app.serverError(w, r, err, "error reading file")
			return
		}

		ctx := context.WithValue(r.Context(), "fileBytes", fileBytes.Bytes())
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (app *application) securityHeadersMiddleware(next http.Handler) http.Handler {
  // TODO: Set the correct values
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Security-Policy", "default-src 'self'; img-src 'self' https://frontend.com; media-src 'self' https://frontend.com; form-action 'self' https://frontend.com")
    w.Header().Set("X-Content-Type-Options", "nosniff")
    w.Header().Set("X-Frame-Options", "DENY")
    next.ServeHTTP(w, r)
  })
}
