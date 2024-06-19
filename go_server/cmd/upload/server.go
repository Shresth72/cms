package main

import (
	"fmt"
	"log"
	"net/http"
)

// Add TLS
func (app *application) Serve() error {
	server := &http.Server{
		Addr:     fmt.Sprintf(":%d", app.cfg.port),
		Handler:  app.routes(),
		ErrorLog: log.New(app.logger, "", 0),
	}

	return server.ListenAndServe()
}
