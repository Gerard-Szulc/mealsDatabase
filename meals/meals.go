package meals

import (
	"fmt"
	"mealsDatabase/database"
	"mealsDatabase/interfaces"
)

func GetMeals() map[string]interface{} {
	// Find and return user
		meals := []interfaces.Meal{}
		database.DB.Find(&meals)

		fmt.Println(meals)
		container := map[string]interface{}{
			"message": "success_response",
			"meals": meals,
		}
		return container
}