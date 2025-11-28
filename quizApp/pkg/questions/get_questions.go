package questions

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stclaird/quizzie.cloud/pkg/common/models"
	"go.uber.org/zap"
)

func (h handler) GetQuestions(ctx *gin.Context) {
	zap.L().Debug("GetQuestions endpoint called")

	var questions []models.Question
	var questionsNoAnswers []models.QuestionNoCorrectAnswer

	if result := h.DB.Model(&models.Question{}).Preload("Answers").Find(&questions); result.Error != nil {
		zap.L().Error("Failed to fetch questions from database", zap.Error(result.Error))
		ctx.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	zap.L().Info("Successfully fetched questions", zap.Int("count", len(questions)))
	questionsNoAnswers = removeCorrectAnswerfield(questions)

	ctx.JSON(http.StatusOK, &questionsNoAnswers)
}

func (h handler) GetQuestionsByCategoryAndSubcategory(ctx *gin.Context) {
	category := ctx.Param("category")
	subcategory := ctx.Param("subcategory")

	zap.L().Debug("GetQuestionsByCategoryAndSubcategory endpoint called",
		zap.String("category", category),
		zap.String("subcategory", subcategory))

	var questions []models.Question
	if err := h.DB.Where("category = ? AND subcategory = ?", category, subcategory).Preload("Answers").Find(&questions).Error; err != nil {
		zap.L().Error("Failed to fetch questions by category and subcategory",
			zap.String("category", category),
			zap.String("subcategory", subcategory),
			zap.Error(err))
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}

	zap.L().Info("Successfully fetched questions by category and subcategory",
		zap.String("category", category),
		zap.String("subcategory", subcategory),
		zap.Int("count", len(questions)))

	questionsNoAnswers := removeCorrectAnswerfield(questions)
	ctx.JSON(http.StatusOK, &questionsNoAnswers)
}
