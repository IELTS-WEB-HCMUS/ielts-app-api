package handlers

import (
	"ielts-web-api/common"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetTagSearches() gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := h.service.GetTagSearches(c)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}
		c.JSON(common.SUCCESS_STATUS, common.ResponseOk(data))
	}
}
