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
	fmt.Printf("questions: %v\n", questions)
	//
	//Use maps to ensure unique categories and subcategories
	//
	//Keyed on category name to ensure uniqueness
	//Value is a pointer to the category struct so we can update it
	//with subcategories as we find them
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
	//by matching on category name.
	//This builds out the subcategories slice in each category
	//with all the subcategories that belong to that category
	//based on the questions we have
	//This assumes that the category name is the first part of the
	//subcat URL prefix before the hyphen
	//e.g. "science-fiction-hitchhikersguide" belongs to category "science-fiction"
	for _, v := range subCategories {
		split := strings.Split(v.URLPrefix, "-")
		cat := split[0]
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
