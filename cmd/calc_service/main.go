package main

import (
	"fmt"
	"log"
	"net/http"

	"calc-service/internal/handler"
)

func main() {
	// Setup routes
	http.HandleFunc("/api/v1/calculate", handler.HandleCalculate)

	// Start server
	port := 8080
	fmt.Printf("Server starting on port %d...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
