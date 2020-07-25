package api

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"mealsDatabase/interfaces"
	"mealsDatabase/meals"
	"mealsDatabase/users"
	"mealsDatabase/utils"
	"net/http"
	"os"
	"strings"
	"time"
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
	router.HandleFunc("/user/{id}", getUser).Methods("GET")
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
	auth := r.Header.Get("Authorization")
	user := users.GetUser(userId, auth)
	apiResponse(user, w)
}
func getMeals(w http.ResponseWriter, r *http.Request) {
	//auth := r.Header.Get("Authorization")
	if !ValidateRequestToken(r) {
		apiResponse(map[string]interface{}{"message": "error:token_not_valid"}, w)
		return
	}
	responseMeals := meals.GetMeals()
	apiResponse(responseMeals, w)
}

func ValidateRequestToken(r *http.Request) bool {
	jwtKey, exists := os.LookupEnv("JWTKEY")
	if !exists {
		fmt.Println(exists)
	}
	jwtToken := r.Header.Get("Authorization")
	cleanJWT := strings.Replace(jwtToken, "Bearer ", "", -1)
	tokenData := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(cleanJWT, tokenData, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	utils.HandleErr(err)

	now := time.Now()
	expiry := tokenData["expiry"].(float64)

	expired := now.After(time.Unix(int64(expiry), 0))
	if expired {
		return false
	}
	if !token.Valid {
		return false
	}

	return true
}
