package handlers

import (
	"ielts-web-api/common"

	"github.com/gin-gonic/gin"
)

func (h *Handler) BuyMoreAiVocabTurn(c *gin.Context) {
	ok, userJWTProfile := common.ProfileFromJwt(c)
	if !ok {
		c.JSON(common.INTERNAL_SERVER_ERR, gin.H{"error": "Authentication failed"})
		return
	}

	orderUrl, err := h.service.BuyMoreAiVocabTurn(c.Request.Context(), userJWTProfile.Id)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(common.SUCCESS_STATUS, gin.H{"message": "Create order successfully", "order_url": *orderUrl})
}
