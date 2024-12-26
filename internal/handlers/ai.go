package handlers

import (
	"ielts-web-api/common"
	"ielts-web-api/internal/models"

	"github.com/gin-gonic/gin"
)

func (h *Handler) LookUpVocab(c *gin.Context) {
	var req models.LookUpVocabRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	ok, userJWTProfile := common.ProfileFromJwt(c)
	if !ok {
		c.JSON(common.INTERNAL_SERVER_ERR, gin.H{"error": "Authentication failed"})
		return
	}
	err := h.service.CheckVocabUsageCount(c, userJWTProfile.Id)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	data, err := h.service.LookUpVocab(c, req, userJWTProfile.Id)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(common.SUCCESS_STATUS, gin.H{"message": "Look up vocabulary successfully", "data": data})
}
