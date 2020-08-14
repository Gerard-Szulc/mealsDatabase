package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Gerard-Szulc/mealsDatabase/ingredients"
	"github.com/Gerard-Szulc/mealsDatabase/interfaces"
	"github.com/Gerard-Szulc/mealsDatabase/meals"
	"github.com/Gerard-Szulc/mealsDatabase/users"
	"github.com/Gerard-Szulc/mealsDatabase/utils"
	"github.com/gorilla/mux"
)

//Login form
type Login struct {
	Username string
	Password string
}

func readBody(r *http.Request) []byte {
	body, err := ioutil.ReadAll(r.Body)
	utils.HandleErr(err)

	return body
}

//StartAPI starts routing
func StartAPI() {
	router := mux.NewRouter()
	router.Use(utils.PanicHandler)
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/register", register).Methods("POST")
	router.HandleFunc("/users", getUsers).Methods("GET")
	router.HandleFunc("/users/{id}", getUser).Methods("GET")
	router.HandleFunc("/users/{id}/meals", getUserMeals).Methods("GET")
	router.HandleFunc("/meals", getMeals).Methods("GET")
	router.HandleFunc("/meals/{id}", getMeal).Methods("GET")
	router.HandleFunc("/meals", addMeal).Methods("POST")
	router.HandleFunc("/ingredients", ingredients.GetIngredientsRoute).Methods("GET")
	router.HandleFunc("/ingredients/{id}", ingredients.GetIngredientRoute).Methods("GET")

	fmt.Println("App is working on port :2137")
	log.Fatal(http.ListenAndServe(":2137", router))
}

func login(w http.ResponseWriter, r *http.Request) {
	body := readBody(r)
	var formattedBody Login
	err := json.Unmarshal(body, &formattedBody)
	utils.HandleErr(err)
	login := users.Login(formattedBody.Username, formattedBody.Password)
	utils.ApiResponse(login, w)
}

func register(w http.ResponseWriter, r *http.Request) {
	body := readBody(r)
	var formattedBody interfaces.Register
	err := json.Unmarshal(body, &formattedBody)
	utils.HandleErr(err)
	register := users.Register(formattedBody.Username, formattedBody.Email, formattedBody.Password)
	utils.ApiResponse(register, w)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]
	//auth := r.Header.Get("Authorization")
	if !utils.ValidateRequestToken(r) {
		utils.ApiResponse(map[string]interface{}{
			"message": "error:token_not_valid",
		}, w)
		return
	}
	user := users.GetUser(userID)
	utils.ApiResponse(user, w)
}
func getUserMeals(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]
	//auth := r.Header.Get("Authorization")
	if !utils.ValidateRequestToken(r) {
		utils.ApiResponse(map[string]interface{}{
			"message": "error:token_not_valid",
		}, w)
		return
	}
	user := meals.GetUserMeals(userID)
	utils.ApiResponse(user, w)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	if !utils.ValidateRequestToken(r) {
		utils.ApiResponse(map[string]interface{}{"message": "error:token_not_valid"}, w)
		return
	}
	user := users.GetUsers()
	utils.ApiResponse(user, w)
}

func getMeals(w http.ResponseWriter, r *http.Request) {

	if !utils.ValidateRequestToken(r) {
		utils.ApiResponse(map[string]interface{}{"message": "error:token_not_valid"}, w)
		return
	}

	q := r.URL.Query()
	label := q.Get("find")
	if label != "" {
		responseMeals := meals.SearchMeals(label)
		utils.ApiResponse(responseMeals, w)
		return
	}

	responseMeals := meals.GetMeals()
	utils.ApiResponse(responseMeals, w)
}
func getMeal(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	mealId := vars["id"]

	if !utils.ValidateRequestToken(r) {
		utils.ApiResponse(map[string]interface{}{
			"message": "error:token_not_valid",
			"code":    400,
		}, w)
		return
	}
	responseMeal := meals.GetMeal(mealId)
	utils.ApiResponse(responseMeal, w)
}
func searchMeals(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	label := vars["label"]

	if !utils.ValidateRequestToken(r) {
		utils.ApiResponse(map[string]interface{}{
			"message": "error:token_not_valid",
			"code":    400,
		}, w)
		return
	}
	responseMeals := meals.SearchMeals(label)
	utils.ApiResponse(responseMeals, w)
}

func addMeal(w http.ResponseWriter, r *http.Request) {
	//auth := r.Header.Get("Authorization")
	body, _ := ioutil.ReadAll(r.Body)
	var meal interfaces.Meal
	err := json.Unmarshal(body, &meal)
	utils.HandleErr(err)

	fmt.Println(meal)
	if !utils.ValidateRequestToken(r) {
		utils.ApiResponse(map[string]interface{}{
			"message": "error:token_not_valid",
			"code":    400,
		}, w)
		return
	}
	responseMeal := meals.AddMeal(meal)
	utils.ApiResponse(responseMeal, w)
}
