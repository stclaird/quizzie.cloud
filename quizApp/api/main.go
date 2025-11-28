package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/stclaird/quizzie.cloud/pkg/common/db"
	"github.com/stclaird/quizzie.cloud/pkg/common/logger"
	"github.com/stclaird/quizzie.cloud/pkg/questions"
	"go.uber.org/zap"
)

func main() {
    //get the app configuration first
    config := GetConfig()

    // Initialize logger with config log level and environment
    if err := logger.InitLoggerWithConfig(config.logLevel, config.environment); err != nil {
        log.Fatal("Failed to initialize logger:", err)
    }
    defer logger.Sync()

    zapLogger := logger.GetLogger()
    zapLogger.Info("Starting Quiz App API")
    zapLogger.Info("Configuration loaded",
        zap.String("port", config.port),
        zap.String("questionPath", config.questionPath))

    //remove any existing db files
    //as we build them from scratch
    e := os.Remove(config.dbUrl)
    if e != nil {
        zapLogger.Info("No existing database found", zap.String("dbUrl", config.dbUrl))
    } else {
        zapLogger.Info("Removed existing database", zap.String("dbUrl", config.dbUrl))
    }

    // Create router without default middleware
    router := gin.New()

    // Add our custom Zap middleware
    router.Use(logger.GinZapMiddleware())
    router.Use(logger.GinRecoveryWithZap(true))

   // router.Use(static.Serve("/", static.LocalFile("../ui/build", true)))

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

    //init the database object
    dbHandler, err := db.Init(config.dbUrl)
    if err != nil {
        zapLogger.Error("Failed to initialize database", zap.Error(err))
        return
    }
    zapLogger.Info("Database initialized successfully")

    //init the questions
    questionsCount := questions.InitQuestions(config.questionPath, dbHandler)
    zapLogger.Info("Questions initialized", zap.Int("count", len(questionsCount)))

    questions.RegisterRoutes(router, dbHandler)
    zapLogger.Info("Routes registered")

    zapLogger.Info("Starting HTTP server", zap.String("port", config.port))
    if err := router.Run(config.port); err != nil {
        zapLogger.Error("Failed to start server", zap.Error(err))
    }
}
