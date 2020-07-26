package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"mealsDatabase/interfaces"
	"mealsDatabase/meals"
	"mealsDatabase/users"
	"mealsDatabase/utils"
	"net/http"
)

type Login struct {
	Username string
	Password string
}

func readBody(r *http.Request) []byte {
	body, err := ioutil.ReadAll(r.Body)
	utils.HandleErr(err)

	return body
}

func apiResponse(call map[string]interface{}, w http.ResponseWriter) {
	if call["message"] == "success.response_success" {
		delete(call, "message")
		resp := call
		json.NewEncoder(w).Encode(resp)
	} else {
		resp := call
		json.NewEncoder(w).Encode(resp)
	}
}

func StartApi() {
	router := mux.NewRouter()
	router.Use(utils.PanicHandler)
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/register", register).Methods("POST")
	router.HandleFunc("/users", getUsers).Methods("GET")
	router.HandleFunc("/users/{id}", getUser).Methods("GET")
	router.HandleFunc("/meals", getMeals).Methods("GET")
	fmt.Println("App is working on port :2137")
	log.Fatal(http.ListenAndServe(":2137", router))
}

func login(w http.ResponseWriter, r *http.Request) {
	body := readBody(r)
	var formattedBody Login
	err := json.Unmarshal(body, &formattedBody)
	utils.HandleErr(err)
	login := users.Login(formattedBody.Username, formattedBody.Password)
	apiResponse(login, w)
}

func register(w http.ResponseWriter, r *http.Request) {
	body := readBody(r)
	var formattedBody interfaces.Register
	err := json.Unmarshal(body, &formattedBody)
	utils.HandleErr(err)
	register := users.Register(formattedBody.Username, formattedBody.Email, formattedBody.Password)
	apiResponse(register, w)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["id"]
	//auth := r.Header.Get("Authorization")
	if !utils.ValidateRequestToken(r) {
		apiResponse(map[string]interface{}{"message": "error:token_not_valid"}, w)
		return
	}
	user := users.GetUser(userId)
	apiResponse(user, w)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	if !utils.ValidateRequestToken(r) {
		apiResponse(map[string]interface{}{"message": "error:token_not_valid"}, w)
		return
	}
	user := users.GetUsers()
	apiResponse(user, w)
}

func getMeals(w http.ResponseWriter, r *http.Request) {
	//auth := r.Header.Get("Authorization")
	if !utils.ValidateRequestToken(r) {
		apiResponse(map[string]interface{}{"message": "error:token_not_valid"}, w)
		return
	}
	responseMeals := meals.GetMeals()
	apiResponse(responseMeals, w)
}
