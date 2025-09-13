package questions

import (
	"github.com/stclaird/quizzie.cloud/pkg/common/models"
)

func (h handler) CreateQuestion(question models.Question) (uint, error) {

	result := h.DB.Create(&question)
	return question.ID, result.Error
}
