package handlers

import (
	"ielts-web-api/common"
	"ielts-web-api/internal/models"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetAnswer(c *gin.Context) {
	ok, userProfile := common.ProfileFromJwt(c)
	if !ok {
		common.AbortWithError(c, common.ErrUserNotFound)
		return
	}

	// validate url params
	var answerParamsUri *models.AnswerParamsUri
	if err := c.ShouldBindUri(&answerParamsUri); err != nil {
		common.AbortWithError(c, err)
		return
	}

	result, err := h.service.GetAnswer(c, userProfile.Id, answerParamsUri.AnswerID)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(common.SUCCESS_STATUS, common.ResponseOk(result))
}

func (h *Handler) GetAnswerStatistic(c *gin.Context) {
	ok, userProfile := common.ProfileFromJwt(c)
	if !ok {
		common.AbortWithError(c, common.ErrNotAuthorized)
		return
	}

	var paramsUri models.AnswerStatisticsQuery
	if err := c.ShouldBindQuery(&paramsUri); err != nil {
		c.AbortWithStatusJSON(common.BAD_REQUEST_STATUS, common.BaseResponse(common.REQUEST_FAILED, "Thông tin không hợp lệ", err.Error(), nil))
		return
	}

	// validate skill and type
	if paramsUri.Type == nil {
		common.AbortWithError(c, common.ErrAnswerStatisticTypeRequired)
		return
	}

	data, err := h.service.GetAnswerStatistic(c, userProfile.Id, &paramsUri)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(common.SUCCESS_STATUS, common.ResponseOk(data))
}
