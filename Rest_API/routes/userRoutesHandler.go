package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rest-api/models"
	"rest-api/util"
)

func signup(context *gin.Context) {
	var user models.User
	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest,
			gin.H{"message": "Could not parse sign up data."})
		return
	}

	err = user.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError,
			gin.H{"message": "Could not save user."})
	}
	context.JSON(http.StatusCreated,
		gin.H{"message": "Sign up successfully!", "user": user})
}

func login(context *gin.Context) {
	var user models.User
	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest,
			gin.H{"message": "Could not parse login data."})
		return
	}

	err = user.ValidateCredentials()
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Authentication Failed."})
		return
	}

	token, err := util.GenerateJwt(user.Email, user.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError,
			gin.H{"message": "Could not generate JWT."})
	}

	context.JSON(http.StatusOK,
		gin.H{
			"message": "Logged in!",
			"user":    user,
			"token":   token,
		})
}
