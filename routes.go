package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"course-app-with-jwt/auth"
)

const (
	insertQuery = "INSERT INTO users (username, password) VALUES (?, ?)"
)

type credentials struct {
	UserName string `json:"username"`
	PassWord string `json:"password"`
}

func adminSignup(w http.ResponseWriter, r *http.Request) {
	pasrseBody := r.Body
	bodyByte, err := io.ReadAll(pasrseBody)

    if err != nil {
		fmt.Println("", err)
	}

	fmt.Println(string(bodyByte))

	cred := credentials{}
	err = json.Unmarshal([]byte(bodyByte), &cred)

	if err != nil {
		fmt.Println("Json Unmarshal :", err)
	}
	result, err := db.Exec(insertQuery, cred.UserName, cred.PassWord)
    
	if err != nil {
		log.Fatal(err)
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
    token := auth.GenerateJwt(cred.UserName)
	err = auth.ValidateToken(token)
	if err != nil{
		fmt.Println("Token validation Failed",err)
	}
    
	fmt.Printf("Inserted row with ID: %d\n", lastInsertID)

}
