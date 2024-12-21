package handlers

import (
	"ielts-web-api/common"
	"ielts-web-api/internal/models"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetQuizzes() gin.HandlerFunc {
	return func(c *gin.Context) {
		var userID string
		ok, userJWTProfile := common.ProfileFromJwt(c)
		if ok {
			userID = userJWTProfile.Id
		}

		var params = models.ListQuizzesParamsUri{}
		if err := c.ShouldBindQuery(&params); err != nil {
			c.AbortWithStatusJSON(common.BAD_REQUEST_STATUS, common.BaseResponse(common.REQUEST_FAILED, "Thông tin không hợp lệ", err.Error(), nil))
			return
		}

		data, err := h.service.GetQuizzes(c, userID, &params)
		if err != nil {
			c.AbortWithStatusJSON(common.INTERNAL_SERVER_ERR, common.BaseResponse(common.REQUEST_FAILED, "Đã xảy ra lỗi!", err.Error(), nil))
			return
		}

		c.JSON(common.SUCCESS_STATUS, common.ResponseOk(data))
	}
}

func (h *Handler) GetQuiz() gin.HandlerFunc {
	return func(c *gin.Context) {
		var quizParamsUri *models.QuizParamsUri
		if err := c.ShouldBindUri(&quizParamsUri); err != nil {
			common.AbortWithError(c, err)
			return
		}
		if err := c.ShouldBindQuery(&quizParamsUri); err != nil {
			common.AbortWithError(c, err)
			return
		}
		userID := ""
		ok, userJWTProfile := common.ProfileFromJwt(c)
		if ok {
			userID = userJWTProfile.Id
		}

		data, err := h.service.GetQuiz(c, quizParamsUri, userID)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}

		c.JSON(common.SUCCESS_STATUS, common.ResponseOk(data))
	}
}

func (h *Handler) SubmitQuiz() gin.HandlerFunc {
	return func(c *gin.Context) {
		var quizParamsUri *models.QuizParamsUri
		if err := c.ShouldBindUri(&quizParamsUri); err != nil {
			common.AbortWithError(c, err)
			return
		}

		var body models.QuizAnswer
		if err := c.ShouldBind(&body); err != nil {
			common.AbortWithError(c, err)
			return
		}

		ok, user := common.ProfileFromJwt(c)
		if !ok {
			common.AbortWithError(c, common.ErrNotAuthorized)
			return
		}

		answer, err := h.service.SubmitQuizAnswer(c, user.Id, quizParamsUri.QuizID, body)
		if err != nil {
			common.AbortWithError(c, err)
			return
		}
		c.JSON(common.SUCCESS_STATUS, common.ResponseOk(answer))
	}
}
