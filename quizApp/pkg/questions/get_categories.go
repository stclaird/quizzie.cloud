package questions

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/stclaird/quizzie.cloud/pkg/common/models"
)

func (h handler) GetCategories(ctx *gin.Context) {
	var questions []models.Question

	err := h.DB.Model(&models.Question{}).Preload("Answers").Find(&questions).Error
	if err != nil {
		fmt.Printf("error %s", err)
	}

	Categories := make(map[string]*models.Category)
	subCategories := make(map[string]models.Subcategory)

	//Populate a map of all Categories and a second map of all subcategories
	for k, v := range questions {
		var category models.Category
		category.Id = strconv.Itoa(k)
		category.CategoryName = v.Category
		Categories[v.Category] = &category

		var subCategory models.Subcategory
		subCategory.SubCategoryName = v.Subcategory
		subCategory.URLPrefix = fmt.Sprintf("%s-%s", v.Category, v.Subcategory)
		subCategories[subCategory.URLPrefix] = subCategory
	}

	//Apply all subcategories to appropriate category
	for _, v := range subCategories {
		splt := strings.Split(v.URLPrefix, "-")
		cat := splt[0]
		for _, value := range Categories {
			if value.CategoryName == cat {
				Categories[cat].SubCategories = append(Categories[cat].SubCategories, v)
			}
		}
	}

	var response []*models.Category
	for _, v := range Categories {
		response = append(response, v)
	}

	ctx.JSON(http.StatusOK, &response)
}
