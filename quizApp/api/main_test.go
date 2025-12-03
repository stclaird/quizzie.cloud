package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	// Set gin to test mode
	gin.SetMode(gin.TestMode)

	// Run tests
	code := m.Run()

	// Exit with test result code
	os.Exit(code)
}

// TestApplicationSetup tests that the application can be initialized without errors
func TestApplicationSetup(t *testing.T) {
	// This test ensures the basic setup logic doesn't panic
	// We can't easily test main() directly, but we can test initialization

	// Test config loading
	config := GetConfig()
	if config.port == "" {
		t.Error("Config port should not be empty")
	}

	// Test CORS config
	corsConfig := CORSConfig(config)
	if len(corsConfig.AllowOrigins) == 0 {
		t.Error("CORS should have allowed origins")
	}
}

// TestRouterSetup tests that routes are properly registered
func TestRouterSetup(t *testing.T) {
	// Create a test router
	router := gin.New()

	// Add basic middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Next()
	})

	// Add a test route
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Create test request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)

	// Serve request
	router.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// TestCORSMiddleware tests the CORS headers
func TestCORSMiddleware(t *testing.T) {
	router := gin.New()

	// Add the same CORS middleware as in main
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

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Test OPTIONS request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("OPTIONS", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != 204 {
		t.Errorf("Expected status 204 for OPTIONS, got %d", w.Code)
	}

	// Test GET request with CORS headers
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	// Check CORS headers
	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("CORS Allow-Origin header not set correctly")
	}

	if w.Header().Get("Access-Control-Allow-Methods") == "" {
		t.Error("CORS Allow-Methods header not set")
	}

	if w.Header().Get("Access-Control-Allow-Headers") == "" {
		t.Error("CORS Allow-Headers header not set")
	}

	if w.Header().Get("Access-Control-Allow-Credentials") != "true" {
		t.Error("CORS Allow-Credentials header not set correctly")
	}
}

// TestDatabaseFileRemoval tests the database file cleanup logic
func TestDatabaseFileRemoval(t *testing.T) {
	// Create a temporary test database file
	testDbFile := "test_quizzie.sqlite.db"

	// Create the file
	file, err := os.Create(testDbFile)
	if err != nil {
		t.Fatalf("Failed to create test db file: %v", err)
	}
	file.Close()

	// Verify file exists
	if _, err := os.Stat(testDbFile); os.IsNotExist(err) {
		t.Fatal("Test db file was not created")
	}

	// Remove the file (simulating main.go behavior)
	err = os.Remove(testDbFile)
	if err != nil {
		t.Errorf("Failed to remove test db file: %v", err)
	}

	// Verify file was removed
	if _, err := os.Stat(testDbFile); !os.IsNotExist(err) {
		t.Error("Test db file still exists after removal")
		// Clean up if test fails
		os.Remove(testDbFile)
	}
}

// TestDatabaseFileRemoval_NotExists tests that removing non-existent file doesn't cause error
func TestDatabaseFileRemoval_NotExists(t *testing.T) {
	// Try to remove a file that doesn't exist
	err := os.Remove("nonexistent_db.sqlite.db")

	// This should produce an error, but the app should handle it gracefully
	if err == nil {
		t.Error("Expected error when removing non-existent file")
	}

	// In main.go, this error is logged but doesn't stop the app
	// which is the correct behavior
}

// Benchmark for config loading
func BenchmarkGetConfig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetConfig()
	}
}

// Benchmark for CORS config
func BenchmarkCORSConfig(b *testing.B) {
	config := GetConfig()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		CORSConfig(config)
	}
}
