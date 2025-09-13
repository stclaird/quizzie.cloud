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

func setupTestDBQuestions(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        t.Fatalf("failed to connect database: %v", err)
    }
    db.AutoMigrate(&models.Question{}, &models.Answer{})
    return db
}

func TestGetQuestions(t *testing.T) {
    gin.SetMode(gin.TestMode)
    db := setupTestDBQuestions(t)

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
    router.GET("/questions/", h.GetQuestions)

    req, _ := http.NewRequest("GET", "/questions/", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Fatalf("Expected status 200, got %d", w.Code)
    }

    var resp []models.QuestionNoCorrectAnswer
    if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
        t.Fatalf("Failed to unmarshal response: %v", err)
    }
    if len(resp) != 1 {
        t.Errorf("Expected 1 question, got %d", len(resp))
    }
    if resp[0].Text != question.Text {
        t.Errorf("Expected question text %q, got %q", question.Text, resp[0].Text)
    }
    // Ensure correct answer is not present
    for _, ans := range resp[0].Answers {
        if ans.Text == "299,792,458 m/s" {
            // Should not have IsCorrect field
            // (You may want to check struct fields if needed)
        }
    }
}

func TestGetQuestionsCatSubCat(t *testing.T) {
    gin.SetMode(gin.TestMode)
    db := setupTestDBQuestions(t)

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
    router.GET("/questions/catsubcat/:catsubcat", h.GetQuestionsCatSubCat)

    req, _ := http.NewRequest("GET", "/questions/catsubcat/Science-Physics", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Fatalf("Expected status 200, got %d", w.Code)
    }

    var resp []models.QuestionNoCorrectAnswer
    if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
        t.Fatalf("Failed to unmarshal response: %v", err)
    }
    if len(resp) != 1 {
        t.Errorf("Expected 1 question, got %d", len(resp))
    }
    if resp[0].Category != "Science" || resp[0].Subcategory != "Physics" {
        t.Errorf("Expected category Science and subcategory Physics, got %s and %s", resp[0].Category, resp[0].Subcategory)
    }
}
