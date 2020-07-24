package main

import (
	"mealsDatabase/api"
	"mealsDatabase/database"
	"mealsDatabase/migrations"
	"os"
)

func start () {
	database.InitDatabase()
	api.StartApi()
}

func main() {
	argsWithProg := os.Args
	if len(argsWithProg) <= 1 {
		start()
	} else {
		switch argsWithProg[1] {
		case "migrate":
			{
				database.InitDatabase()
				migrations.Migrate()
				return
			}
		case "start":
			{
				start()
			}
		}
	}
}
