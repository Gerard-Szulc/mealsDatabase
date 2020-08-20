package ingredients

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/Gerard-Szulc/mealsDatabase/database"
	"github.com/Gerard-Szulc/mealsDatabase/interfaces"
	"github.com/Gerard-Szulc/mealsDatabase/utils"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

//GetIngredients gets all ingredients from database
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

//GetIngredient gets single ingredient with associated meals
func GetIngredient(id string, DB *gorm.DB) map[string]interface{} {

	var ingredient interfaces.Ingredient

	if DB.Where("id = ? ", id).Preload("Meals").Preload("Meals.Ingredients").Preload("Meals.Categories").Find(&ingredient).RecordNotFound() {
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

//GetIngredientsRoute is for getting all or find list of ingredients by its labels
func GetIngredientsRoute(w http.ResponseWriter, r *http.Request) {
	if !utils.ValidateRequestToken(r) {
		utils.ApiResponse(map[string]interface{}{"message": "error:token_not_valid"}, w)
		return
	}

	q := r.URL.Query()
	label := q.Get("search")
	if label != "" {
		responseIngredients := searchIngredients(label)
		utils.ApiResponse(responseIngredients, w)
		return
	}

	responseIngredients := GetIngredients()
	utils.ApiResponse(responseIngredients, w)
}

//GetIngredientRoute siple getting single ingredient
func GetIngredientRoute(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	mealID := vars["id"]

	if !utils.ValidateRequestToken(r) {
		utils.ApiResponse(map[string]interface{}{
			"message": "error:token_not_valid",
			"code":    400,
		}, w)
		return
	}
	responseIngredient := GetIngredient(mealID, database.DB)
	utils.ApiResponse(responseIngredient, w)
}

func searchIngredients(label string) map[string]interface{} {

	var ingredients []interfaces.Ingredient

	if database.DB.Where("label LIKE ? ", label+"%").Preload("Meals").Find(&ingredients).RecordNotFound() {
		return map[string]interface{}{
			"message": "error.ingredient_not_found",
			"code":    404,
		}
	}

	container := map[string]interface{}{
		"data":    ingredients,
		"message": "success_ingredient_found",
	}
	return container
}

func AddIngredientRoute(w http.ResponseWriter, r *http.Request) {
	//auth := r.Header.Get("Authorization")
	body, _ := ioutil.ReadAll(r.Body)
	var ingredient interfaces.Ingredient
	err := json.Unmarshal(body, &ingredient)
	utils.HandleErr(err)

	if !utils.ValidateRequestToken(r) {
		utils.ApiResponse(map[string]interface{}{
			"message": "error:token_not_valid",
			"code":    400,
		}, w)
		return
	}
	responseIngredient := AddIngredient(ingredient, w)
	utils.ApiResponse(responseIngredient, w)
}

func AddIngredient(ingredient interfaces.Ingredient, w http.ResponseWriter) map[string]interface{} {

	isAllowed := database.DB.NewRecord(ingredient)
	if !isAllowed {
		container := map[string]interface{}{
			"message": "error_ingredient_id_exists",
			"code":    403,
		}
		utils.ApiResponse(container, w)
		utils.HandleErr(errors.New("error_ingredient_id_exists"))
	}

	database.DB.Create(&ingredient)
	container := map[string]interface{}{
		"data":    ingredient,
		"message": "success_ingredient_added",
	}
	return container
}
