//Package questions
// Provide api functionailty for question objects

package questions

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type handler struct {
	DB *gorm.DB
}

func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
	h := &handler{
		DB: db,
	}

	questionRoutes := router.Group("/questions")
	questionRoutes.GET("/", h.GetQuestions)
	questionRoutes.GET("/catsubcat/:catsubcat", h.GetQuestionsCatSubCat)
	questionRoutes.GET("/id/:id", h.GetQuestion)
	questionRoutes.GET("/answer/:id/:answer", h.Answers)

	categoryRoutes := router.Group("/categories")
	categoryRoutes.GET("/", h.GetCategories)
}
