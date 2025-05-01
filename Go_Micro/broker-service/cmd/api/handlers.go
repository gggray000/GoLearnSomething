package main

import (
	"broker/event"
	"broker/logs"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"net/rpc"
	"time"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayLoad  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayLoad struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJson(w, http.StatusOK, payload)

	/*out, _ := json.MarshalIndent(payload, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(out)*/
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJson(w, r, &requestPayload)
	if err != nil {
		app.writeJsonError(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	case "log":
		//app.logItem(w, requestPayload.Log)
		//app.logEventViaRabbit(w, requestPayload.Log)
		app.logItemViaRPC(w, requestPayload.Log)
	case "mail":
		app.sendMail(w, requestPayload.Mail)
	default:
		app.writeJsonError(w, errors.New("unknown action"))
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	// create json that will be sent to auth service
	jsonData, _ := json.MarshalIndent(a, "", "\t")
	// call the service
	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.writeJsonError(w, err)
		return
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		app.writeJsonError(w, err)
		return
	}
	defer response.Body.Close()

	// get back correct status code
	if response.StatusCode == http.StatusUnauthorized {
		app.writeJsonError(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.writeJsonError(w, errors.New("error calling auth service"))
		return
	}

	// create a variable for reading response.Body
	var jsonFromService jsonResponse

	// decode json from auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.writeJsonError(w, err)
		return
	}

	if jsonFromService.Error {
		app.writeJsonError(w, err, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated successfully"
	payload.Data = jsonFromService.Data

	app.writeJson(w, http.StatusAccepted, payload)

	authLog := fmt.Sprintf("Rabbit: User %v has logged in.", payload.Data.(map[string]interface{})["email"])
	app.pushToQueue("auth", authLog)
}

func (app *Config) logItem(w http.ResponseWriter, entry LogPayLoad) {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceURL := "http://logger-service/log"
	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))

	if err != nil {
		app.writeJsonError(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		app.writeJsonError(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.writeJsonError(w, errors.New("error calling logger service"))
		return
	}

	var jsonFromService jsonResponse
	jsonFromService.Error = false
	jsonFromService.Message = "Logged successfully"
	app.writeJson(w, http.StatusAccepted, jsonFromService)
}

// In reality, broker service should not talk to mail service directly, because malicious users can send JSON payload to send spams.
// When sending email is necessary, the request should be routed to auth service, which will talk to mail service.
func (app *Config) sendMail(w http.ResponseWriter, msg MailPayload) {
	jsonData, _ := json.MarshalIndent(msg, "", "\t")

	mailServiceURL := "http://mail-service/send"
	request, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))

	if err != nil {
		app.writeJsonError(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		app.writeJsonError(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.writeJsonError(w, errors.New("error calling mail service"))
		return
	}

	var jsonFromService jsonResponse
	jsonFromService.Error = false
	jsonFromService.Message = "Mail sent successfully to " + msg.To
	app.writeJson(w, http.StatusAccepted, jsonFromService)
}

// Replacement for logItem(), it calls logger-service by pushing the payload to RabbitMQ
func (app *Config) logEventViaRabbit(w http.ResponseWriter, l LogPayLoad) {
	err := app.pushToQueue(l.Name, l.Data)
	if err != nil {
		app.writeJsonError(w, err)
		return
	}

	var jsonFromService jsonResponse
	jsonFromService.Error = false
	jsonFromService.Message = "Logged via RabbitMQ successfully"
	app.writeJson(w, http.StatusAccepted, jsonFromService)
}

func (app *Config) pushToQueue(name, msg string) error {
	emitter, err := event.NewEventEmitter(app.Rabbit)
	if err != nil {
		return err
	}

	payload := LogPayLoad{
		Name: name,
		Data: msg,
	}

	j, _ := json.MarshalIndent(&payload, "", "\t")
	err = emitter.Push(string(j), "log.INFO")
	if err != nil {
		return err
	}

	return nil
}

type RPCPayload struct {
	Name string
	Data string
}

func (app *Config) logItemViaRPC(w http.ResponseWriter, l LogPayLoad) {
	client, err := rpc.Dial("tcp", "logger-service:5001")
	if err != nil {
		app.writeJsonError(w, err)
		return
	}

	rpcPayload := RPCPayload{
		Name: l.Name,
		Data: l.Data,
	}

	var result string
	err = client.Call("RPCServer.LogInfo", rpcPayload, &result)
	if err != nil {
		app.writeJsonError(w, err)
		return
	}

	jsonResultFromRPCServer := jsonResponse{
		Error:   false,
		Message: result,
	}
	app.writeJson(w, http.StatusAccepted, jsonResultFromRPCServer)
}

func (app *Config) logViaGRPC(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJson(w, r, &requestPayload)
	if err != nil {
		app.writeJsonError(w, err)
		return
	}

	conn, err := grpc.Dial(
		"logger-service:50001",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		app.writeJsonError(w, err)
		return
	}
	defer conn.Close()

	c := logs.NewLogServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = c.WriteLog(ctx, &logs.LogRequest{
		LogEntry: &logs.Log{
			Name: requestPayload.Log.Name,
			Data: requestPayload.Log.Data,
		},
	})
	if err != nil {
		app.writeJsonError(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Logged successfully"

	app.writeJson(w, http.StatusAccepted, payload)
}
