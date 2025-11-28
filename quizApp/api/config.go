package main

import (
	"fmt"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/spf13/viper"
)

type Config struct {
	port         string
	dbUrl        string
	questionPath string
	logLevel     string
	environment  string
}

func GetConfig() Config {
	// Set up viper to read from multiple sources
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	// Enable reading from environment variables
	viper.AutomaticEnv()

	// Set defaults
	viper.SetDefault("PORT", ":8080")
	viper.SetDefault("DB_URL", "quizzie.sqlite.db")
	viper.SetDefault("QUESTION_PATH", "./questionPack")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("APP_ENV", "development")

	// Get values (env vars will override .env file values, which override defaults)
	port := viper.GetString("PORT")
	dbUrl := viper.GetString("DB_URL")
	questionPath := viper.GetString("QUESTION_PATH")
	logLevel := viper.GetString("LOG_LEVEL")
	environment := viper.GetString("APP_ENV")

	// Ensure port has colon prefix
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	fmt.Println("questionPath:", questionPath)
	fmt.Println("port:", port)

	return Config{
		port:         port,
		dbUrl:        dbUrl,
		questionPath: questionPath,
		logLevel:     logLevel,
		environment:  environment,
	}
}

func CORSConfig(config Config) cors.Config {
	corsConfig := cors.DefaultConfig()

	// Use the port from config to build the allowed origin

	corsConfig.AllowOrigins = []string{
        "http://localhost:3000",  // React dev server
        "http://localhost:8081",  // Expo dev server
        "http://localhost:19000", // Expo web
        "http://localhost:19006", // Expo web alternative
        "http://127.0.0.1:8081",  // Alternative localhost
        "exp://192.168.1.0:19000", // Expo mobile
        "https://stclaird.github.io", // GitHub Pages
    }

	corsConfig.AllowCredentials = true
	corsConfig.AddAllowHeaders("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers", "Content-Type", "X-XSRF-TOKEN", "Accept", "Origin", "X-Requested-With", "Authorization")
	corsConfig.AddAllowMethods("GET", "POST", "PUT", "DELETE")
	return corsConfig
}
