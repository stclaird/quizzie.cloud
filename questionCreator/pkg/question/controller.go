// Package questions
// Provide api functionailty for question objects
package question

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
    questionRoutes := router.Group("/questions")
    questionRoutes.POST("/", generateQuestion)
    questionRoutes.GET("/health", Health)
}
