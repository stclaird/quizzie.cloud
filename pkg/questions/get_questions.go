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

    //Database retrieves objects including their correct answers - we need to copy this to an object with no correct answers
    //before sending it to the user. otherwise they can cheat!
    for _, question := range questions { //Copying by loop
        var questionnoanswer models.QuestionNoCorrectAnswer
        questionnoanswer.ID =  question.ID
        questionnoanswer.Category = question.Category
        questionnoanswer.Subcategory = question.Subcategory
        questionnoanswer.Text = question.Text
        questionnoanswer.Type = question.Type
        for _, answer := range question.Answers{
            a := struct{
                ID    uint
                Text   string `json:"text"`
            }{
                ID:    answer.ID,
                Text:  answer.Text,
            }
            questionnoanswer.Answers = append(questionnoanswer.Answers, a )
        }
        questionsNoAnswers = append(questionsNoAnswers, questionnoanswer)
     }

    ctx.JSON(http.StatusOK, &questionsNoAnswers)
}
