package models

type Category struct {
	Id string	`json:"id"`
	CategoryName string   `json:"Category"`
	SubCategories []Subcategory `json:"SubCategories"`
}

type Subcategory struct {
	SubCategoryName string `json:"SubCategoryName"`
	URLPrefix string `json:"URLPrefix"`
}