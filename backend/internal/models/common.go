package models

// Pagination represents pagination metadata
type Pagination struct {
	Page       int  `json:"page"`
	PageSize   int  `json:"pageSize"`
	TotalItems int  `json:"totalItems"`
	TotalPages int  `json:"totalPages"`
	HasNext    bool `json:"hasNext"`
	HasPrev    bool `json:"hasPrev"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse[T any] struct {
	Items      []T        `json:"items"`
	Pagination Pagination `json:"pagination"`
}

// ErrorResponse represents an error API response
type ErrorResponse struct {
	Error       string       `json:"error"`
	Message     string       `json:"message"`
	Details     string       `json:"details,omitempty"`
	FieldErrors []FieldError `json:"fieldErrors,omitempty"`
}

// FieldError represents a single field validation error
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
