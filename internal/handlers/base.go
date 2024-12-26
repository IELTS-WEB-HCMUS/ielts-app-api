package handlers

import (
	services "ielts-web-api/internal/services"
	"ielts-web-api/middleware"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *services.Service
}

func NewHandler(service *services.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// Define API route here
func (h *Handler) RegisterRoutes(c *gin.Engine) {
	userRoutes := c.Group("/api/users")
	{
		userRoutes.GET("", middleware.UserAuthentication, h.GetUserProfile)
		userRoutes.POST("/signup", h.SignUp)
		userRoutes.POST("/login", h.LogIn)
		userRoutes.GET("/target", middleware.UserAuthentication, h.GetTarget)
		userRoutes.POST("/target", middleware.UserAuthentication, h.CreateTarget)
		userRoutes.PATCH("/target", middleware.UserAuthentication, h.UpdateTarget)
	}
	authRoutes := c.Group("/api/auth")
	{
		authRoutes.POST("/gen-otp", h.GenerateOTP)
		authRoutes.POST("/validate-otp", h.ValidateOTP)
		authRoutes.POST("/reset-password", h.ResetPassword)
	}

	health := c.Group("api/health")
	{
		health.GET("/status", h.CheckStatusHealth)
	}

	quizzes := c.Group("/v1/quizzes")
	{
		//API Get Quiz Detail
		quizzes.GET("/:quiz_id", middleware.UserAuthentication, h.GetQuiz())
		//API Listing Quiz
		quizzes.GET("", middleware.OptionalUserAuthentication(), h.GetQuizzes())

		quizzes.POST("/:quiz_id/answer", middleware.UserAuthentication, h.SubmitQuiz())

	}

	tagSearches := c.Group("/v1/tag-searches")
	{
		//API Get Tag Search
		tagSearches.GET("", h.GetTagSearches())
	}

	answerRoutes := c.Group("/v1/answers")
	{
		answerRoutes.GET("/:answer_id", middleware.UserAuthentication, h.GetAnswer)
		answerRoutes.GET("/statistics", middleware.UserAuthentication, h.GetAnswerStatistic)
	}

	vocabRoutes := c.Group("/api/vocabs")
	{
		vocabRoutes.GET("/get-categories", middleware.UserAuthentication, h.GetVocabCategorires)
		vocabRoutes.POST("/update-category", middleware.UserAuthentication, h.UpdateVocabCategory)
		vocabRoutes.POST("/add", middleware.UserAuthentication, h.AddVocab)
		vocabRoutes.POST("/update", middleware.UserAuthentication, h.UpdateVocab)
		vocabRoutes.DELETE("/:id", middleware.UserAuthentication, h.DeleteVocab)
		vocabRoutes.POST("/", middleware.UserAuthentication, h.GetVocabs)
	}

	aiRoutes := c.Group("/api/ai")
	{
		aiRoutes.POST("/look-up-vocab", middleware.UserAuthentication, h.LookUpVocab)
	}
}
