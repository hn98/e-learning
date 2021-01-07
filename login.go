package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"context"
	"go.mongodb.org/mongo-driver/bson"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func StudentLogin(w http.ResponseWriter, r *http.Request) {
	var user LoginRequest
	err := HandleRequest(w, r, &user)
	fmt.Println(user)

	if err != nil {
		return
	}

	resp := FindStudent(user.Username, user.Password)
	Respond(w, resp)
}

func FindStudent(username, password string) map[string]interface{} {
	var student Student

	studentCollection := dbClient.Database("learning").Collection("Students")
	err := studentCollection.FindOne(context.Background(), bson.M{"username": username}).Decode(&student)
	if err != nil {
		var resp = map[string]interface{}{"status": false, "message": "Username not found"}
		return resp
	}
	expiresAt := time.Now().Add(time.Minute * 1000).Unix()

	errf := bcrypt.CompareHashAndPassword([]byte(student.Password), []byte(password))

	if errf != nil && errf == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		var resp = map[string]interface{}{"status": false, "message": "Invalid login credentials. Please try again"}
		return resp
	}

	tk := &Token{
		ID:       student.ID.Hex(),
		Username: student.Username,
		Role:     "Student",
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)

	tokenString, error := token.SignedString([]byte("secret"))
	if error != nil {
		fmt.Println(error)
	}

	var resp = map[string]interface{}{"status": false, "message": "logged in"}
	resp["token"] = tokenString //Store the token in the response
	return resp
}

func InstructorLogin(w http.ResponseWriter, r *http.Request) {
	var user LoginRequest
	err := HandleRequest(w, r, &user)
	if err != nil {
		return
	}

	resp := FindInstructor(user.Username, user.Password)
	Respond(w, resp)
}

func FindInstructor(username, password string) map[string]interface{} {
	var instructor Instructor

	instructorCollection := dbClient.Database("learning").Collection("Instructors")
	err := instructorCollection.FindOne(context.Background(), bson.M{"username": username}).Decode(&instructor)

	if err != nil {
		var resp = map[string]interface{}{"status": false, "message": "Username not found"}
		return resp
	}
	expiresAt := time.Now().Add(time.Minute * 1000).Unix()

	errf := bcrypt.CompareHashAndPassword([]byte(instructor.Password), []byte(password))

	if errf != nil && errf == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		var resp = map[string]interface{}{"status": false, "message": "Invalid login credentials. Please try again"}
		return resp
	}

	tk := &Token{
		ID:       instructor.ID.Hex(),
		Username: instructor.Username,
		Role:     "Instructor",
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)

	tokenString, error := token.SignedString([]byte("secret"))
	if error != nil {
		fmt.Println(error)
	}

	var resp = map[string]interface{}{"status": false, "message": "logged in"}
	resp["token"] = tokenString //Store the token in the response
	return resp
}

func JwtVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var header = r.Header.Get("x-access-token") //Grab the token from the header

		header = strings.TrimSpace(header)

		if header == "" {
			header = r.Header.Get("Authorization")
			splitToken := strings.Split(header, "Bearer ")
			header = splitToken[1]
		}

		if header == "" {
			//Token is missing, returns with error code 403 Unauthorized
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(Exception{Message: "Missing auth token"})
			return
		}
		tk := &Token{}

		_, err := jwt.ParseWithClaims(header, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})

		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(Exception{Message: err.Error()})
			return
		}

		ctx := context.WithValue(r.Context(), "user", tk)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
