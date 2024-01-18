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
	insertQuery    = "INSERT INTO users (username, password, role) VALUES (?, ?, ?)"
	selectQuery    = "SELECT username, password, role From users Where username = ?"
	insertCourse   = "INSERT INTO courses (title, description, price, imageLink, published) VALUES (?, ?, ?, ?, ?)"
	selectCourseId = "SELECT id from courses where title = ?"
	updateCourse   = "UPDATE courses SET title = ?, description = ?, price = ?, imageLink = ?, published = ? WHERE id = ? "
	getCourses     = "SELECT * from courses"
	admin          = "admin"

	//user
	user           = "user"

	//common
	getUsrAdminId  = "Select id from users where username = ? And Password = ?"   
)

// credentials struct representing credentials data
type credentials struct {
	Id       int    `json:"id,omitempty"`
	UserName string `json:"username"`
	PassWord string `json:"password"`
	Role     string `json:"role,omitempty"`
}

// response struct representing response data
type response struct {
	Message  string `json:"message"`
	Token    string `json:"token,omitempty"`
	CourseId int    `json:"CourseId,omitempty"`
}

// course struct representing course data
type course struct {
	Id          int     `json:"id,omitempty"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	ImageLink   string  `json:"imageLink"`
	Published   bool    `json:"published"`
}

// adminSignup creates an account for an admin.
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
	row,_ := db.Query(getUsrAdminId,cred.UserName,cred.PassWord)
	row.Scan(&cred.Id)
	token := auth.GenerateJwt(cred.Id)

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

// adminLogin handles the validation of login credentials.
// On successful validation, it returns a token.
// On failed validation, it returns a bad request.
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
	ro,_ := db.Query(getUsrAdminId,dbCred.UserName,dbCred.PassWord)
	ro.Scan(&cred.Id)
	token := auth.GenerateJwt(cred.Id)

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

// createCourse creates a new course.
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
		fmt.Println("parse error", err)
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println(string(bodyByte))

	newCourse := course{}
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
	row := db.QueryRow(selectCourseId, newCourse.Title)

	var crse course
	err = row.Scan(&crse.Id)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	res := response{
		Message:  "Course created successfully",
		CourseId: crse.Id,
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

// updateCourses handles the PUT request to update a course.
// It expects a course ID as a URL parameter and the updated course details in the request body.
// Responds with a success status code (204 No Content) if the update is successful.
func updateCourses(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	err := auth.ValidateToken(token)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	courseId := r.URL.Query().Get("courseId")

	pasrseBody := r.Body
	bodyByte, err := io.ReadAll(pasrseBody)

	if err != nil {
		fmt.Println("parse error", err)
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println(string(bodyByte))

	newCourse := course{}
	err = json.Unmarshal([]byte(bodyByte), &newCourse)

	if err != nil {
		fmt.Println("Json Unmarshal :", err)
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = db.Exec(updateCourse, newCourse.Title, newCourse.Description,
		newCourse.Price, newCourse.ImageLink, newCourse.Published, courseId)

	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	res := response{
		Message: "Course updated successfully",
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

// getAllCourses will return all the courses
func getAllCourses(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	err := auth.ValidateToken(token)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var courses []course

	rows, err := db.Query(getCourses)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var c course
		err := rows.Scan(&c.Id, &c.Title, &c.Description,
			&c.Price, &c.ImageLink, &c.Published)

		courses = append(courses, c)
		if err != nil {
			fmt.Println("Json Unmarshal :", err)
			log.Fatal(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	jsonRes, err := json.Marshal(courses)

	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write(jsonRes)
}

// ******* User Routes ******


// userSignup creates an account for an admin.
func userSignup(w http.ResponseWriter, r *http.Request) {
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
	result, err := db.Exec(insertQuery, cred.UserName, cred.PassWord, user)

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
	ro,_ := db.Query(getUsrAdminId,cred.UserName,cred.PassWord)
	ro.Scan(&cred.Id)
	token := auth.GenerateJwt(cred.Id)

	res := response{
		Message: "User created successfully",
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

// userLogin handles the validation of login credentials.
// On successful validation, it returns a token.
// On failed validation, it returns a bad request.
func userLogin(w http.ResponseWriter, r *http.Request) {
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
	if cred.PassWord == dbCred.PassWord && dbCred.Role == "User" {
		fmt.Println("User logged in")
	}
	ro,_ := db.Query(getUsrAdminId,cred.UserName,cred.PassWord)
	ro.Scan(&cred.Id)
	token := auth.GenerateJwt(cred.Id)

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

