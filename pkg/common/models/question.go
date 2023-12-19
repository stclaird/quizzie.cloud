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
	Text        string `json:"text"`
	Type        string `json:"type"`
	Category    string `json:"category"`
	Subcategory string `json:"subcategory"`
	DateAdded   string
	Answers     []Answer
}

type Answer struct {
	ID        uint           `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	QuestionID uint
	Text string `json:"text"`
	IsCorrect bool   `json:"iscorrect"`
}

type QuestionNoCorrectAnswer struct {
	ID          uint
	Text        string `json:"text"`
	Type        string `json:"type"`
	Category    string `json:"category"`
	Subcategory string `json:"subcategory"`
	DateAdded   string
	Answers     []struct {
		ID        uint
		Text      string `json:"text"`
	} `json:"answers"`
}

type AnswerResponse struct {
	IsCorrect bool
	CorrectAnswer []Answer
}
