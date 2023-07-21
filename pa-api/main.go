package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"

	"math"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
)

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

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

	// Create a pool with go-redis (or redigo) which is the pool redisync will
	// use while communicating with Redis. This can also be any pool that
	// implements the `redis.Pool` interface.
	client := goredislib.NewClient(&goredislib.Options{
		Addr: "redis:6379",
	})
	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)

	// Create an instance of redisync to be used to obtain a mutual exclusion
	// lock.
	rs := redsync.New(pool)

	// Obtain a new mutex by using the same name for all instances wanting the
	// same lock.

	instanceId := 0

	for {
		mutexname := fmt.Sprintf("global-api-instance-id-%d", instanceId)
		mutex := rs.NewMutex(mutexname)

		// Obtain a lock for our given mutex. After this is successful, no one else
		// can obtain the same lock (the same mutex name) until we unlock it.
		if err := mutex.Lock(); err != nil {
			instanceId++
			continue
		}

		break
	}

	log.Printf("instanceId : %d", instanceId)

	// Do your work that requires the lock.

	// Release the lock so other processes or threads can obtain a lock.
	/*
		if ok, err := mutex.Unlock(); !ok || err != nil {
			panic("unlock failed")
		}
	*/

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

	router.HandleFunc("/locations/nearby", getNearByLocations).Methods("GET")

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

// getNearByLocations retrieves nearest locations
func getNearByLocations(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	latitude, _ := strconv.ParseFloat(params["latitude"], 64)
	longitude, _ := strconv.ParseFloat(params["longitude"], 64)
	//distance, _ := strconv.ParseFloat(params["distance"], 64)

	locations := []Location{
		{Latitude: 37.7749, Longitude: -122.4194}, // San Francisco, CA
		{Latitude: 40.7128, Longitude: -74.0060},  // New York, NY
		{Latitude: 34.0522, Longitude: -118.2437}, // Los Angeles, CA
		// Add more locations as needed
	}

	n := 2
	nearestLocations := findNearestLocations(locations, latitude, longitude, n)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nearestLocations)
}

const earthRadius = 6371 // Earth's radius in kilometers

func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	// Convert latitude and longitude from degrees to radians
	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	// Calculate differences in latitude and longitude
	dLat := lat2Rad - lat1Rad
	dLon := lon2Rad - lon1Rad

	// Calculate the central angle using the Haversine formula
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	// Calculate the distance using the Earth's radius
	distance := earthRadius * c
	return distance
}

func findNearestLocations(locations []Location, refLat, refLon float64, n int) []Location {
	// Create a slice to store distances and their corresponding indices
	type DistanceIndex struct {
		distance float64
		index    int
	}
	distances := make([]DistanceIndex, len(locations))

	// Calculate distances and store them with their indices
	for i, loc := range locations {
		distances[i] = DistanceIndex{
			distance: haversine(refLat, refLon, loc.Latitude, loc.Longitude),
			index:    i,
		}
	}

	// Sort the distances in ascending order
	sort.Slice(distances, func(i, j int) bool {
		return distances[i].distance < distances[j].distance
	})

	// Create a slice to store the n nearest locations
	nearestLocations := make([]Location, n)

	// Extract the n nearest locations from the sorted distances
	for i := 0; i < n; i++ {
		nearestLocations[i] = locations[distances[i].index]
	}

	return nearestLocations
}
