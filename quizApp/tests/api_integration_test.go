//go:build integration
// +build integration

package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stclaird/quizzie.cloud/pkg/common/models"
)

var apiURL string

func init() {
	apiURL = os.Getenv("API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080"
	}
}

func waitForAPI(t *testing.T) {
	client := &http.Client{Timeout: 2 * time.Second}
	maxAttempts := 30

	for i := 0; i < maxAttempts; i++ {
		resp, err := client.Get(apiURL + "/categories/")
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(1 * time.Second)
	}

	t.Fatal("API did not become ready in time")
}

func TestMain(m *testing.M) {
	// Wait for API to be ready
	fmt.Println("Waiting for API to be ready...")
	code := m.Run()
	os.Exit(code)
}

func TestAPI_GetCategories(t *testing.T) {
	waitForAPI(t)

	resp, err := http.Get(apiURL + "/categories/")
	if err != nil {
		t.Fatalf("Failed to get categories: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var categories []models.Category
	if err := json.Unmarshal(body, &categories); err != nil {
		t.Fatalf("Failed to parse categories: %v", err)
	}

	if len(categories) == 0 {
		t.Error("Expected at least one category")
	}

	t.Logf("Found %d categories", len(categories))
}

func TestAPI_GetAllQuestions(t *testing.T) {
	waitForAPI(t)

	resp, err := http.Get(apiURL + "/questions/")
	if err != nil {
		t.Fatalf("Failed to get questions: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var questions []models.QuestionNoCorrectAnswer
	if err := json.Unmarshal(body, &questions); err != nil {
		t.Fatalf("Failed to parse questions: %v", err)
	}

	if len(questions) == 0 {
		t.Error("Expected at least one question")
	}

	// Verify structure
	for _, q := range questions {
		if q.Text == "" {
			t.Error("Question text should not be empty")
		}
		if len(q.Answers) == 0 {
			t.Error("Question should have answers")
		}
	}

	t.Logf("Found %d questions", len(questions))
}

func TestAPI_GetQuestionsByCategory(t *testing.T) {
	waitForAPI(t)

	// First get a category
	resp, err := http.Get(apiURL + "/categories/")
	if err != nil {
		t.Fatalf("Failed to get categories: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var categories []models.Category
	json.Unmarshal(body, &categories)

	if len(categories) == 0 {
		t.Skip("No categories available")
	}

	category := categories[0]
	if len(category.SubCategories) == 0 {
		t.Skip("Category has no subcategories")
	}

	subcategory := category.SubCategories[0]
	url := fmt.Sprintf("%s/questions/%s/%s", apiURL, category.CategoryName, subcategory.SubCategoryName)

	resp, err = http.Get(url)
	if err != nil {
		t.Fatalf("Failed to get questions by category: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	body, _ = io.ReadAll(resp.Body)
	var questions []models.QuestionNoCorrectAnswer
	if err := json.Unmarshal(body, &questions); err != nil {
		t.Fatalf("Failed to parse questions: %v", err)
	}

	t.Logf("Found %d questions in category %s/%s", len(questions), category.CategoryName, subcategory.SubCategoryName)
}

func TestAPI_AnswerQuestion(t *testing.T) {
	waitForAPI(t)

	// Get all questions first to find one with an ID and answers
	resp, err := http.Get(apiURL + "/questions/")
	if err != nil {
		t.Fatalf("Failed to get questions: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var questions []models.QuestionNoCorrectAnswer
	if err := json.Unmarshal(body, &questions); err != nil {
		t.Fatalf("Failed to parse questions: %v", err)
	}

	if len(questions) == 0 {
		t.Skip("No questions available")
	}

	// Use first question that has answers
	question := questions[0]
	if len(question.Answers) == 0 {
		t.Skip("Question has no answers")
	}

	// Try answering with first answer
	answerID := question.Answers[0].ID
	url := fmt.Sprintf("%s/questions/answer/%d/%d", apiURL, question.ID, answerID)

	resp, err = http.Get(url)
	if err != nil {
		t.Fatalf("Failed to submit answer: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	body, _ = io.ReadAll(resp.Body)
	var answerResp models.AnswerResponse
	if err := json.Unmarshal(body, &answerResp); err != nil {
		t.Fatalf("Failed to parse answer response: %v", err)
	}

	t.Logf("Answer was correct: %v", answerResp.IsCorrect)
}

func TestAPI_CORSHeaders(t *testing.T) {
	waitForAPI(t)

	// Test OPTIONS request
	req, _ := http.NewRequest("OPTIONS", apiURL+"/questions/", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to make OPTIONS request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		t.Errorf("Expected status 204 for OPTIONS, got %d", resp.StatusCode)
	}

	// Check CORS headers on regular request
	resp, err = http.Get(apiURL + "/questions/")
	if err != nil {
		t.Fatalf("Failed to get questions: %v", err)
	}
	defer resp.Body.Close()

	origin := resp.Header.Get("Access-Control-Allow-Origin")
	if origin == "" {
		t.Error("CORS Allow-Origin header not set")
	}

	t.Logf("CORS Allow-Origin: %s", origin)
}

func TestAPI_ResponseTimes(t *testing.T) {
	waitForAPI(t)

	endpoints := []string{
		"/categories/",
		"/questions/",
	}

	for _, endpoint := range endpoints {
		start := time.Now()
		resp, err := http.Get(apiURL + endpoint)
		duration := time.Since(start)

		if err != nil {
			t.Errorf("Failed to get %s: %v", endpoint, err)
			continue
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200 for %s, got %d", endpoint, resp.StatusCode)
		}

		t.Logf("%s responded in %v", endpoint, duration)

		if duration > 2*time.Second {
			t.Errorf("%s took too long: %v", endpoint, duration)
		}
	}
}
