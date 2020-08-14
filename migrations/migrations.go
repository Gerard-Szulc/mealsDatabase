package migrations

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/Gerard-Szulc/mealsDatabase/database"
	"github.com/Gerard-Szulc/mealsDatabase/interfaces"
	"github.com/Gerard-Szulc/mealsDatabase/utils"
)

func createAccounts() {
	users := &[1]interfaces.User{
		{Username: "Gerard", Email: "gerszulc05@gmail.com"},
	}
	for i := 0; i < len(users); i++ {
		generatedPassword := utils.HashAndSalt([]byte(users[i].Username))
		user := &interfaces.User{Username: users[i].Username, Email: users[i].Email, Password: generatedPassword}
		database.DB.Create(&user)
	}
}

func Migrate() {
	User := &interfaces.User{}
	Meal := &interfaces.Meal{}
	Ingredient := &interfaces.Ingredient{}
	IngredientCategory := &interfaces.IngredientTypeCategory{}
	MealCategory := &interfaces.MealTypeCategory{}
	database.DB.AutoMigrate(&User, &Meal, &Ingredient, &IngredientCategory, &MealCategory)

	createAccounts()
}
