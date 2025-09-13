package questions

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stclaird/quizzie.cloud/pkg/common/models"
)

func (h handler) GetQuestions(ctx *gin.Context) {
	var questions []models.Question
	var questionsNoAnswers []models.QuestionNoCorrectAnswer

	if result := h.DB.Model(&models.Question{}).Preload("Answers").Find(&questions); result.Error != nil {
		ctx.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	questionsNoAnswers = removeCorrectAnswerfield(questions)

	ctx.JSON(http.StatusOK, &questionsNoAnswers)
}

func (h handler) GetQuestionsByCategoryAndSubcategory(ctx *gin.Context) {
	category := ctx.Param("category")
	subcategory := ctx.Param("subcategory")

	var questions []models.Question
	if err := h.DB.Where("category = ? AND subcategory = ?", category, subcategory).Preload("Answers").Find(&questions).Error; err != nil {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}

	questionsNoAnswers := removeCorrectAnswerfield(questions)
	ctx.JSON(http.StatusOK, &questionsNoAnswers)
}
