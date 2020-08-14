package ingredients

import (
	"github.com/Gerard-Szulc/mealsDatabase/database"
	"github.com/Gerard-Szulc/mealsDatabase/interfaces"
)

func GetIngredients() map[string]interface{} {

	var ingredients []interfaces.Ingredient

	database.DB.Find(&ingredients)
	//TODO add users property when user has admin privileges

	container := map[string]interface{}{
		"message": "success_response",
		"meals":   ingredients,
	}
	return container
}

func GetIngredient(id string) map[string]interface{} {

	var ingredient interfaces.Ingredient

	if database.DB.Where("id = ? ", id).Preload("Meals").Preload("Meals.Ingredients").Preload("Meals.Categories").Find(&ingredient).RecordNotFound() {
		return map[string]interface{}{
			"message": "error.ingredient_not_found",
			"code":    404,
		}
	}

	container := map[string]interface{}{
		"data":    ingredient,
		"message": "success_ingredient_found",
	}
	return container
}
