package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/stclaird/quizzie.cloud/pkg/common/db"
	"github.com/stclaird/quizzie.cloud/pkg/questions"
)

func main() {
    //get the app configuration
    config := GetConfig()

    //remove any existing db files
    //as we build them from scratch
    e := os.Remove(config.dbUrl)
    if e != nil {
        log.Println("No existing db found")
    }

    router := gin.Default()
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
        log.Printf("main %s", err)
    }

    //init the questions
    questions.InitQuestions(config.questionPath, dbHandler)
    questions.RegisterRoutes(router, dbHandler)

    router.Run(config.port)
}
