package question

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/google/uuid"

	"github.com/stclaird/quizzie.cloud/questionCreator/pkg/models"
)

func SortString(w string) string {
    s := strings.Split(w, "")
    sort.Strings(s)
    return strings.Join(s, "")
}

func generateQuestionFileName(q models.QuestionOut) (name string ){
    //Generate File name for Question
    slug := generateQuestionID(q)

    cat := strings.ToLower(q.Category)
    catsafe := safeFileName(cat)

    subcat := strings.ToLower(q.Subcategory)
    subcatsafe := safeFileName(subcat)

    name = fmt.Sprintf("%s-%s-%s.json", catsafe, subcatsafe, slug )

    return name
}

func generateQuestionID(q models.QuestionOut) string {
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
