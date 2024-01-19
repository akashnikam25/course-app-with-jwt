package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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

	defer db.Close()
	initDB()
	// Set up a channel to listen for interrupt signals
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)

	router := mux.NewRouter()
	router.HandleFunc("/admin/signup", adminSignup).Methods(http.MethodPost)
	router.HandleFunc("/admin/login", adminLogin).Methods(http.MethodPost)
	router.HandleFunc("/admin/course", createCourse).Methods(http.MethodPost)
	router.HandleFunc("/admin/courses/{courseId}", updateCourses).Methods(http.MethodPut)
	router.HandleFunc("/admin/courses", getAllCourses).Methods(http.MethodGet)
	router.HandleFunc("/user/signup", userSignup).Methods(http.MethodPost)
	router.HandleFunc("/user/login", userLogin).Methods(http.MethodPost)
	router.HandleFunc("/user/courses", getAllCourses).Methods(http.MethodGet)
	router.HandleFunc("/user/courses/{courseId}", purchaseCourse).Methods(http.MethodPost)
	// Start the server

	fmt.Printf("Server is listening on port %s...\n", ":8000")
	go http.ListenAndServe(":8000", router)

	// Wait for an interrupt signal
	<-interruptChan

	// Handle cleanup and shutdown logic here...

	fmt.Println("Server is shutting down.")
	os.Exit(0)

}
