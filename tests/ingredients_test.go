package mealdbtests

import (
	"testing"
	"time"

	"github.com/Gerard-Szulc/mealsDatabase/ingredients"
	"github.com/Gerard-Szulc/mealsDatabase/interfaces"
	mocket "github.com/Selvatico/go-mocket"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

func SetupTests() *gorm.DB {
	mocket.Catcher.Register()
	mocket.Catcher.Logging = true
	// GORM
	db, err := gorm.Open(mocket.DriverName, "example")
	db = db
	if err != nil {
		panic(err)
	}
	return db
}

func TestGetIngredient(t *testing.T) {
	db = SetupTests()

	type args struct {
		id string
		DB *gorm.DB
	}

	commonReply := []map[string]interface{}{{
		"id":                    uint(2),
		"name":                  "",
		"label":                 "",
		"calories":              nil,
		"ingredient_categories": nil,
		"meals":                 nil,
		"created_at":            time.Time{},
		"updated_at":            time.Time{},
		"deleted_at":            time.Time{},
	}}

	t.Run("Checks length of rerurned map is 2", func(t *testing.T) {
		mocket.Catcher.Logging = true

		mocket.Catcher.Reset().NewMock().WithID(2).WithReply(commonReply)

		result := ingredients.GetIngredient("2", db)

		if len(result) != 2 {
			t.Errorf("Returned sets is not equal to 2. Received %d", len(result))
		}
	})

	t.Run("Check id returned from database is the same as on mock", func(t *testing.T) {
		mocket.Catcher.Logging = true
		mocket.Catcher.Reset().NewMock().WithID(2).WithReply(commonReply)

		result := ingredients.GetIngredient("2", db)
		data := result["data"]

		ingre := data.(interfaces.Ingredient)

		if ingre.ID != commonReply[0]["id"] {
			t.Errorf("Returned ingredient id %d is not equal to id of commonReply %d .  ", ingre.ID, commonReply[0]["id"])
		}
	})
}
