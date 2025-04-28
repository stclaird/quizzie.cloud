package models

type Category struct {
	Id            string        `json:"id"`
	CategoryName  string        `json:"category"`
	SubCategories []Subcategory `json:"subcategories"`
}

type Subcategory struct {
	SubCategoryName string `json:"subcategoryname"`
	URLPrefix       string `json:"urlprefix"`
}
