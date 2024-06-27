package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/IBM/sarama"
)

// Security Headers for Router
func (app *application) securityHeadersMiddleware(next http.Handler) http.Handler {
	// TODO: Set the correct values
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; img-src 'self' https://frontend.com; media-src 'self' https://frontend.com; form-action 'self' https://frontend.com")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		next.ServeHTTP(w, r)
	})
}

// Auth Middleware
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

// Uploaded File Chunk Context Middleware
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

// AccessLogMiddleware
type AccessLogEntry struct {
	Method       string  `json:"method"`
	Host         string  `json:"host"`
	Path         string  `json:"path"`
	IP           string  `json:"ip"`
	ResponseTime float64 `json:"response_time"`

	encoded []byte
	err     error
}

func (ale *AccessLogEntry) ensureEncoded() {
	if ale.encoded == nil && ale.err == nil {
		ale.encoded, ale.err = json.Marshal(ale)
	}
}

func (ale *AccessLogEntry) Length() int {
	ale.ensureEncoded()
	return len(ale.encoded)
}

func (ale *AccessLogEntry) Encode() ([]byte, error) {
	ale.ensureEncoded()
	return ale.encoded, ale.err
}

func (app *application) withAccessLogs(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()

		next.ServeHTTP(w, r)

		entry := &AccessLogEntry{
			Method:       r.Method,
			Host:         r.Host,
			Path:         r.RequestURI,
			IP:           r.RemoteAddr,
			ResponseTime: float64(time.Since(started)) / float64(time.Second),
		}

		// Using Client IP Address as Key to be stored on same partition
		app.logproducer.Input() <- &sarama.ProducerMessage{
			Topic: "access_log",
			Key:   sarama.StringEncoder(r.RemoteAddr),
			Value: entry,
		}
	})
}
