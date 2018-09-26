package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/gorilla/mux"
	"github.com/mohammedajao/rest-api/app"
	"github.com/mohammedajao/rest-api/app/models"
)

type UserController struct {
	DB *sql.DB
}

type message struct {
	Data string
}

func (u *UserController) Init(a *app.App) {
	// Define user routes for the API through our App structure
	c := a.Router.PathPrefix("/user").Subrouter()
	c.HandleFunc("/{id}", u.getUser).Methods("GET")
	c.HandleFunc("/create", u.createUser).Methods("POST")
	c.HandleFunc("/login", u.login).Methods("POST")
}

func (u *UserController) hashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, 10)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func (u *UserController) comparePassword(hashedPassword string, pwd []byte) bool {
	byteHash := []byte(hashedPassword)
	err := bcrypt.CompareHashAndPassword(byteHash, pwd)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func (u *UserController) getUser(w http.ResponseWriter, r *http.Request) {
	// Break down r.URL path to get URL variables
	vars := mux.Vars(r)

	// Grab new empty instance/values of user model to reference and fill in
	var user models.User
	w.WriteHeader(http.StatusOK)

	// Ensure the id passed is an int for the mySQL Query
	id, _ := strconv.ParseInt(vars["id"], 10, 64)

	// Commence search
	rows, err := u.DB.Query("SELECT id, username, email, createdAt FROM users WHERE id = ?", id)

	// Ensure it is closed once the function is over
	// Prevents wasted resources and too many open connections
	defer rows.Close()
	if err != nil {
		log.Println(err)
	}

	// Loop through rows of the unique IDs to find the user
	for rows.Next() {
		// Fill in user instance/values for comparison
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
		if err != nil {
			log.Println(err)
		}

		// Convert to int64 due to strconv library
		if int64(user.ID) == id {
			// Send JSON back to client for the GET request
			payload, _ := json.Marshal(user)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(payload))
		}
	}
}

func (u *UserController) login(w http.ResponseWriter, r *http.Request) {
	// Decode json password to Golang's requirements
	var user models.User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	plainPwd := []byte(user.Password)

	// Error checking
	if err != nil {
		log.Println(err)
	}

	// Query database to search for a unique email and a user with the name
	rows, err := u.DB.Query("SELECT * FROM users WHERE email = ?", user.Email)
	defer rows.Close()
	for rows.Next() { // Loop through rows to scan
		// Error checking
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
		if err != nil {
			log.Println(err)
		}

		// Compare hashedPassword with submitted password for authentication
		if u.comparePassword(user.Password, plainPwd) {
			// Write back to the client that the login was successful
			msg := message{"success"}
			payload, _ := json.Marshal(msg)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(payload))
		}
	}
}

func (u *UserController) createUser(w http.ResponseWriter, r *http.Request) {
	// Decode the POST request from its JSON
	var user models.User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		log.Println(err)
	}
	w.WriteHeader(http.StatusOK)

	// Prepare statement to prevent injection
	stmt, err := u.DB.Prepare("INSERT INTO users (username, password, email, createdAt) VALUES(?, ?, ?, ?)")
	hashedPassword := u.hashAndSalt([]byte(user.Password))

	// Store as UNIX to calculate time for 1970 and easy date conversion
	_, err = stmt.Exec(user.Name, hashedPassword, user.Email, time.Now().Unix())
	if err != nil {
		log.Println(err)
	}
}
