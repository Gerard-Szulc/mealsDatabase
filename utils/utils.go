package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Gerard-Szulc/mealsDatabase/interfaces"
	"github.com/golang-jwt/jwt"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func HandleErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func HashAndSalt(pass []byte) string {
	hashed, err := bcrypt.GenerateFromPassword(pass, bcrypt.MinCost)
	HandleErr(err)
	return string(hashed)
}

func init() {
	// loads values from .env into the system
	env := os.Getenv("FINT_ENV")
	fmt.Println(env)

	if "" == env {
		env = "development"
	}

	godotenv.Load(".env." + env + ".local")
	fmt.Println(".env." + env + ".local")

	if "test" != env {
		godotenv.Load(".env.local")
	}
	godotenv.Load(".env." + env)
	fmt.Println(".env." + env)

	err := godotenv.Load()
	HandleErr(err)
}

func Validation(values []interfaces.Validation) bool {
	username := regexp.MustCompile("^([A-Za-z0-9]{5,})+$")
	email := regexp.MustCompile("^[A-Za-z0-9]+[@]+[A-Za-z0-9]+[.]+[A-Za-z]+$")
	for i := 0; i < len(values); i++ {
		switch values[i].Valid {
		case "username":
			if !username.MatchString(values[i].Value) {
				return false
			}
		case "email":
			if !email.MatchString(values[i].Value) {
				return false
			}
		case "password":
			if len(values[i].Value) < 5 {
				return false
			}
		}
	}
	return true
}

func PanicHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			error := recover()
			if error != nil {
				log.Println(error)

				resp := interfaces.ErrResponse{Message: "Internal server error"}
				json.NewEncoder(w).Encode(resp)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func ValidateToken(id string, jwtToken string) bool {
	jwtKey, exists := os.LookupEnv("JWTKEY")
	if !exists {
		fmt.Println(exists)
	}
	cleanJWT := strings.Replace(jwtToken, "Bearer ", "", -1)
	tokenData := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(cleanJWT, tokenData, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	HandleErr(err)
	var userId, _ = strconv.ParseFloat(id, 8)
	if token.Valid && tokenData["user_id"] == userId {
		return true
	} else {
		return false
	}

}

func ValidateRequestToken(r *http.Request) bool {
	jwtKey, exists := os.LookupEnv("JWTKEY")
	if !exists {
		fmt.Println(exists)
		return false
	}
	jwtToken := r.Header.Get("Authorization")
	cleanJWT := strings.Replace(jwtToken, "Bearer ", "", -1)
	cleanJWTHeader := strings.Split(cleanJWT, ".")[0]
	cleanJWTPayload := strings.Split(cleanJWT, ".")[1]
	cleanJWTSecret := strings.Split(cleanJWT, ".")[2]
	_, err := jwt.DecodeSegment(cleanJWTHeader)
	if err != nil {
		if _, ok := err.(base64.CorruptInputError); ok {
			panic("\nbase64 input is corrupt, check service Key")
		}
		panic(err)
	}
	_, err = jwt.DecodeSegment(cleanJWTPayload)
	if err != nil {
		if _, ok := err.(base64.CorruptInputError); ok {
			panic("\nbase64 input is corrupt, check service Key")
		}
		panic(err)
	}
	_, err = jwt.DecodeSegment(cleanJWTSecret)
	if err != nil {
		if _, ok := err.(base64.CorruptInputError); ok {
			panic("\nbase64 input is corrupt, check service Key")
		}
		panic(err)
	}

	tokenData := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(cleanJWT, tokenData, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	HandleErr(err)
	//HandleErrRequest(err)

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

func ApiResponse(call map[string]interface{}, w http.ResponseWriter) {
	str := fmt.Sprintf("%v", call["message"])
	if strings.Contains(str, "success") {
		delete(call, "message")
		resp := call
		json.NewEncoder(w).Encode(resp)
	} else {
		resp := call
		w.WriteHeader(call["code"].(int))
		json.NewEncoder(w).Encode(resp)
	}
}
