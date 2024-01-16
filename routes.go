package main

import (
	"course-app-with-jwt/auth"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
	insertQuery  = "INSERT INTO users (username, password, role) VALUES (?, ?, ?)"
	selectQuery  = "SELECT * From users Where username = ?"
	insertCourse = "INSERT INTO courses (title, description, price, imageLink, published) VALUES (?, ?, ?, ?, ?)"
	selectCourse = "SELECT id from courses where title = ?"
	admin        = "admin"
)

type credentials struct {
	UserName string `json:"username"`
	PassWord string `json:"password"`
	Role     string `json:"role,omitempty"`
}
type response struct {
	Message  string `json:"message"`
	Token    string `json:"token,omitempty"`
	CourseId int    `json:"CourseId,omitempty"`
}

// Course struct representing course data
type Course struct {
	Id          int     `json:"id,omitempty"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	ImageLink   string  `json:"imageLink"`
	Published   bool    `json:"published"`
}

func adminSignup(w http.ResponseWriter, r *http.Request) {
	pasrseBody := r.Body
	bodyByte, err := io.ReadAll(pasrseBody)

	if err != nil {
		fmt.Println("", err)
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println(string(bodyByte))

	cred := credentials{}
	err = json.Unmarshal([]byte(bodyByte), &cred)

	if err != nil {
		fmt.Println("Json Unmarshal :", err)
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	result, err := db.Exec(insertQuery, cred.UserName, cred.PassWord, admin)

	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	token := auth.GenerateJwt(cred.UserName)

	res := response{
		Message: "Admin created successfully",
		Token:   token,
	}
	jsonRes, err := json.Marshal(res)

	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")

	w.Write(jsonRes)
	fmt.Printf("Inserted row with ID: %d\n", lastInsertID)
}

func adminLogin(w http.ResponseWriter, r *http.Request) {
	pasrseBody := r.Body
	bodyByte, err := io.ReadAll(pasrseBody)

	if err != nil {
		fmt.Println("", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println(string(bodyByte))

	cred := credentials{}
	err = json.Unmarshal([]byte(bodyByte), &cred)

	if err != nil {
		fmt.Println("Json Unmarshal :", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	row := db.QueryRow(selectQuery, cred.UserName)

	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
	}

	dbCred := credentials{}
	err = row.Scan(&dbCred.UserName, &dbCred.PassWord, &dbCred.Role)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if cred.PassWord == dbCred.PassWord && dbCred.Role == "admin" {
		fmt.Println("admin logged in")
	}
	token := auth.GenerateJwt(cred.UserName)

	res := response{
		Message: "Logged in successfully",
		Token:   token,
	}
	jsonRes, err := json.Marshal(res)

	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")

	w.Write(jsonRes)
}

func createCourse(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	err := auth.ValidateToken(token)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pasrseBody := r.Body
	bodyByte, err := io.ReadAll(pasrseBody)

	if err != nil {
		fmt.Println("", err)
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println(string(bodyByte))

	newCourse := Course{}
	err = json.Unmarshal([]byte(bodyByte), &newCourse)

	if err != nil {
		fmt.Println("Json Unmarshal :", err)
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = db.Exec(insertCourse, newCourse.Title, newCourse.Description, newCourse.Price, newCourse.ImageLink, newCourse.Published)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	row := db.QueryRow(selectCourse, newCourse.Title)

	var course Course
	err = row.Scan(&course.Id)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	res := response{
		Message:  "Course created successfully",
		CourseId: course.Id,
	}
	jsonRes, err := json.Marshal(res)

	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")

	w.Write(jsonRes)

}
