package main

import (
	"net/http"
	"encoding/json"
	"io"
	"log"
	"fmt"
)

type User struct {
	Name		string `json:"name"`
	Email		string `json:"email"`
	Password	string `json:"password"`
}

func addUser(w http.ResponseWriter, r *http.Request) {
	var user User

	fmt.Println("Adding User...")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal("Err reading request body")
	}

	println(string(body))
	json.Unmarshal(body, &user)
	fmt.Println("Name: ", user.Name)
	fmt.Println("Email: ", user.Email)
	fmt.Println("Password: ", user.Password)
}