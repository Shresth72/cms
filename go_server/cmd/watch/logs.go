package main

import (
	"fmt"
	"net/http"

	"github.com/Shresth72/server/internals/utils"
)

// Base Logs
func (app *application) logMsgFmt(msg string, args ...interface{}) {
	app.logger.Info().Msgf(msg, args...)
}

func (app *application) logSuccess(_ *http.Request, msg string) {
	app.logger.Info().Msg(msg)
}

// Error Logs
func (app *application) logError(_ *http.Request, err error) {
	app.logger.Error().Err(err).Msg("")
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, data any) {
	msg := utils.Envelope{
		"error": data,
	}

	err := utils.WriteJSON(w, r, status, msg)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error, msg ...string) {
	app.logError(r, err)
	var message string
	if len(msg) > 0 {
		message = msg[0]
	} else {
		message = "something went wrong"
	}

	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err ...error) {
  var erresp []error
  if len(err) > 0 {
    erresp = err
  }
	app.errorResponse(w, r, http.StatusBadRequest, fmt.Sprint(erresp))
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request) {
	msg := `the requested resource cannot be found`
	app.errorResponse(w, r, http.StatusNotFound, msg)
}

func (app *application) unproccessableEntityError(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (app *application) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	app.errorResponse(w, r, http.StatusConflict, message)
}
