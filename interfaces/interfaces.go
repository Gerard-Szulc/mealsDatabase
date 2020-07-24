package interfaces

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Username string
	Email    string
	Password string
	Active   bool    `gorm:"default:false"`
	Meals    []*Meal `gorm:"many2many:user_meals;"`
}

type ResponseUser struct {
	ID       uint
	Username string
	Email    string
}

type Validation struct {
	Value string
	Valid string
}

type Register struct {
	Username string
	Email    string
	Password string
}

type ErrResponse struct {
	Message string
}

type Meal struct {
	gorm.Model
	Name        string
	Label       string
	Ingredients []*Ingredient `gorm:"many2many:meal_ingredients;"`
	Description string
	Categories  []*MealTypeCategory `gorm:"many2many:meal_categories;"`
	Users       []*User             `gorm:"many2many:user_meals;"`
}

type MealTypeCategory struct {
	gorm.Model
	Name  string
	Label string
	Meals []*Meal `gorm:"many2many:meal_categories;"`
}

type Ingredient struct {
	gorm.Model
	Name                 string
	Label                string
	Calories             float64
	IngredientCategories []*IngredientTypeCategory `gorm:"many2many:ingredient_categories;"`
	Meals                []*Meal                   `gorm:"many2many:meal_ingredients;"`
}

type IngredientTypeCategory struct {
	gorm.Model
	Name        string
	Label       string
	Ingredients []*Ingredient `gorm:"many2many:ingredient_categories;"`
}
