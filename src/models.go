package main

// Product represents the product model as defined in OpenAPI spec
type Product struct {
	ProductID    int32  `json:"product_id"`
	SKU          string `json:"sku"`
	Manufacturer string `json:"manufacturer"`
	CategoryID   int32  `json:"category_id"`
	Weight       int32  `json:"weight"`
	SomeOtherID  int32  `json:"some_other_id"`
}

// ErrorResponse represents the error response model
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// ValidateProduct checks if all required fields are present and valid
func (p *Product) Validate() *ErrorResponse {
	if p.ProductID < 1 {
		return &ErrorResponse{
			Error:   "INVALID_INPUT",
			Message: "Invalid product data",
			Details: "product_id must be >= 1",
		}
	}
	if len(p.SKU) == 0 || len(p.SKU) > 100 {
		return &ErrorResponse{
			Error:   "INVALID_INPUT",
			Message: "Invalid product data",
			Details: "sku must be between 1 and 100 characters",
		}
	}
	if len(p.Manufacturer) == 0 || len(p.Manufacturer) > 200 {
		return &ErrorResponse{
			Error:   "INVALID_INPUT",
			Message: "Invalid product data",
			Details: "manufacturer must be between 1 and 200 characters",
		}
	}
	if p.CategoryID < 1 {
		return &ErrorResponse{
			Error:   "INVALID_INPUT",
			Message: "Invalid product data",
			Details: "category_id must be >= 1",
		}
	}
	if p.Weight < 0 {
		return &ErrorResponse{
			Error:   "INVALID_INPUT",
			Message: "Invalid product data",
			Details: "weight must be >= 0",
		}
	}
	if p.SomeOtherID < 1 {
		return &ErrorResponse{
			Error:   "INVALID_INPUT",
			Message: "Invalid product data",
			Details: "some_other_id must be >= 1",
		}
	}
	return nil
}
