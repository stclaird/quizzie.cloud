package questions

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

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

func (h handler) Answers(ctx *gin.Context) {
    id := ctx.Param("id")
    submittedAnswer := ctx.Param("answer")

    var question models.Question

    if result := h.DB.Preload("Answers").First(&question, id); result.Error != nil {
        ctx.AbortWithError(http.StatusNotFound, result.Error)
        return
    }

    isCorrect, CorrectAnswer := checkAnswer(question, submittedAnswer)

	response := models.AnswerResponse{
		IsCorrect : isCorrect,
		CorrectAnswer: CorrectAnswer,
	}

    ctx.JSON(http.StatusOK, response)
}

//compare the real answer with the user submitted answer
func checkAnswer(question models.Question, submittedAnswer string) (bool, []models.Answer) {
    var answers []string // store the answers from the question so we can compare with submitted answer
	var correctAnswersResp []models.Answer //this is sent back to the user

	for _,v := range question.Answers{
		if v.IsCorrect == true {
			fmt.Printf("Correct AnswerID %v \n", v.ID)
			correctAnswer := models.Answer{
                ID :  v.ID,
                Text : v.Text,
				IsCorrect:  true,
			}
			correctAnswersResp = append(correctAnswersResp, correctAnswer)
			answers = append(answers, strconv.Itoa(int(v.ID)))
		}
	}

	submittedAnswer = SortString(submittedAnswer) //string
	var answersStr string

	answersStr = strings.Join(answers,"")
	answersStr = SortString(answersStr)

	fmt.Printf("Final: answersStr: %s, submittedAnswer: %s", answersStr, submittedAnswer)

	if answersStr != submittedAnswer {
		return false, correctAnswersResp
	}

	return true, correctAnswersResp
}