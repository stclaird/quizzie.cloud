//Package questions
// Provide api functionailty for question objects

package questions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/stclaird/quizzie.cloud/pkg/common/models"
	"gorm.io/gorm"
)

type handler struct {
    DB *gorm.DB
}

func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
    h := &handler{
        DB: db,
    }

    questionRoutes := router.Group("/questions")
    questionRoutes.GET("/", h.GetQuestions)
    questionRoutes.GET("/:id", h.GetQuestion)
	questionRoutes.GET("/answer/:id/:answer", h.Answers)

    categoryRoutes := router.Group("/categories")
    categoryRoutes.GET("/", h.GetCategories)
}

func InitQuestions(questionPack string, db *gorm.DB) (allQuestions []models.Question) {
	//import questions from a json file ready for adding to the DB
	//returns a slice of question stuct types

    h := &handler{
        DB: db,
    }

	files, err := ioutil.ReadDir(questionPack)
	if err != nil {
		log.Fatal(err)
	}

	for _, File := range files {
		fileExtension := filepath.Ext(File.Name())
		if fileExtension == ".json" {
			var questionsObj []models.Question
			fmt.Printf("Loading %s\n", File.Name())
			filePath := fmt.Sprintf("%s/%s", questionPack, File.Name())
			jsonFile, err := os.Open(filePath)
			if err != nil {
				log.Println("Error", err)
			}
			defer jsonFile.Close()
			byteValue, _ := ioutil.ReadAll(jsonFile)
			json.Unmarshal(byteValue, &questionsObj)
			for _, question := range questionsObj {
				allQuestions = append(allQuestions, question)
                h.CreateQuestion(question)
			}
		}
	}


	return allQuestions
}