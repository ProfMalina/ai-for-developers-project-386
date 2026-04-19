package handlers

import (
	"net/http"

	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/models"
	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/repositories"
	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/services"
	"github.com/gin-gonic/gin"
)

// OwnerHandler handles HTTP requests for owners
type OwnerHandler struct {
	service *services.OwnerService
}

// NewOwnerHandler creates a new owner handler
func NewOwnerHandler() *OwnerHandler {
	return NewOwnerHandlerWithService(services.NewOwnerService(repositories.NewOwnerRepository()))
}

func NewOwnerHandlerWithService(service *services.OwnerService) *OwnerHandler {
	return &OwnerHandler{
		service: service,
	}
}

// Create handles POST /api/owners
func (h *OwnerHandler) Create(c *gin.Context) {
	var req models.CreateOwnerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "Invalid request body: "+err.Error())
		return
	}

	owner, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		BadRequest(c, "Failed to create owner: "+err.Error())
		return
	}

	SuccessResponse(c, http.StatusCreated, owner)
}

// GetByID handles GET /api/owners/{id}
func (h *OwnerHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	owner, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		NotFound(c, "Owner")
		return
	}

	SuccessResponse(c, http.StatusOK, owner)
}
