package questions

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stclaird/quizzie.cloud/pkg/common/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDBWithAnswers(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        t.Fatalf("failed to connect database: %v", err)
    }
    db.AutoMigrate(&models.Question{}, &models.Answer{})
    return db
}

func TestGetQuestion(t *testing.T) {
    gin.SetMode(gin.TestMode)
    db := setupTestDBWithAnswers(t)

    question := models.Question{
        Category:    "Science",
        Subcategory: "Physics",
        Text:        "What is the speed of light?",
        Type:        "multiple-choice",
    }
    db.Create(&question)

    h := handler{DB: db}
    router := gin.New()
    router.GET("/questions/id/:id", h.GetQuestion)

    req, _ := http.NewRequest("GET", "/questions/id/1", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Fatalf("Expected status 200, got %d", w.Code)
    }

    var resp models.Question
    if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
        t.Fatalf("Failed to unmarshal response: %v", err)
    }
    if resp.Text != question.Text {
        t.Errorf("Expected question text %q, got %q", question.Text, resp.Text)
    }
}

func TestAnswers(t *testing.T) {
    gin.SetMode(gin.TestMode)
    db := setupTestDBWithAnswers(t)

    question := models.Question{
        Category:    "Science",
        Subcategory: "Physics",
        Text:        "What is the speed of light?",
        Type:        "multiple-choice",
        Answers: []models.Answer{
            {Text: "299,792,458 m/s", IsCorrect: true},
            {Text: "150,000,000 m/s", IsCorrect: false},
        },
    }
    db.Create(&question)
    for i := range question.Answers {
        question.Answers[i].QuestionID = question.ID
        db.Create(&question.Answers[i])
    }

    h := handler{DB: db}
    router := gin.New()
    router.GET("/questions/answer/:id/:answer", h.Answers)

    // Find the correct answer ID
    var correctAnswer models.Answer
    db.Where("is_correct = ?", true).First(&correctAnswer)

    // Test correct answer
    req, _ := http.NewRequest("GET", "/questions/answer/1/"+string(rune(correctAnswer.ID)), nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Fatalf("Expected status 200, got %d", w.Code)
    }

    var resp models.AnswerResponse
    if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
        t.Fatalf("Failed to unmarshal response: %v", err)
    }
    if !resp.IsCorrect {
        t.Errorf("Expected IsCorrect true, got false")
    }

    // Test incorrect answer
    req, _ = http.NewRequest("GET", "/questions/answer/1/999", nil)
    w = httptest.NewRecorder()
    router.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Fatalf("Expected status 200, got %d", w.Code)
    }
    if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
        t.Fatalf("Failed to unmarshal response: %v", err)
    }
    if resp.IsCorrect {
        t.Errorf("Expected IsCorrect false, got true")
    }
}
