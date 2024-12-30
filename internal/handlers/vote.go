package handlers

import (
	"ielts-web-api/common"
	"ielts-web-api/internal/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) Vote(c *gin.Context){
	var req models.VoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	message, err := h.service.Vote(c, req)
	if err != nil {
		common.AbortWithError(c,err)
		return
	}

	c.JSON(common.SUCCESS_STATUS, common.ResponseOk(message))
}