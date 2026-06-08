package fake

import (
	"forum/database"
	"forum/models"
)

func GetAllCategories() []models.Category {
	return database.GetAllCategories()
}

func GetCategoryById(id int) (models.Category, bool) {
	return database.GetCategoryByID(id)
}
