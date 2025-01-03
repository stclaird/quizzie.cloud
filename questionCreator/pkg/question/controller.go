// Package questions
// Provide api functionailty for question objects
package question

import (
	"github.com/gin-gonic/gin"
	"github.com/stclaird/quizzie.cloud/questionCreator/pkg/util"
)

func RegisterRoutes(router *gin.Engine, config util.Config) {
    questionRoutes := router.Group("/questions")
    questionRoutes.POST("/", createQuestion)
    questionRoutes.GET("/health", Health)
}
