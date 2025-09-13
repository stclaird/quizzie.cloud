package questions

import (
	"testing"

	"github.com/stclaird/quizzie.cloud/pkg/common/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateQuestion(t *testing.T) {
	// Setup in-memory SQLite DB for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	db.AutoMigrate(&models.Question{})

	h := handler{DB: db}

	question := models.Question{
		Category:    "Science",
		Subcategory: "Physics",
		Text:        "What is the speed of light?",
		Type:        "multiple-choice",
	}

	id, err := h.CreateQuestion(question)
	if err != nil {
		t.Errorf("CreateQuestion failed: %v", err)
	}
	if id == 0 {
		t.Errorf("Expected non-zero ID, got %d", id)
	}

	// Verify question was created
	var q models.Question
	if err := db.First(&q, id).Error; err != nil {
		t.Errorf("Question not found in DB: %v", err)
	}
	if q.Text != question.Text {
		t.Errorf("Expected question text %q, got %q", question.Text, q.Text)
	}
}
