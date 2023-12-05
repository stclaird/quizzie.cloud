package questions

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stclaird/quizzie.cloud/pkg/common/models"
)

func (h handler) GetQuestion(ctx *gin.Context) {
    id := ctx.Param("id")

    var book models.Book

    if result := h.DB.First(&book, id); result.Error != nil {
        ctx.AbortWithError(http.StatusNotFound, result.Error)
        return
    }

    ctx.JSON(http.StatusOK, &book)
}