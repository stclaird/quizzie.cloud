package models

import "gorm.io/gorm"

type Question struct {
	gorm.Model
	Text        string `json:"text"`
	Type        string `json:"type"`
	Category    string `json:"category"`
	Subcategory string `json:"subcategory"`
	DateAdded   string
	Answers     []Answer
}

type Answer struct {
	gorm.Model
	QuestionID uint
	Text string `json:"text"`
	IsCorrect bool   `json:"iscorrect"`
}