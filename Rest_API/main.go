package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rest-api/database"
	"rest-api/models"
)

func main() {
	database.InitDB()
	// Initialize HTTP server with basic features
	server := gin.Default()

	// Handling a GET request
	server.GET("/events", getEvents)
	server.POST("/events", createEvent)

	server.Run("localhost:8080")

}

func getEvents(context *gin.Context) {
	// Sending back a response in JSON format
	// gin.H{} is alias of map[string]any
	//context.JSON(http.StatusOK, gin.H{"message": "Hello!"})

	events, err := models.GetAllEvents()
	if err != nil {
		context.JSON(http.StatusInternalServerError,
			gin.H{"message": "Could not fetch events. "})
	}
	context.JSON(http.StatusOK, events)

}

func createEvent(context *gin.Context) {
	var event models.Event
	// Converting request to JSON data
	err := context.ShouldBindJSON(&event)
	if err != nil {
		context.JSON(http.StatusBadRequest,
			gin.H{"message": "Could not parse request data."})
		return
	}

	event.ID = 1
	event.UserID = 1
	err = event.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError,
			gin.H{"message": "Could not create event. "})
	}
	context.JSON(http.StatusCreated,
		gin.H{"message": "Event created!", "event": event})
}
