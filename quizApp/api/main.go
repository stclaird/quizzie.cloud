package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
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
    router.Use(static.Serve("/", static.LocalFile("../ui/build", true)))

    router.Use(cors.New(CORSConfig(config)))

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
