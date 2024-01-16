package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	
)

var (
	db  *sql.DB
	err error
)

func initDB() {
	connectionString := "root:akash@tcp(127.0.0.1:3306)/courseapp"

	// Open a connection to the database
	db, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MySQL database!")

}

func main() {
	fmt.Println("Hello world")
	defer db.Close()
	initDB()
	router := mux.NewRouter()
	router.HandleFunc("/admin/signup", adminSignup).Methods("POST")
	router.HandleFunc("/admin/login", adminLogin).Methods("POST")
	router.HandleFunc("/admin/course", createCourse).Methods("POST")
	// Start the server

	fmt.Printf("Server is listening on port %s...\n", ":8000")
	http.ListenAndServe(":8000", router)

}
