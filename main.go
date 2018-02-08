package main

import (
	"fmt"
	"net/http"
	"io"
	"strings"

	"github.com/gorilla/mux"
	"github.com/dgrijalva/jwt-go"
)

const (
    signingKey = "secret"
)

func httpServer() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/token", getToken).Methods("GET")
	r.HandleFunc("/api", getData).Methods("GET")

	return r
}

func createToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
    "username": "foobar",
    "password": "password",
	})
	return (token.SignedString([]byte(signingKey)))
}

func getToken(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getToken is called.")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	token, err := createToken()
	if err != nil {
		fmt.Println("Failed to create token", err)
	}
	fmt.Println("Token is created:", token)

	str := `{"token":"` + token +`"}`
	io.WriteString(w, str)
}

func parseToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(signingKey), nil
	})

  if err == nil && token.Valid {
		claims := token.Claims.(jwt.MapClaims)
		msg := fmt.Sprintf("Hello, %s!!", claims["username"])
		fmt.Println(msg)
	} else {
		fmt.Println("Token is invalid!")
		return err
	}

	return nil
}

func getData(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getData is called.")

	tokenString := r.Header.Get("Authorization")
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	if err := parseToken(tokenString); err != nil {
		fmt.Println("Can not access data.")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	str := `{"data": "12345"}`
	io.WriteString(w, str)
}

func main() {
	fmt.Println("Server started.")
	http.ListenAndServe(":8080", httpServer())
}
