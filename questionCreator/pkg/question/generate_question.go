package question

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"github.com/stclaird/questionCreator/questionCreator/pkg/models"
	"google.golang.org/api/option"
)

func generateQuestion(ctx *gin.Context) {
	//Generate a question return as API response and write copy to file system
    var questionDir = "questions"
    var questionIn models.QuestionIn
    if err := ctx.BindJSON(&questionIn); err != nil {
        ctx.JSON(400, gin.H{"error": err.Error()})
        return
    }

    questionsOut := AskAi(questionIn)

	for _, question := range questionsOut.Questions {
		questionJsonBytes, _ := json.Marshal(question)
		questionName := generateQuestionFileName(question)
		err := os.MkdirAll(questionDir, 0755)

		if err != nil {
			panic(err)
		}
		fmt.Printf("Writing File %s\n", questionName)
		var writeFileName = fmt.Sprintf("%s/%s", questionDir, questionName)
		os.WriteFile(writeFileName, questionJsonBytes, os.ModePerm)
	}

    ctx.JSON(http.StatusOK, questionsOut)
}

func AskAi (questionIn models.QuestionIn) models.Questions {

	ctx := context.Background()
	// Access your API key as an environment variable (see "Set up your API key" above)
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	model := client.GenerativeModel("gemini-1.5-pro-latest")
	// Ask the model to respond with JSON.
	model.ResponseMIMEType = "application/json"
	prompt :=  generatePrompt(questionIn)
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		log.Fatal(err)
	}

	var bytes []byte
	var questions models.Questions
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			bytes = []byte(txt)
		}
	//	fmt.Println(string(bytes))
	}

	json.Unmarshal(bytes, &questions)

	return questions

}

func generatePrompt(questionIn models.QuestionIn) string {
    // create the prompt to make use of AI's JSONs RESPONSE
    var answersStr []string

	var numAns int

	numCorrectAns := questionIn.NumCorrectAns
	numInCorrectAns := questionIn.NumInCorrectAns

    if numCorrectAns == 0 {
        numCorrectAns = 1
    }

	if numInCorrectAns == 0 {
        numInCorrectAns = 1
    }

    //How many answers do we want in our json template
    numAns = numCorrectAns + numInCorrectAns


    for i := 0; i < numAns; i++ {
        answersStr = append(answersStr, "{'text': string, 'iscorrect':bool}" )
    }


    promptPrefix := fmt.Sprintf("Ask me %v questions regarding", questionIn.NumQuestions)
    promptAnswers := fmt.Sprintf("give me %v correct answers and %v incorrect answers using this JSON schema:", questionIn.NumCorrectAns, questionIn.NumInCorrectAns)
    promptJson := fmt.Sprintf("Questions = {'question': string, 'answerReference': string, 'answers':[ %s ]}", answersStr)
    promptSuffix := "Return: <Question>"


    fullQuestion := fmt.Sprintf("%s %s %s %s %s", promptPrefix, questionIn.QuestionText, promptAnswers, promptJson, promptSuffix )

    return fullQuestion
}
