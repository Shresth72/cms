package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Shresth72/server/internals/utils"
	"github.com/go-redis/redis/v8"
)

func (app *application) oauthGoogleCallback(w http.ResponseWriter, r *http.Request) {
	// code := r.URL.Query().Get("code")
	code := r.Header.Get("X-OAuth-Code")
	if code == "" {
		app.badRequestError(w, r)
		return
	}

	token, err := app.oauth.Exchange(context.Background(), code)
	if err != nil {
		app.serverError(w, r, err, "failed to exchange code for token")
		return
	}

	// Fetch user info from google
	userInfo, err := GetUserInfo(token.AccessToken)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Self Signed JWT, cannot send AccessToken to the Client
	signedToken, id, err := SignJWT(userInfo)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Save Refresh Token (FIXME)
	err = app.saveRefreshToken(id, token.RefreshToken)
	if err != nil {
		app.serverError(w, r, err, "failed to save refresh token")
		return
	}

	err = utils.WriteJSON(w, r, http.StatusCreated, utils.Envelope{
		"signed_token": signedToken,
	})
	if err != nil {
		app.serverError(w, r, err)
		return
	}
}

type getNewAccessToken struct {
	UserId string `json:"user_id"`
}

// Currently not implementing for security reasons
func (app *application) getNewAccessToken(w http.ResponseWriter, r *http.Request) {
	var body getNewAccessToken
	err := utils.ReadJSON(w, r, &body)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	userId := body.UserId
	key := fmt.Sprintf("refresh_token_%s", userId)

	refreshToken, err := app.redis.Get(context.Background(), key).Result()
	if err == redis.Nil {
		app.notFoundError(w, r)
		return
	} else if err != nil {
		app.serverError(w, r, err, "failed to get refresh token from redis")
		return
	}

	newToken, err := app.refreshToken(refreshToken)
	if err != nil {
		app.serverError(w, r, err, "failed to get refresh token")
		return
	}

	if newToken.RefreshToken != "" && newToken.RefreshToken != refreshToken {
		err = app.saveRefreshToken(userId, newToken.RefreshToken)
		if err != nil {
			app.serverError(w, r, err, "failed to update refresh token")
			return
		}
}

	userInfo, err := GetUserInfo(newToken.AccessToken)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	signedToken, _, err := SignJWT(userInfo)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	err = utils.WriteJSON(w, r, http.StatusCreated, utils.Envelope{
		"signed_token": signedToken,
	})
	if err != nil {
		app.serverError(w, r, err)
		return
	}
}
