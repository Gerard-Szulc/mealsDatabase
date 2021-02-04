package database

import (
	"fmt"
	"github.com/Gerard-Szulc/mealsDatabase/utils"
	"github.com/jinzhu/gorm"
	"os"
)

var DB *gorm.DB

func InitDatabase() {
	dbhost, exists := os.LookupEnv("DBHOST")
	if !exists {
		fmt.Println(exists)
	}

	dbport, exists := os.LookupEnv("DBPORT")
	if !exists {
		fmt.Println(exists)
	}
	user, exists := os.LookupEnv("DBUSER")
	if !exists {
		fmt.Println(exists)
	}
	password, exists := os.LookupEnv("DBPASSWORD")
	if !exists {
		fmt.Println(exists)
	}
	dbname, exists := os.LookupEnv("DBNAME")
	if !exists {
		fmt.Println(exists)
	}
	dbargs := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", dbhost, dbport, user, dbname, password)
	db, err := gorm.Open("postgres", dbargs)
	utils.HandleErr(err)
	db.DB().SetMaxIdleConns(20)
	db.DB().SetMaxOpenConns(200)
	DB = db
}
