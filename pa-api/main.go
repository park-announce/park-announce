package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"sync"
	"syscall"
	"time"

	"math"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"github.com/google/uuid"

	"github.com/gorilla/websocket"
)

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type SocketMessage struct {
	ClientId    string      `json:"client_id"`
	MessageType string      `json:"message_type"`
	Data        interface{} `json:"data"`
}

type SendSocketMessage struct {
	ClientId string `json:"client_id"`
	Message  string `json:"message"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var connections map[string]*websocket.Conn = make(map[string]*websocket.Conn)

func main() {

	httpServerQuit := make(chan os.Signal, 1)
	signal.Notify(httpServerQuit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, os.Interrupt)

	redisLockQuit := make(chan os.Signal, 1)
	signal.Notify(redisLockQuit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, os.Interrupt)

	redisHeartBeatQuit := make(chan os.Signal, 1)
	signal.Notify(redisHeartBeatQuit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, os.Interrupt)

	wgMain := &sync.WaitGroup{}
	wgMain.Add(3)

	redisLockObtained := make(chan bool)

	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	var instanceId int = 0
	var lockName string
	var lockValue string = uuid.NewString()

	go func(_wgMain *sync.WaitGroup) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in redis go rouitine", r)
			}
		}()
		defer func() {
			_wgMain.Done()
		}()

		for {

			lockName = fmt.Sprintf("global-api-instance-id-%d", instanceId)

			set, setErr := rdb.SetNX(lockName, lockValue, 0).Result()

			if setErr != nil {
				log.Fatal(setErr)
				instanceId++
				continue
			}
			if !set {
				log.Println("value could not be set")
				instanceId++
				continue
			}

			break
		}

		log.Printf("instanceId : %d, lockValue : %s", instanceId, lockValue)
		redisLockObtained <- true
		<-redisLockQuit

		log.Println("interrupt detect in redis go rouitine")
		lockValueFromRedis, getErr := rdb.Get(lockName).Result()

		log.Printf("lockName : %s, lockValue : %s, lockValueFromRedis : %s", lockName, lockValue, lockValueFromRedis)

		if getErr != nil {
			log.Fatal(getErr)
		}

		if lockValueFromRedis == lockValue {
			deleteResult, delErr := rdb.Del(lockName).Result()

			if delErr != nil {
				log.Fatal(delErr)
			}
			if deleteResult > 0 {
				log.Printf("redis key delete result. key : %s, value : %s, result : %d", lockName, lockValue, deleteResult)
			}

		}

	}(wgMain)

	<-redisLockObtained

	go func(_wgMain *sync.WaitGroup) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in redis heartbeat goroutine", r)
			}
		}()
		defer func() {
			_wgMain.Done()
		}()

		for {

			select {

			case <-redisHeartBeatQuit:

				return

			default:

				set, setErr := rdb.Expire(lockName, time.Second*30).Result()

				if setErr != nil {
					log.Fatal(setErr)
					continue
				}
				if !set {
					log.Println("expire could not be set")
					continue
				}

				time.Sleep(time.Second)

			}

		}

	}(wgMain)

	go func(_wgMain *sync.WaitGroup) {
		defer func() {
			_wgMain.Done()
		}()

		go func() {

			defer func() {
				if r := recover(); r != nil {
					fmt.Println("Recovered in http go routine", r)
				}
			}()
			// Initialize router
			router := mux.NewRouter()

			// Define API endpoints
			router.HandleFunc("/locations/nearby", getNearByLocations).Methods("GET")
			router.HandleFunc("/socket/connect", connectSocket)
			router.HandleFunc("/socket/messages", sendSampleMessage)

			log.Println("Server started on port 8000")
			log.Fatal(http.ListenAndServe(":8000", router))
		}()

		<-httpServerQuit
		log.Println("interrupt detect in http go rouitine")

	}(wgMain)

	wgMain.Wait()
}

func sendSampleMessage(w http.ResponseWriter, r *http.Request) {

	var sendSocketMessage SendSocketMessage
	json.NewDecoder(r.Body).Decode(&sendSocketMessage)

	conn := connections[sendSocketMessage.ClientId]

	// Write message back to browser
	if err := conn.WriteMessage(websocket.TextMessage, []byte(sendSocketMessage.Message)); err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("OK")
}

func connectSocket(w http.ResponseWriter, r *http.Request) {
	conn, connErr := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity

	if connErr != nil {
		log.Println(connErr)
	}

	for {
		// Read message from browser
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		// Print the message to the console
		fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

		var socketMessage SocketMessage
		deSerializationError := json.Unmarshal(msg, &socketMessage)
		if deSerializationError != nil {
			log.Println(deSerializationError)
			continue
		}

		//validate ClientId

		if socketMessage.ClientId == "" {
			log.Println("invalid client id")
			return
		}

		connections[socketMessage.ClientId] = conn

		// Write message back to browser
		if err = conn.WriteMessage(msgType, msg); err != nil {
			return
		}
	}
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
