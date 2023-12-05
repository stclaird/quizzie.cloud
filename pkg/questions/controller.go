package questions

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type handler struct {
    DB *gorm.DB
}

func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
    h := &handler{
        DB: db,
    }

    routes := router.Group("/questions")
    routes.GET("/", h.GetQuestion)
    routes.GET("/:id", h.GetQuestion)
}