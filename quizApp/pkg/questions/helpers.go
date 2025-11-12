package questions

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
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

		// Create answers slice
		var answers []struct {
			ID   uint
			Text string `json:"text"`
		}

		for _, answer := range question.Answers {
			a := struct {
				ID   uint
				Text string `json:"text"`
			}{
				ID:   answer.ID,
				Text: answer.Text,
			}
			answers = append(answers, a)
		}

		// Shuffle the answers to prevent predictable ordering
		rand.Shuffle(len(answers), func(i, j int) {
			answers[i], answers[j] = answers[j], answers[i]
		})

		questionnoanswer.Answers = answers
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
		// Handle HTTP URL with fallback
		allQuestions = loadQuestionsFromURLWithFallback(questionPack, h)
	} else {
		// Handle local directory with fallback
		allQuestions = loadQuestionsFromLocalWithFallback(questionPack, h)
	}

	// Final fallback if no questions loaded at all
	if len(allQuestions) == 0 {
		log.Println("No questions loaded from any source, using emergency fallback")
		allQuestions = loadEmergencyFallback(h)
	}

	return allQuestions
}

func loadQuestionsFromLocalWithFallback(questionPack string, h *handler) []models.Question {
	// Try the specified local path first
	allQuestions := loadQuestionsFromLocal(questionPack, h)

	// If that failed and it's not already the default path, try the default
	if len(allQuestions) == 0 && questionPack != "./questionPack" {
		log.Printf("Failed to load from %s, trying ./questionPack", questionPack)
		allQuestions = loadQuestionsFromLocal("./questionPack", h)
	}

	// If still no questions, try the specific defaultQuestion.json file
	if len(allQuestions) == 0 {
		log.Println("Failed to load from directory, trying defaultQuestion.json")
		questions := processJSONFile("./questionPack/defaultQuestion.json", h, false)
		if questions != nil {
			allQuestions = append(allQuestions, questions...)
		}
	}

	return allQuestions
}

func loadQuestionsFromLocal(questionPack string, h *handler) []models.Question {
	var allQuestions []models.Question

	files, err := os.ReadDir(questionPack)
	if err != nil {
		log.Printf("Error reading directory %s: %v", questionPack, err)
		return allQuestions // Return empty slice instead of crashing
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			filePath := filepath.Join(questionPack, file.Name())
			questions := processJSONFile(filePath, h, true)
			if questions != nil {
				allQuestions = append(allQuestions, questions...)
			}
		}
	}

	return allQuestions
}

func loadQuestionsFromURLWithFallback(url string, h *handler) []models.Question {
	// Try to load from URL first
	allQuestions := loadQuestionsFromURL(url, h)

	// If URL loading failed, fallback to local defaultQuestion.json
	if len(allQuestions) == 0 {
		log.Println("Failed to load questions from URL, falling back to local defaultQuestion.json")
		fallbackPath := "./questionPack/defaultQuestion.json"
		questions := processJSONFile(fallbackPath, h, false)
		if questions != nil {
			allQuestions = append(allQuestions, questions...)
		}
	}

	return allQuestions
}

func loadQuestionsFromURL(url string, h *handler) []models.Question {
	var allQuestions []models.Question

	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error fetching from URL:", err)
		return allQuestions // Return empty slice instead of crashing
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("HTTP error: %d %s", resp.StatusCode, resp.Status)
		return allQuestions
	}

	// Check if this is a ZIP file by URL or content type
	contentType := resp.Header.Get("Content-Type")
	isZipFile := strings.Contains(url, ".zip") ||
		strings.Contains(contentType, "application/zip") ||
		strings.Contains(contentType, "application/x-zip") ||
		strings.Contains(contentType, "application/octet-stream") // GitHub releases often use this

	fmt.Printf("URL: %s, Content-Type: %s, isZipFile: %t\n", url, contentType, isZipFile)

	if isZipFile {
		allQuestions = processZipReader(resp.Body, h, url)
	} else {
		questions := processJSONReader(resp.Body, h, url)
		if questions != nil {
			allQuestions = append(allQuestions, questions...)
		}
	}

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

func processZipReader(reader io.Reader, h *handler, source string) []models.Question {
	var allQuestions []models.Question

	// Read the entire response body into memory
	data, err := io.ReadAll(reader)
	if err != nil {
		log.Printf("Error reading ZIP data from %s: %v", source, err)
		return allQuestions
	}

	fmt.Printf("Downloaded %d bytes from %s\n", len(data), source)

	// Create a zip reader from the data
	zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		log.Printf("Error creating ZIP reader for %s: %v", source, err)
		return allQuestions
	}

	fmt.Printf("Processing ZIP file from %s with %d files\n", source, len(zipReader.File))

	// Process each file in the ZIP
	for _, file := range zipReader.File {
		// Only process JSON files
		if filepath.Ext(file.Name) == ".json" {
			fmt.Printf("Processing JSON file: %s\n", file.Name)

			fileReader, err := file.Open()
			if err != nil {
				log.Printf("Error opening file %s in ZIP: %v", file.Name, err)
				continue
			}
			defer fileReader.Close()

			questions := processJSONReader(fileReader, h, fmt.Sprintf("%s/%s", source, file.Name))
			if questions != nil {
				allQuestions = append(allQuestions, questions...)
			}
		}
	}

	fmt.Printf("Loaded %d questions from ZIP file %s\n", len(allQuestions), source)
	return allQuestions
}

func loadEmergencyFallback(h *handler) []models.Question {
	// Create a minimal hardcoded question as absolute last resort
	question := models.Question{
		Category:    "General",
		Subcategory: "Default",
		Text:        "What is 2 + 2?",
		Type:        "multiple-choice",
		Answers: []models.Answer{
			{Text: "3", IsCorrect: false},
			{Text: "4", IsCorrect: true},
			{Text: "5", IsCorrect: false},
		},
	}

	h.CreateQuestion(question)
	return []models.Question{question}
}
