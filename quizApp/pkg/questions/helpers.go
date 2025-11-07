package questions

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/stclaird/quizzie.cloud/pkg/common/models"
	"gorm.io/gorm"
)

func SortString(w string) string {
	s := strings.Split(w, "")
	sort.Strings(s)
	return strings.Join(s, "")
}

func removeCorrectAnswerfield(questions []models.Question) (questionsNoAnswers []models.QuestionNoCorrectAnswer) {
	//Database retrieves objects including their correct answers - we need to copy this to an object with no correct answers
	//before sending it to the user. otherwise they can cheat!
	for _, question := range questions { //Copying by loop
		var questionnoanswer models.QuestionNoCorrectAnswer
		questionnoanswer.ID = question.ID
		questionnoanswer.Category = question.Category
		questionnoanswer.Subcategory = question.Subcategory
		questionnoanswer.Text = question.Text
		questionnoanswer.Type = question.Type
		for _, answer := range question.Answers {
			a := struct {
				ID   uint
				Text string `json:"text"`
			}{
				ID:   answer.ID,
				Text: answer.Text,
			}
			questionnoanswer.Answers = append(questionnoanswer.Answers, a)
		}
		questionsNoAnswers = append(questionsNoAnswers, questionnoanswer)
	}

	return questionsNoAnswers
}

func InitQuestions(questionPack string, db *gorm.DB) (allQuestions []models.Question) {
	h := &handler{
		DB: db,
	}

	// Check if questionPack is a URL
	if strings.HasPrefix(questionPack, "http://") || strings.HasPrefix(questionPack, "https://") {
		// Handle HTTP URL
		allQuestions = loadQuestionsFromURL(questionPack, h)
	} else {
		// Handle local directory
		allQuestions = loadQuestionsFromLocal(questionPack, h)
	}

	return allQuestions
}

func loadQuestionsFromLocal(questionPack string, h *handler) []models.Question {
	var allQuestions []models.Question

	files, err := os.ReadDir(questionPack)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			filePath := filepath.Join(questionPack, file.Name())
			questions := processJSONFile(filePath, h, true)
			allQuestions = append(allQuestions, questions...)
		}
	}

	return allQuestions
}

func loadQuestionsFromURL(url string, h *handler) []models.Question {
	var allQuestions []models.Question

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Error fetching from URL:", err)
	}
	defer resp.Body.Close()

	questions := processJSONReader(resp.Body, h, url)
	allQuestions = append(allQuestions, questions...)

	return allQuestions
}

func processJSONFile(filePath string, h *handler, isLocal bool) []models.Question {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		log.Println("Error opening file:", err)
		return nil
	}
	defer jsonFile.Close()

	return processJSONReader(jsonFile, h, filePath)
}

func processJSONReader(reader io.Reader, h *handler, source string) []models.Question {
	var questionsObj []models.Question
	var allQuestions []models.Question

	fmt.Printf("Loading from %s\n", source)

	byteValue, err := io.ReadAll(reader)
	if err != nil {
		log.Println("Error reading:", err)
		return nil
	}

	if err := json.Unmarshal(byteValue, &questionsObj); err != nil {
		log.Println("Error unmarshaling JSON:", err)
		return nil
	}

	for i, question := range questionsObj {
		questionsObj[i].Category = strings.Replace(question.Category, "-", " ", -1)
		fmt.Printf("Adding Question Category: %s\n", questionsObj[i].Category)

		allQuestions = append(allQuestions, questionsObj[i])
		h.CreateQuestion(questionsObj[i])
	}

	fmt.Printf("Loaded %d questions from %s\n", len(questionsObj), source)
	return allQuestions
}
