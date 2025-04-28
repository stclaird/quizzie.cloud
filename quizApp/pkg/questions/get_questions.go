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

func (h handler) GetQuestionsCatSubCat(ctx *gin.Context) {
	//get questions that match cat and subcat
	catSubCat := ctx.Param("catsubcat")

	cat, subcat := splitCatSubcat(catSubCat)

	var questions []models.Question
	var questionsNoAnswers []models.QuestionNoCorrectAnswer

	if result := h.DB.Model(&models.Question{}).Preload("Answers").Where("category = ? AND subcategory = ?", cat, subcat).Find(&questions); result.Error != nil {
		ctx.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	questionsNoAnswers = removeCorrectAnswerfield(questions)

	ctx.JSON(http.StatusOK, &questionsNoAnswers)

}
