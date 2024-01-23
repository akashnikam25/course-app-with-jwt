package main

import (
	"course-app-with-jwt/auth"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
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

// validAuthuser will validate jwtToken
func validateAuthUser(h http.Handler) http.Handler {
	fmt.Println("Inside user ValidateAuthUser")
	var role string

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("Authorization")
		if strings.Contains(r.URL.Path, admin) {
			role = admin
		} else {
			role = user
		}
		userId, err := auth.ValidateToken(token, role)
		if err != nil {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		r.Header.Add("userId", strconv.Itoa(userId))
		h.ServeHTTP(w, r)
	})
}

func main() {

	defer db.Close()
	initDB()
	// Set up a channel to listen for interrupt signals
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)

	mainRouter := mux.NewRouter()

	mainRouter.HandleFunc("/admin/signup", adminSignup).Methods(http.MethodPost)
	mainRouter.HandleFunc("/admin/login", adminLogin).Methods(http.MethodPost)
	mainRouter.HandleFunc("/user/signup", userSignup).Methods(http.MethodPost)
	mainRouter.HandleFunc("/user/login", userLogin).Methods(http.MethodPost)

	adminRouter := mainRouter.PathPrefix("/admin").Subrouter()
	adminRouter.HandleFunc("/courses", createCourse).Methods(http.MethodPost)
	adminRouter.HandleFunc("/courses/{courseId}", updateCourses).Methods(http.MethodPut)
	adminRouter.HandleFunc("/courses", getAllCourses).Methods(http.MethodGet)
	adminRouter.Use(validateAuthUser)

	userRouter := mainRouter.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/courses", getAllCourses).Methods(http.MethodGet)
	userRouter.HandleFunc("/courses/{courseId}", purchaseCourse).Methods(http.MethodPost)
	userRouter.HandleFunc("/purchasedCourses", getAllPurchaseCourse).Methods(http.MethodGet)
	userRouter.Use(validateAuthUser)

	// Start the server
	fmt.Printf("Server is listening on port %s...\n", ":8000")
	go startServer(mainRouter)

	// Wait for an interrupt signal
	<-interruptChan

	// Handle cleanup and shutdown logic here...

	fmt.Println("Server is shutting down.")
	os.Exit(0)
}

func startServer(mainRouter *mux.Router) {
	if err := http.ListenAndServe(":8000", mainRouter); err != nil {
		fmt.Println("err :	", err)
	}
}
