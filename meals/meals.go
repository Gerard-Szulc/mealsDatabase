package meals

import (
	"mealsDatabase/database"
	"mealsDatabase/interfaces"
)

func GetMeals() map[string]interface{} {
	// Find and return user
	var meals []interfaces.Meal
	database.DB.Find(&meals)

	container := map[string]interface{}{
		"message": "success_response",
		"meals":   meals,
	}
	return container
}
