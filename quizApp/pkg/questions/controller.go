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

	//Question route group.
	questionRoutes := router.Group("/questions")

	//Question routes.
	questionRoutes.GET("/", h.GetQuestions)
	//Get questions by category and subcategory
	questionRoutes.GET("/:category/:subcategory", h.GetQuestionsByCategoryAndSubcategory)
	//Get question by id
	questionRoutes.GET("/id/:id", h.GetQuestion)
	//Submit an answer to a question
	questionRoutes.GET("/answer/:id/:answer", h.Answers)

	categoryRoutes := router.Group("/categories")
	categoryRoutes.GET("/", h.GetCategories)
}
