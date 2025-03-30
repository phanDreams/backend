package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func main() {

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8000"
	}

	http.HandleFunc("/api/v1", helloWorldHandler)

	log.Printf("Server is running on port %s...\n", serverPort)
	err := http.ListenAndServe(":"+serverPort, nil)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
