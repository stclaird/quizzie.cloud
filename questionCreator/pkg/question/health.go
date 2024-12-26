package question

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Health(ctx *gin.Context) {

	response := gin.H{
		"status" : "ok",
	}

	ctx.JSON(http.StatusOK, response)
}
