package models

import "time"

// Owner represents a calendar owner
type Owner struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Timezone  string    `json:"timezone"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// CreateOwnerRequest represents the request body for creating an owner
type CreateOwnerRequest struct {
	Name     string `json:"name" binding:"required,min=1,max=100"`
	Email    string `json:"email" binding:"required,email,max=255"`
	Timezone string `json:"timezone" binding:"required,max=50"`
}

// UpdateOwnerRequest represents the request body for updating an owner
type UpdateOwnerRequest struct {
	Name     *string `json:"name,omitempty" binding:"omitempty,min=1,max=100"`
	Email    *string `json:"email,omitempty" binding:"omitempty,email,max=255"`
	Timezone *string `json:"timezone,omitempty" binding:"omitempty,max=50"`
}
