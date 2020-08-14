package users

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"github.com/Gerard-Szulc/mealsDatabase/database"
	"github.com/Gerard-Szulc/mealsDatabase/interfaces"
	"github.com/Gerard-Szulc/mealsDatabase/utils"
	"os"
	"time"
)

func prepareToken(user *interfaces.User) string {
	jwtKey, exists := os.LookupEnv("JWTKEY")
	if !exists {
		fmt.Println(exists)
	}
	tokenContent := jwt.MapClaims{
		"user_id": user.ID,
		"expiry":  time.Now().Add(time.Minute * 60).Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenContent)
	token, err := jwtToken.SignedString([]byte(jwtKey))
	utils.HandleErr(err)
	return token
}

func prepareResponse(user *interfaces.User, withToken bool) map[string]interface{} {
	responseUser := &interfaces.ResponseUser{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
	var response = map[string]interface{}{"message": "success.response_success"}
	// Add withToken feature to prepare response
	if withToken {
		var token = prepareToken(user)
		response["jwt"] = token
	}
	response["data"] = responseUser
	return response
}

func Login(username string, pass string) map[string]interface{} {
	// Add validation to login
	valid := utils.Validation(
		[]interfaces.Validation{
			{Value: username, Valid: "username"},
			{Value: pass, Valid: "password"},
		})
	if valid {
		// Connect DB
		user := &interfaces.User{}
		if database.DB.Where("username = ? ", username).First(&user).RecordNotFound() {
			return map[string]interface{}{"message": "error.user_not_found"}
		}
		// Verify password
		passErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass))

		if passErr == bcrypt.ErrMismatchedHashAndPassword && passErr != nil {
			return map[string]interface{}{"message": "error.wrong_password"}
		}
		var response = prepareResponse(user, true)
		if !user.Active {
			return map[string]interface{}{"message": "error.account_not_active"}
		}
		return response
	} else {
		return map[string]interface{}{"message": "not valid values"}
	}
}

// Create registration function
func Register(username string, email string, pass string) map[string]interface{} {
	// Add validation to registration
	valid := utils.Validation(
		[]interfaces.Validation{
			{Value: username, Valid: "username"},
			{Value: email, Valid: "email"},
			{Value: pass, Valid: "password"},
		})
	if valid {
		// Create registration logic
		// Connect DB
		generatedPassword := utils.HashAndSalt([]byte(pass))
		user := &interfaces.User{Username: username, Email: email, Password: generatedPassword}
		database.DB.Create(&user)
		var response = prepareResponse(user, true)

		return response
	} else {
		return map[string]interface{}{"message": "error.values_not_valid"}
	}

}

func GetUser(id string) map[string]interface{} {
	// Find and return user
	user := &interfaces.User{}
	if database.DB.Where("id = ? ", id).First(&user).RecordNotFound() {
		return map[string]interface{}{"message": "error.user_not_found"}
	}
	var response = prepareResponse(user, false)
	return response
}
func GetUsers() map[string]interface{} {
	// Find and return user
	var users []interfaces.User
	database.DB.Select("id, username, email").Find(&users)

	responseUsers := []interfaces.ResponseUser{}
	for _, user := range users {
		responseUser := interfaces.ResponseUser{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		}
		responseUsers = append(responseUsers, responseUser)
	}

	container := map[string]interface{}{
		"message": "success_response",
		"users":   responseUsers,
	}
	return container
}
