//Package questions
// Provide api functionailty for question objects

package questions

import (
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SortString(w string) string {
    s := strings.Split(w, "")
    sort.Strings(s)
    return strings.Join(s, "")
}

type handler struct {
    DB *gorm.DB
}

func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
    h := &handler{
        DB: db,
    }

    questionRoutes := router.Group("/questions")
    questionRoutes.GET("/", h.GetQuestions)
	questionRoutes.GET("/bycatsubcat/:catsubcat", h.GetQuestions)
    questionRoutes.GET("/byid/:id", h.GetQuestion)
	questionRoutes.GET("/answer/:id/:answer", h.Answers)

    categoryRoutes := router.Group("/categories")
    categoryRoutes.GET("/", h.GetCategories)
}
