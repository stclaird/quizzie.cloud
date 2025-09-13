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

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	db.AutoMigrate(&models.Question{})
	return db
}

func TestGetCategories(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB(t)

	// Insert sample questions
	questions := []models.Question{
		{Category: "Science", Subcategory: "Physics", Text: "Q1"},
		{Category: "Science", Subcategory: "Chemistry", Text: "Q2"},
		{Category: "Literature", Subcategory: "Poetry", Text: "Q3"},
	}
	for _, q := range questions {
		db.Create(&q)
	}

	h := handler{DB: db}
	router := gin.New()
	router.GET("/categories", h.GetCategories)

	req, _ := http.NewRequest("GET", "/categories", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}

	var resp []*models.Category
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(resp) != 2 {
		t.Errorf("Expected 2 categories, got %d", len(resp))
	}

	// Check that subcategories are present
	for _, cat := range resp {
		if cat.CategoryName == "Science" && len(cat.SubCategories) != 2 {
			t.Errorf("Expected 2 subcategories for Science, got %d", len(cat.SubCategories))
		}
		if cat.CategoryName == "Literature" && len(cat.SubCategories) != 1 {
			t.Errorf("Expected 1 subcategory for Literature, got %d", len(cat.SubCategories))
		}
	}
}
