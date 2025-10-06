package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	// Create product store
	store := NewProductStore()

	// Create handlers
	handlers := NewHandlers(store)

	// Create router
	router := mux.NewRouter()

	// Register routes
	router.HandleFunc("/products/{productId}", handlers.GetProduct).Methods("GET")
	router.HandleFunc("/products/{productId}/details", handlers.AddProductDetails).Methods("POST")

	// Health check endpoint (optional but useful)
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(addr, router))
}
