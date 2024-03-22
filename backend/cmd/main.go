package main

import (
	"dsproject/backend/handlers"
	"log"
	"net/http"
	"os"
)

// Default values
const DEFAULT_LOGGING_PATH = "./logs"
const DEFAULT_PORT = "8080"

// Control whether stdout console output should be suppressed (only works if logging is deactivated)
const SUPPRESS_CONSOLE_OUTPUT = false

// API Paths (embed trailing slashes to retain all URL control here)

const VERSION = "v1/"
const API_PATH = "api/" + VERSION
const START_PATH = API_PATH + "start/"

func main() {

	// Check for custom port
	port := os.Getenv("port")
	if port == "" {
		port = DEFAULT_PORT
	}

	log.Println("Server running on port " + port)
	http.HandleFunc("/start", handlers.SortRequest)
	http.HandleFunc("/ws", handlers.ServeWs)

	// HTTP
	err := http.ListenAndServe("localhost:"+port, nil)
	if err != nil {
		log.Fatal("Failed server launch:", err)
	}
}
