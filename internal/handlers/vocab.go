package handlers

import (
	"ielts-web-api/common"
	"ielts-web-api/internal/models"

	"net/http"

	"github.com/gin-gonic/gin"

	"strconv"
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

func (h *Handler) AddVocab(c *gin.Context) {
	var req models.UserVocabBankAddRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	message, err := h.service.AddVocab(c, req)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.ResponseOk(message))
}

func (h *Handler) UpdateVocab(c *gin.Context) {
	var req models.UserVocabBankUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	message, err := h.service.UpdateVocab(c, req)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.ResponseOk(message))
}

func (h *Handler) DeleteVocab(c *gin.Context) {
	// Get ID from the URL parameter
	idParam := c.Param("id")
	if idParam == "" {
		common.AbortWithError(c, common.ErrIdRequired)
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		common.AbortWithError(c, common.ErrIdMustBeInt)
		return
	}

	// Call the service layer to delete the vocab
	message, err := h.service.DeleteVocab(c.Request.Context(), id)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, common.ResponseOk(message))
}

func (h *Handler) GetVocabs(c *gin.Context) {
	var req models.UserVocabBankGetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.AbortWithError(c, common.ErrInvalidInput)
		return
	}

	data, err := h.service.GetVocabs(c, req)
	if err != nil {
		c.JSON(common.INTERNAL_SERVER_ERR, gin.H{"error": err.Error()})
		return
	}
	c.JSON(common.SUCCESS_STATUS, gin.H{"message": "Get user vocabularies successfully", "data": data})
}
