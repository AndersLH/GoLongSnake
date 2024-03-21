package handlers

import (
	customhttp "dsproject/backend/customHTTP"
	"fmt"
	"net/http"
	"strings"
)

type player struct {
	Username string `json:"username"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
}

type board struct {
	Size    int      `json:"size"`
	Players []player `json:players`
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
		getBoard(w, r)
	case http.MethodPost:
		moveSnake(w, r)
	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}

}

func createUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("cu")
}

func getBoard(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get")
}

func inititalizeBoard(w http.ResponseWriter, r *http.Request) {
	var board board
	customhttp.Decode(r, &board)

}

func moveSnake(w http.ResponseWriter, r *http.Request) {
	var snakeUser snake
	customhttp.Decode(r, &snakeUser)

	//ensure lowercase
	snakeUser.Direction = strings.ToLower(snakeUser.Direction)
	//TODO add logic for less than 0 or more than maximum. teleport on border?
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

	fmt.Println(snakeUser)
	customhttp.Encode(w, &snakeUser)
	// if snakeUser.Name == "" {
	// 	http.Error(w, "Dataset name can not be blank", http.StatusNotAcceptable)
	// 	return
	// }
}

func cors(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS") // Allow all methods
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")                    // Allow content-type header
	w.Header().Set("Access-Control-Allow-Credentials", "true")                        // Allow credential header
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081")            // Allow port 8080 to make requests (where frontend runs in dev mode)

	return r.Method == http.MethodOptions // If method is Options, return true to allow logic in handler return early - prevents double work.
}
