package handlers

import (
	services "hotel-booking-system/internal/services"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *services.Service
}

func NewHandler(service *services.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// Declare All the API Routes at here!!!
func (h *Handler) RegisterRoutes(c *gin.Engine) {

}
