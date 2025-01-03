package question

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

func generateQuestionFileName(cat string, subcat string) (name string ){
    //Generate File name for Question
    slug := generateQuestionID()

    catlower := strings.ToLower(cat)
    catsafe := safeFileName(catlower)

    subcatlower := strings.ToLower(subcat)
    subcatsafe := safeFileName(subcatlower)

    name = fmt.Sprintf("%s-%s-%s.json", catsafe, subcatsafe, slug )

    return name
}

func generateQuestionID() string {
    se := uuid.New().String()
    fmt.Println(se)
    return se
}

func safeFileName(s string) string {
    //Remove any undesirable characters from a file name
    re := regexp.MustCompile(`[\\/:*?"<>|=]`)
    cleanString := re.ReplaceAllString(s, "")

    return cleanString
}

func sortString(w string) string {
    s := strings.Split(w, "")
    sort.Strings(s)
    return strings.Join(s, "")
}

func splitString(w string) []string {
    s := strings.Split(w, " ")
    return s
}

func createDate() string {
    //Return date in format 30-April-2018
    currentTime := time.Now()
    dateOut := fmt.Sprintf("%v-%v-%v", currentTime.Day(), currentTime.Month(), currentTime.Year())

    return dateOut
}
