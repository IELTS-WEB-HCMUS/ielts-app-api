package handlers

import (
	"ielts-web-api/common"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetUserProfile(c *gin.Context) {
	ok, userJWTProfile := common.ProfileFromJwt(c)
	if !ok {
		c.JSON(common.INTERNAL_SERVER_ERR, gin.H{"error": "User not found"})
		return
	}
	data, err := h.service.GetUserProfileById(c, userJWTProfile.Id)
	if err != nil {
		c.JSON(common.INTERNAL_SERVER_ERR, gin.H{"error": err.Error()})
		return
	}
	c.JSON(common.SUCCESS_STATUS, gin.H{"message": "Get user succerfully", "data": data})
}
