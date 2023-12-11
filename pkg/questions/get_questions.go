package questions

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stclaird/quizzie.cloud/pkg/common/models"
)

func (h handler) GetQuestions(ctx *gin.Context) {
    var questions []models.Question

    if result := h.DB.Model(&models.Question{}).Preload("Answers").Find(&questions); result.Error != nil {
        ctx.AbortWithError(http.StatusNotFound, result.Error)
        return
    }

    ctx.JSON(http.StatusOK, &questions)
}
