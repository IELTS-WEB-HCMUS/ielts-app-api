package handlers

import (
	services "ielts-web-api/internal/services"
	"ielts-web-api/middleware"

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

// Define API route here
func (h *Handler) RegisterRoutes(c *gin.Engine) {
	userRoutes := c.Group("/api/users")
	{
		userRoutes.GET("", middleware.UserAuthentication, h.GetUserProfile)
		// userRoutes.POST("/signup", h.SignUp)
		// userRoutes.POST("/login", h.LogIn)
		userRoutes.GET("/target", middleware.UserAuthentication, h.GetTarget)
	}
}
