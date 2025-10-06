package main

import (
	"sync"
)

// ProductStore manages in-memory storage of products
type ProductStore struct {
	mu       sync.RWMutex
	products map[int32]*Product
}

// NewProductStore creates a new product store
func NewProductStore() *ProductStore {
	return &ProductStore{
		products: make(map[int32]*Product),
	}
}

// Get retrieves a product by ID
func (s *ProductStore) Get(productID int32) (*Product, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	product, exists := s.products[productID]
	return product, exists
}

// Set adds or updates a product
func (s *ProductStore) Set(product *Product) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.products[product.ProductID] = product
}
