package question

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/stclaird/quizzie.cloud/questionCreator/pkg/models"
	"github.com/stclaird/quizzie.cloud/questionCreator/pkg/util"
)

func createQuestion(ctx *gin.Context) {
	//Generate a question return as API response and write copy to file system
	config := ctx.MustGet("config").(util.Config)

	var questionDir = config.Questionsdir
    var questionIn models.QuestionIn
    if err := ctx.BindJSON(&questionIn); err != nil {
        ctx.JSON(400, gin.H{"error": err.Error()})
        return
    }

	var questionCategory string
	var questionSubCategory string

	questionTextSplit := splitString(questionIn.QuestionText)

	if questionIn.Category == "" {
		questionCategory = questionTextSplit[0]
	} else {
		questionCategory = questionIn.Category
	}

	if questionIn.Subcategory == "" {
		questionSubCategory = questionTextSplit[1]
	} else {
		questionSubCategory = questionIn.Subcategory
	}

    questionsOut := askAi(questionIn)

	var questionsOutFile []models.QuestionOut

	for i := range questionsOut.Questions {

		questionOut := questionsOut.Questions[i]
        questionOut.Subcategory = questionSubCategory
		questionOut.Category = questionCategory
		questionOut.DateAdded = createDate()

		questionsOutFile = append(questionsOutFile, questionOut)

	}

	questionsOutFileBytes, _ := json.Marshal(questionsOutFile)

    questionName := generateQuestionFileName(questionCategory, questionSubCategory )
	err := os.MkdirAll(questionDir, 0755)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Writing File %s\n", questionName)
	var writeFileName = fmt.Sprintf("%s/%s", questionDir, questionName)
	os.WriteFile(writeFileName, questionsOutFileBytes, os.ModePerm)

    ctx.JSON(http.StatusOK, questionsOut)
}
