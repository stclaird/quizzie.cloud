package questions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/stclaird/quizzie.cloud/pkg/common/models"
	"gorm.io/gorm"
)

func splitCatSubcat( catSubCat string) ( string, string ) {
    catSubCatSplit := strings.Split(catSubCat, "-")
    return catSubCatSplit[0], catSubCatSplit[1]
}

func SortString(w string) string {
    s := strings.Split(w, "")
    sort.Strings(s)
    return strings.Join(s, "")
}

func removeCorrectAnswerfield(questions []models.Question) (questionsNoAnswers []models.QuestionNoCorrectAnswer) {
    //Database retrieves objects including their correct answers - we need to copy this to an object with no correct answers
    //before sending it to the user. otherwise they can cheat!
    for _, question := range questions { //Copying by loop
        var questionnoanswer models.QuestionNoCorrectAnswer
        questionnoanswer.ID =  question.ID
        questionnoanswer.Category = question.Category
        questionnoanswer.Subcategory = question.Subcategory
        questionnoanswer.Text = question.Text
        questionnoanswer.Type = question.Type
        for _, answer := range question.Answers{
            a := struct{
                ID    uint
                Text   string `json:"text"`
            }{
                ID:    answer.ID,
                Text:  answer.Text,
            }
            questionnoanswer.Answers = append(questionnoanswer.Answers, a )
        }
        questionsNoAnswers = append(questionsNoAnswers, questionnoanswer)
     }

     return questionsNoAnswers
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