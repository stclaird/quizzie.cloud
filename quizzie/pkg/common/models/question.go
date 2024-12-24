package models

import (
	"time"

	"gorm.io/gorm"
)

type Question struct {
	ID        uint           `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Text        string `json:"questionText"`
	Type        string `json:"type"`
	Category    string `json:"category"`
	Subcategory string `json:"subcategory"`
	DateAdded   string `json:"dateAdded"`
	Answers     []Answer `json:"answers"`
}

type Answer struct {
	ID        uint           `gorm:"primaryKey"`
	QuestionID uint	`gorm:"questionid"`
	Text string `json:"text"`
	IsCorrect bool   `json:"iscorrect"`
}

type QuestionNoCorrectAnswer struct {
	ID          uint
	Text        string `json:"questionText"`
	Type        string `json:"type"`
	Category    string `json:"category"`
	Subcategory string `json:"subcategory"`
	DateAdded   string `json:"dateAdded"`
	Answers     []struct {
		ID        uint
		Text      string `json:"text"`
	} `json:"answers"`
}

type AnswerResponse struct {
	IsCorrect bool `json:"iscorrect"`
	CorrectAnswer []Answer `json:"correctanswer"`
}
