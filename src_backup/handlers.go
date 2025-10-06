package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Handlers struct {
	store *ProductStore
}

func NewHandlers(store *ProductStore) *Handlers {
	return &Handlers{store: store}
}

// GetProduct handles GET /products/{productId}
func (h *Handlers) GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productIDStr := vars["productId"]

	// Parse and validate productId
	productID, err := strconv.ParseInt(productIDStr, 10, 32)
	if err != nil || productID < 1 {
		respondWithError(w, http.StatusBadRequest, &ErrorResponse{
			Error:   "INVALID_INPUT",
			Message: "Invalid product ID",
			Details: "Product ID must be a positive integer",
		})
		return
	}

	// Get product from store
	product, exists := h.store.Get(int32(productID))
	if !exists {
		respondWithError(w, http.StatusNotFound, &ErrorResponse{
			Error:   "NOT_FOUND",
			Message: "Product not found",
			Details: "No product exists with the given ID",
		})
		return
	}

	// Return product
	respondWithJSON(w, http.StatusOK, product)
}

// AddProductDetails handles POST /products/{productId}/details
func (h *Handlers) AddProductDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productIDStr := vars["productId"]

	// Parse and validate productId from path
	productID, err := strconv.ParseInt(productIDStr, 10, 32)
	if err != nil || productID < 1 {
		respondWithError(w, http.StatusBadRequest, &ErrorResponse{
			Error:   "INVALID_INPUT",
			Message: "Invalid product ID",
			Details: "Product ID must be a positive integer",
		})
		return
	}

	// Parse request body
	var product Product
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&product); err != nil {
		respondWithError(w, http.StatusBadRequest, &ErrorResponse{
			Error:   "INVALID_INPUT",
			Message: "Invalid JSON format",
			Details: err.Error(),
		})
		return
	}

	// Validate product data
	if validationErr := product.Validate(); validationErr != nil {
		respondWithError(w, http.StatusBadRequest, validationErr)
		return
	}

	// Ensure product_id in body matches path parameter
	if product.ProductID != int32(productID) {
		respondWithError(w, http.StatusBadRequest, &ErrorResponse{
			Error:   "INVALID_INPUT",
			Message: "Product ID mismatch",
			Details: "Product ID in path must match product_id in request body",
		})
		return
	}

	// Save product to store
	h.store.Set(&product)

	// Return 204 No Content
	w.WriteHeader(http.StatusNoContent)
}

// Helper functions
func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, statusCode int, err *ErrorResponse) {
	respondWithJSON(w, statusCode, err)
}
