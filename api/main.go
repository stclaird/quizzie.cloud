package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/stclaird/quizzie.cloud/pkg/common/db"
	"github.com/stclaird/quizzie.cloud/pkg/questions"
)

func main() {
    //get the app configuration
    config := GetConfig()
    router := gin.Default()

    //init the database object
    dbHandler,err := db.Init(config.dbUrl)
	if err != nil {
		log.Printf("main %s", err)
	}

    //init the questions
    questions.InitQuestions(config.questionPath, dbHandler)
    questions.RegisterRoutes(router, dbHandler)

    router.Run(config.port)
}