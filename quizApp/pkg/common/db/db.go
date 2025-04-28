package db

import (
	"github.com/stclaird/quizzie.cloud/pkg/common/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Init(dbPath string) (*gorm.DB, error) {
	//Function for Initing the DataBase

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		QueryFields: true,
	})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.Question{}, &models.Answer{})

	return db, err
}
