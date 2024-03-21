package handlers

import (
	customhttp "dsproject/backend/customHTTP"
	"fmt"
	"net/http"
	"strings"
)

const GRIDX int = 16
const GRIDY int = 10

type grid struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type player struct {
	Username string `json:"username"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
}

type snake struct {
	Name      string `json:"name"`
	Direction string `json:"direction"`
	X         int    `json:"x"`
	Y         int    `json:"y"`
}

func SortRequest(w http.ResponseWriter, r *http.Request) {
	if cors(w, r) {
		return
	}

	switch r.Method {
	case http.MethodGet:
		getGridSize(w, r)
	case http.MethodPost:
		moveSnake(w, r)
	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}

}

func getGridSize(w http.ResponseWriter, r *http.Request) {
	var grid grid
	grid.X = GRIDX
	grid.Y = GRIDY
	customhttp.Encode(w, &grid)
}

func moveSnake(w http.ResponseWriter, r *http.Request) {
	var snakeUser snake
	customhttp.Decode(r, &snakeUser)

	//Ensure lowercase
	snakeUser.Direction = strings.ToLower(snakeUser.Direction)

	switch snakeUser.Direction {
	case "arrowup":
		snakeUser.Y -= 1
	case "arrowdown":
		snakeUser.Y += 1
	case "arrowright":
		snakeUser.X += 1
	case "arrowleft":
		snakeUser.X -= 1
	default:
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	switch {
	case snakeUser.X >= GRIDX:
		snakeUser.X = 0
	case snakeUser.X < 0:
		snakeUser.X = GRIDX - 1
	}

	switch {
	case snakeUser.Y >= GRIDY:
		snakeUser.Y = 0
	case snakeUser.Y < 0:
		snakeUser.Y = GRIDY - 1
	}

	fmt.Println(snakeUser)
	customhttp.Encode(w, &snakeUser)
}

func cors(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS") // Allow all methods
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")                    // Allow content-type header
	w.Header().Set("Access-Control-Allow-Credentials", "true")                        // Allow credential header
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081")            // Allow port 8080 to make requests (where frontend runs in dev mode)

	return r.Method == http.MethodOptions // If method is Options, return true to allow logic in handler return early - prevents double work.
}
