package handlers

import (
	"ielts-web-api/common"
	"ielts-web-api/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) SignUp(c *gin.Context) {
	var req models.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	err := h.service.SignupUser(c, req)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusCreated, common.ResponseOk("User created successfully"))
}

func (h *Handler) LogIn(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	if req.Email == nil && req.IdToken == nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	token, err := h.service.LoginUser(c, req)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(common.SUCCESS_STATUS, common.BaseResponseMess(common.SUCCESS_STATUS, "Login successfully", token))
}

func (h *Handler) RequestResetPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		common.AbortWithError(c, err)
		return
	}

	otp, err := h.service.GenerateOTP(c.Request.Context(), req.Email)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	err = common.SendOTPEmail(common.FromEmail, req.Email, otp)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(common.SUCCESS_STATUS, common.BaseResponse(http.StatusOK, "OTP email sent successfully", "", ""))
}

func (h *Handler) ValidateOTP(c *gin.Context) {
	var req models.OTPValidateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		common.AbortWithError(c, err)
		return
	}

	err := h.service.ValidateOTP(c.Request.Context(), req.Email, req.OTP)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(common.SUCCESS_STATUS, common.BaseResponse(http.StatusOK, "OTP validated successfully", "", ""))
}

func (h *Handler) ResetPassword(c *gin.Context) {
	var req models.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.AbortWithError(c, err)
		return
	}
	err := h.service.ResetPassword(c.Request.Context(), req.Email, req.NewPassword)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(common.SUCCESS_STATUS, common.BaseResponse(http.StatusOK, "Password Reset successfully", "", ""))
}
