package questions

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stclaird/quizzie.cloud/pkg/common/models"
)

func (h handler) GetQuestion(ctx *gin.Context) {
    id := ctx.Param("id")

    var question models.Question

    if result := h.DB.First(&question, id); result.Error != nil {
        ctx.AbortWithError(http.StatusNotFound, result.Error)
        return
    }

    ctx.JSON(http.StatusOK, &question)
}