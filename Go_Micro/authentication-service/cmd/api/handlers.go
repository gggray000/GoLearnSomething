package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJson(w, r, &requestPayload)
	if err != nil {
		log.Println(err)
		app.writeJsonError(w, err, http.StatusBadRequest)
		return
	}

	// Validate user
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		log.Println("some message", err)
		app.writeJsonError(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		log.Println("some message", err)
		app.writeJsonError(w, errors.New("invalid credentials"), http.StatusBadRequest)
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged on user: %s", user.Email),
		Data:    user,
	}

	app.writeJson(w, http.StatusAccepted, payload)
}
