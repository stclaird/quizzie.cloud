package question

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Health(ctx *gin.Context) {
	//Return Health Status Response
	response := gin.H{
		"status" : "ok",
	}

	ctx.JSON(http.StatusOK, response)
}
