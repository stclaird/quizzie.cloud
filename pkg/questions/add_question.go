package questions

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stclaird/quizzie.cloud/pkg/common/models"
)

type AddQuestionRequestBody struct {
    Title       string `json:"title"`
    Author      string `json:"author"`
    Description string `json:"description"`
}


func (h handler) AddQuestion(ctx *gin.Context) {
    body := AddQuestionRequestBody{}

    if err := ctx.BindJSON(&body); err != nil {
        ctx.AbortWithError(http.StatusBadRequest, err)
        return
    }

    var question models.Question

    question.Title = body.Title
    question.Author = body.Author
    question.Description = body.Description

    if result := h.DB.Create(&question); result.Error != nil {
        ctx.AbortWithError(http.StatusNotFound, result.Error)
        return
    }

    ctx.JSON(http.StatusCreated, &question)
}
