package question

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/stclaird/quizzie.cloud/questionCreator/pkg/models"
	"github.com/stclaird/quizzie.cloud/questionCreator/pkg/util"
	"google.golang.org/api/option"
)

func createPrompt(questionIn models.QuestionIn) string {
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
    promptAnswers := fmt.Sprintf("give me %v correct answers and %v incorrect answers using this JSON schema:", numCorrectAns,numInCorrectAns)
    promptJson := fmt.Sprintf("Questions = {'questionText': string, 'answerReference': string, 'answers':[ %s ]}", answersStr)
    promptSuffix := "Return: <Question>"

    fullQuestion := fmt.Sprintf("%s %s %s %s %s", promptPrefix, questionIn.QuestionText, promptAnswers, promptJson, promptSuffix )

    return fullQuestion
}

func askAi (questionIn models.QuestionIn) models.Questions {
	config,err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config")
	}

	//Sent prompt to AI API
	ctx := context.Background()
	// Access your API key as an environment variable (see "Set up your API key" above)
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	model := client.GenerativeModel(config.Geminimodel)
	// Ask the model to respond with JSON.
	model.ResponseMIMEType = "application/json"
	prompt :=  createPrompt(questionIn)
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
	}

	json.Unmarshal(bytes, &questions)

	return questions

}
