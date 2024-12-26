package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/stclaird/quizzie.cloud/questionCreator/pkg/question"
	"github.com/stclaird/quizzie.cloud/questionCreator/pkg/util"
)

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
    question.RegisterRoutes(router)
    router.Run(config.Port)
}
