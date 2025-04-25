package main

import (
	"bytes"
	"encoding/json"
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

	// log authentication
	err = app.logRequest("authentication", fmt.Sprintf("%s logger in", user.Email))
	if err != nil {
		app.writeJsonError(w, err)
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged on user: %s", user.Email),
		Data:    user,
	}

	app.writeJson(w, http.StatusAccepted, payload)
}

func (app *Config) logRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil {
		return err
	}

	return nil
}
