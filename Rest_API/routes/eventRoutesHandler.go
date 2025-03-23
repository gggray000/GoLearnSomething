package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rest-api/models"
	"strconv"
)

func getEvents(context *gin.Context) {
	// Sending back a response in JSON format
	// gin.H{} is alias of map[string]any
	//context.JSON(http.StatusOK, gin.H{"message": "Hello!"})

	events, err := models.GetAllEvents()
	if err != nil {
		context.JSON(http.StatusInternalServerError,
			gin.H{"message": "Could not fetch events."})
	}
	context.JSON(http.StatusOK, events)

}

func getEvent(context *gin.Context) {
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest,
			gin.H{"message": "Could not parse event ID."})
		return
	}
	event, err := models.GetEventById(eventId)
	if err != nil {
		context.JSON(http.StatusInternalServerError,
			gin.H{"message": "Could not fetch requested events."})
	}

	context.JSON(http.StatusOK, event)
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
	userId := context.GetInt64("userId")
	event.UserID = userId

	err = event.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError,
			gin.H{"message": "Could not create event."})
	}
	context.JSON(http.StatusCreated,
		gin.H{"message": "Event created!", "event": event})
}

func updateEvent(context *gin.Context) {
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest,
			gin.H{"message": "Could not parse event ID."})
		return
	}

	userId := context.GetInt64("userId")
	event, err := models.GetEventById(eventId)
	if err != nil {
		context.JSON(http.StatusInternalServerError,
			gin.H{"message": "Could not fetch requested event."})
		return
	}

	if event.UserID != userId {
		context.JSON(http.StatusUnauthorized,
			gin.H{"message": "Unauthorized to update event."})
		return
	}

	var updatedEvent models.Event
	err = context.ShouldBindJSON(&updatedEvent)
	if err != nil {
		context.JSON(http.StatusBadRequest,
			gin.H{"message": "Could not parse request data."})
		return
	}
	updatedEvent.ID = eventId
	err = updatedEvent.Update()
	if err != nil {
		context.JSON(http.StatusInternalServerError,
			gin.H{"message": "Could not fetch requested event."})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Event updated successfully."})
}

func deleteEvent(context *gin.Context) {
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest,
			gin.H{"message": "Could not parse event ID."})
		return
	}

	userId := context.GetInt64("userId")
	event, err := models.GetEventById(eventId)
	if err != nil {
		context.JSON(http.StatusInternalServerError,
			gin.H{"message": "Could not fetch requested event."})
		return
	}
	if event.UserID != userId {
		context.JSON(http.StatusUnauthorized,
			gin.H{"message": "Unauthorized to delete event."})
		return
	}

	err = event.Delete()
	if err != nil {
		context.JSON(http.StatusInternalServerError,
			gin.H{"message": "Could not delete event."})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully."})
}
