package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/stclaird/quizzie.cloud/questionCreator/pkg/logger"
	"github.com/stclaird/quizzie.cloud/questionCreator/pkg/question"
	"github.com/stclaird/quizzie.cloud/questionCreator/pkg/util"
)

func Configs(config util.Config) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Set("config", config)
    }
}

func main() {
    //get the app configuration first
    config,err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

    // Initialize structured logging with config log level and environment
    if err := logger.InitLoggerWithConfig(config.LogLevel, config.Environment); err != nil {
        log.Fatal("Failed to initialize logger:", err)
    }
    defer logger.Sync()

    // Create router without default middleware
    router := gin.New()

    // Add our custom Zap middleware
    router.Use(logger.GinZapMiddleware())
    router.Use(logger.GinRecoveryWithZap(true))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
		  "status": "OK",
		})
	  })

    router.Use(cors.New(util.CORSConfig()))
	router.Use(Configs(config))
    question.RegisterRoutes(router, config)

    zap.L().Info("Starting question creator service", zap.String("port", config.Port))
    router.Run(config.Port)
}
