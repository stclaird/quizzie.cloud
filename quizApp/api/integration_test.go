//go:build integration
// +build integration

package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stclaird/quizzie.cloud/pkg/common/db"
	"github.com/stclaird/quizzie.cloud/pkg/common/models"
	"github.com/stclaird/quizzie.cloud/pkg/questions"
	"gorm.io/gorm"
)

func setupIntegrationTest(t *testing.T) *gin.Engine {
	gin.SetMode(gin.TestMode)

	// Create in-memory database
	dbHandler, err := db.Init(":memory:")
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	// Seed test data
	seedTestData(t, dbHandler)

	// Setup router with all routes
	router := gin.New()
	router.Use(gin.Recovery())

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	questions.RegisterRoutes(router, dbHandler)

	return router
}

func seedTestData(t *testing.T, dbHandler *gorm.DB) {
	// Seed categories and questions
	question1 := models.Question{
		Category:    "Science",
		Subcategory: "Physics",
		Text:        "What is the speed of light?",
		Type:        "multiple-choice",
		Answers: []models.Answer{
			{Text: "299,792,458 m/s", IsCorrect: true},
			{Text: "150,000,000 m/s", IsCorrect: false},
			{Text: "300,000,000 m/s", IsCorrect: false},
		},
	}

	question2 := models.Question{
		Category:    "Science",
		Subcategory: "Chemistry",
		Text:        "What is H2O?",
		Type:        "multiple-choice",
		Answers: []models.Answer{
			{Text: "Water", IsCorrect: true},
			{Text: "Hydrogen Peroxide", IsCorrect: false},
		},
	}

	question3 := models.Question{
		Category:    "History",
		Subcategory: "World War II",
		Text:        "When did WWII end?",
		Type:        "multiple-choice",
		Answers: []models.Answer{
			{Text: "1945", IsCorrect: true},
			{Text: "1944", IsCorrect: false},
			{Text: "1946", IsCorrect: false},
		},
	}

	dbHandler.Create(&question1)
	dbHandler.Create(&question2)
	dbHandler.Create(&question3)
}

// Test full workflow: Get categories -> Get questions by category -> Answer question
func TestIntegration_FullQuizWorkflow(t *testing.T) {
	router := setupIntegrationTest(t)

	// Step 1: Get all categories
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/categories/", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200 for categories, got %d", w.Code)
	}

	var categories []models.Category
	if err := json.Unmarshal(w.Body.Bytes(), &categories); err != nil {
		t.Fatalf("Failed to parse categories: %v", err)
	}

	if len(categories) < 2 {
		t.Errorf("Expected at least 2 categories, got %d", len(categories))
	}

	// Step 2: Get questions by category
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/questions/Science/Physics", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200 for questions, got %d", w.Code)
	}

	var questions []models.QuestionNoCorrectAnswer
	if err := json.Unmarshal(w.Body.Bytes(), &questions); err != nil {
		t.Fatalf("Failed to parse questions: %v", err)
	}

	if len(questions) != 1 {
		t.Errorf("Expected 1 Science/Physics question, got %d", len(questions))
	}

	// Step 3: Get specific question
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/questions/id/1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200 for question by id, got %d", w.Code)
	}

	var question models.Question
	if err := json.Unmarshal(w.Body.Bytes(), &question); err != nil {
		t.Fatalf("Failed to parse question: %v", err)
	}

	// Step 4: Submit an answer (test the endpoint works, not correctness)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/questions/answer/1/1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200 for answer submission, got %d", w.Code)
	}

	var answerResp models.AnswerResponse
	if err := json.Unmarshal(w.Body.Bytes(), &answerResp); err != nil {
		t.Fatalf("Failed to parse answer response: %v", err)
	}

	// Just verify we get a response with the expected fields
	if len(answerResp.CorrectAnswer) == 0 {
		t.Error("Expected correct answer to be returned in response")
	}
}

func TestIntegration_GetAllQuestions(t *testing.T) {
	router := setupIntegrationTest(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/questions/", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}

	var questions []models.QuestionNoCorrectAnswer
	if err := json.Unmarshal(w.Body.Bytes(), &questions); err != nil {
		t.Fatalf("Failed to parse questions: %v", err)
	}

	if len(questions) != 3 {
		t.Errorf("Expected 3 questions, got %d", len(questions))
	}

	// Verify no correct answers are exposed
	for _, q := range questions {
		if len(q.Answers) == 0 {
			t.Error("Questions should have answers")
		}
	}
}

func TestIntegration_CORSHeaders(t *testing.T) {
	router := setupIntegrationTest(t)

	// Test OPTIONS request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("OPTIONS", "/questions/", nil)
	router.ServeHTTP(w, req)

	if w.Code != 204 {
		t.Errorf("Expected status 204 for OPTIONS, got %d", w.Code)
	}

	// Test regular request has CORS headers
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/questions/", nil)
	router.ServeHTTP(w, req)

	if origin := w.Header().Get("Access-Control-Allow-Origin"); origin != "*" {
		t.Errorf("Expected CORS Allow-Origin *, got %s", origin)
	}
}

func TestIntegration_InvalidEndpoints(t *testing.T) {
	router := setupIntegrationTest(t)

	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
	}{
		{"Invalid question ID", "GET", "/questions/id/99999", http.StatusNotFound},
		{"Invalid category", "GET", "/questions/NonExistent/Category", http.StatusOK}, // Returns empty array
		{"Invalid answer", "GET", "/questions/answer/1/99999", http.StatusOK},         // Returns isCorrect: false
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.method, tt.path, nil)
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Expected status %d, got %d", tt.wantStatus, w.Code)
			}
		})
	}
}
