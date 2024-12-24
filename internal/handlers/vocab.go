package handlers

import (
	"ielts-web-api/common"
	"ielts-web-api/internal/models"

	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetVocabCategorires(c *gin.Context) {
	ok, userJWTProfile := common.ProfileFromJwt(c)
	if !ok {
		c.JSON(common.INTERNAL_SERVER_ERR, gin.H{"error": "Authentication failed"})
		return
	}
	data, err := h.service.GetVocabCategoriresByUserId(c, userJWTProfile.Id)
	if err != nil {
		c.JSON(common.INTERNAL_SERVER_ERR, gin.H{"error": err.Error()})
		return
	}
	c.JSON(common.SUCCESS_STATUS, gin.H{"message": "Get user vocabulary categories succerfully", "data": data})
}

func (h *Handler) UpdateVocabCategory(c *gin.Context) {
	var req models.UpdateVocabCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	message, err := h.service.UpdateVocabCategory(c, req)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.ResponseOk(message))
}
