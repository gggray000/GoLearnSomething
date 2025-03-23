package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rest-api/util"
)

func Authenticate(context *gin.Context) {
	token := context.Request.Header.Get("authorization")

	if token == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized,
			gin.H{"message": "Not authorized."})
		return
	}

	userId, err := util.VerifyJwt(token)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized,
			gin.H{"message": "Invalid token."})
		return
	}

	context.Set("userId", userId)
	context.Next()
}
