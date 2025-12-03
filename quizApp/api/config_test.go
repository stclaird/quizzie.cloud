package main

import (
	"os"
	"testing"

	"github.com/spf13/viper"
)

func TestGetConfig_Defaults(t *testing.T) {
	viper.Reset()
	config := GetConfig()

	if config.port != ":8080" {
		t.Errorf("Expected default port :8080, got %s", config.port)
	}

	if config.dbUrl != "quizzie.sqlite.db" {
		t.Errorf("Expected default dbUrl quizzie.sqlite.db, got %s", config.dbUrl)
	}

	if config.questionPath != "./questionPack" {
		t.Errorf("Expected default questionPath ./questionPack, got %s", config.questionPath)
	}
}

func TestGetConfig_EnvironmentOverrides(t *testing.T) {
	viper.Reset()

	os.Setenv("PORT", "9090")
	os.Setenv("DB_URL", "test.db")
	os.Setenv("QUESTION_PATH", "/custom/path")

	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("DB_URL")
		os.Unsetenv("QUESTION_PATH")
	}()

	config := GetConfig()

	if config.port != ":9090" {
		t.Errorf("Expected port :9090, got %s", config.port)
	}

	if config.dbUrl != "test.db" {
		t.Errorf("Expected dbUrl test.db, got %s", config.dbUrl)
	}

	if config.questionPath != "/custom/path" {
		t.Errorf("Expected questionPath /custom/path, got %s", config.questionPath)
	}
}

func TestGetConfig_PortPrefix(t *testing.T) {
	viper.Reset()

	os.Setenv("PORT", "8080")
	defer os.Unsetenv("PORT")

	config := GetConfig()

	if config.port != ":8080" {
		t.Errorf("Expected port to have colon prefix :8080, got %s", config.port)
	}
}

func TestCORSConfig(t *testing.T) {
	viper.Reset()
	config := GetConfig()

	corsConfig := CORSConfig(config)

	expectedOrigins := []string{
		"http://localhost:3000",
		"http://localhost:8081",
		"http://localhost:19000",
		"https://stclaird.github.io",
	}

	for _, expected := range expectedOrigins {
		found := false
		for _, origin := range corsConfig.AllowOrigins {
			if origin == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected origin %s not found in CORS config", expected)
		}
	}

	if !corsConfig.AllowCredentials {
		t.Error("Expected CORS to allow credentials")
	}
}
