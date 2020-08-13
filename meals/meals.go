package meals

import (
	"mealsDatabase/database"
	"mealsDatabase/interfaces"
)

func GetMeals() map[string]interface{} {

	var meals []interfaces.Meal

	database.DB.Preload("Ingredients").Preload("Categories").Find(&meals)
	//TODO add users property when user has admin privileges

	container := map[string]interface{}{
		"message": "success_response",
		"meals":   meals,
	}
	return container
}

func GetMeal(id string) map[string]interface{} {

	var meal interfaces.Meal

	if database.DB.Where("id = ? ", id).Preload("Categories").Preload("Ingredients").Find(&meal).RecordNotFound() {
		return map[string]interface{}{
			"message": "error.meal_not_found",
			"code":    404,
		}
	}

	container := map[string]interface{}{
		"data":    meal,
		"message": "success_meal_found",
	}
	return container
}
func SearchMeals(label string) map[string]interface{} {

	var meals []interfaces.Meal

	if database.DB.Where("label LIKE ? ", label+"%").Preload("Categories").Preload("Ingredients").Find(&meals).RecordNotFound() {
		return map[string]interface{}{
			"message": "error.meal_not_found",
			"code":    404,
		}
	}

	container := map[string]interface{}{
		"data":    meals,
		"message": "success_meal_found",
	}
	return container
}

func GetUserMeals(id string) map[string]interface{} {

	user := &interfaces.User{}

	if database.DB.Where("id = ? ", id).Preload("Meals.Categories").Preload("Meals.Ingredients").Preload("Meals").Find(&user).RecordNotFound() {
		return map[string]interface{}{
			"message": "error.user_not_found",
			"code":    404,
		}
	}

	container := map[string]interface{}{
		"data":    user.Meals,
		"message": "success_users_meal_found",
	}
	return container
}

func AddMeal(meal interfaces.Meal) map[string]interface{} {

	database.DB.NewRecord(meal)
	database.DB.Create(&meal)
	container := map[string]interface{}{
		"data":    meal,
		"message": "success_meal_added",
	}
	return container
}
