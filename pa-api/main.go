package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// User struct represents the user model
type User struct {
	ID       int    `json:"id"`
	Phone    string `json:"phone"`
	Nickname string `json:"nickname"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Email    string `json:"email"`
}

// Geolocation struct represents the geolocation model
type Geolocation struct {
	ID        int     `json:"id"`
	UserID    int     `json:"user_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

var db *sql.DB

func main() {
	// Establish database connection
	connStr := "postgres://park-announce:PosgresDb1591*@db/padb?sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize router
	router := mux.NewRouter()

	// Define API endpoints
	router.HandleFunc("/users", createUser).Methods("POST")
	router.HandleFunc("/users/{id}", getUser).Methods("GET")
	router.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	// Start the server
	log.Println("Server started on port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}

// createUser creates a new user
func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	json.NewDecoder(r.Body).Decode(&user)

	// Insert the user into the database
	err := db.QueryRow("INSERT INTO users (phone, nickname, name, surname, email) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		user.Phone, user.Nickname, user.Name, user.Surname, user.Email).Scan(&user.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// getUser retrieves a user by ID
func getUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var user User
	err := db.QueryRow("SELECT id, phone, nickname, name, surname, email FROM users WHERE id = $1", id).Scan(
		&user.ID, &user.Phone, &user.Nickname, &user.Name, &user.Surname, &user.Email)
	if err != nil {
		log.Println(err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// updateUser updates a user by ID
func updateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var user User
	json.NewDecoder(r.Body).Decode(&user)

	// Update the user in the database
	_, err := db.Exec("UPDATE users SET phone = $1, nickname = $2, name = $3, surname = $4, email = $5 WHERE id = $6",
		user.Phone, user.Nickname, user.Name, user.Surname, user.Email, id)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// deleteUser deletes a user by ID
func deleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	// Delete the user from the database
	_, err := db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
