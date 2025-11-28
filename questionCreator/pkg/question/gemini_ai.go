package question

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/stclaird/quizzie.cloud/questionCreator/pkg/models"
	"github.com/stclaird/quizzie.cloud/questionCreator/pkg/util"
	"go.uber.org/zap"
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
        answersStr = append(answersStr, "{'text': string, 'iscorrect': boolean}" )
    }

    promptPrefix := fmt.Sprintf("Generate %v questions about", questionIn.NumQuestions)
    promptAnswers := fmt.Sprintf("Each question should have %v correct answers and %v incorrect answers.", numCorrectAns, numInCorrectAns)
    promptJson := fmt.Sprintf("Return JSON in this exact format: {\"questions\": [{\"questionText\": string, \"category\": string, \"subcategory\": string, \"answerReference\": string, \"answers\": [%s]}]}. Make sure each question has appropriate category and subcategory based on its content.", strings.Join(answersStr, ", "))

    fullQuestion := fmt.Sprintf("%s %s. %s %s", promptPrefix, questionIn.QuestionText, promptAnswers, promptJson)

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
	zap.L().Debug("Sending prompt to AI", zap.String("prompt", prompt))
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		zap.L().Error("Error calling AI API", zap.Error(err))
		log.Fatal(err)
	}

	zap.L().Info("Response received from AI", zap.Int("candidate_count", len(resp.Candidates)))

	if len(resp.Candidates) == 0 {
		fmt.Printf("No candidates returned from AI\n")
		return models.Questions{}
	}

	fmt.Printf("Number of parts in first candidate: %d\n", len(resp.Candidates[0].Content.Parts))

	var bytes []byte
	var questions models.Questions
	for i, part := range resp.Candidates[0].Content.Parts {
		fmt.Printf("Processing part %d, type: %T\n", i, part)
		if txt, ok := part.(genai.Text); ok {
			bytes = []byte(txt)
			fmt.Printf("Found text part with %d bytes\n", len(bytes))
		}
	}

	if len(bytes) == 0 {
		fmt.Printf("No text content found in AI response\n")
		return models.Questions{}
	}

	fmt.Printf("AI Response (%d bytes): %s\n", len(bytes), string(bytes))

	err = json.Unmarshal(bytes, &questions)
	if err != nil {
		fmt.Printf("Error unmarshaling AI response: %v\n", err)
		fmt.Printf("Attempting to parse as raw JSON...\n")
		// Try to see if it's valid JSON at all
		var raw interface{}
		if jsonErr := json.Unmarshal(bytes, &raw); jsonErr != nil {
			fmt.Printf("Not valid JSON at all: %v\n", jsonErr)
		} else {
			fmt.Printf("Valid JSON but wrong structure. Raw content: %+v\n", raw)
		}
		return models.Questions{}
	}

	fmt.Printf("Successfully parsed %d questions from AI\n", len(questions.Questions))
	return questions

}
