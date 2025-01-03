package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/stclaird/quizzie.cloud/questionCreator/pkg/question"
	"github.com/stclaird/quizzie.cloud/questionCreator/pkg/util"
)

func Configs(config util.Config) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Set("config", config)
    }
}

func main() {
    //get the app configuration
    config,err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config")
	}

    router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
		  "status": "OK",
		})
	  })

    router.Use(cors.New(util.CORSConfig()))
	router.Use(Configs(config))
    question.RegisterRoutes(router, config)
    router.Run(config.Port)
}
