package fake

import "forum/models"

func GetAllCategories() []models.Category { 
	var allCategories []models.Category = []models.Category{
		{
			ID: 1,
			Name: "Général",
		},
		{
			ID: 2,
			Name: "Jeux vidéos",
		},
		{
			ID: 3,
			Name: "Manga/Animé",
		},
		{
			ID: 4,
			Name: "Actualités",
		},
	}
	return allCategories
}

func GetCategoryById(id int) (models.Category, bool) {
	var allCategories []models.Category = []models.Category{
		{
			ID: 1,
			Name: "Général",
		},
		{
			ID: 2,
			Name: "Jeux vidéos",
		},
		{
			ID: 3,
			Name: "Manga/Animé",
		},
		{
			ID: 4,
			Name: "Actualités",
		},
	}

	for _, cat := range allCategories {
		if cat.ID == id {
			return cat, true
		}
	}

	return models.Category{}, false
}