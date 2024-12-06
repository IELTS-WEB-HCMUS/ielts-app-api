package services

import (
	"ielts-web-api/internal/repositories"
)

type Service struct {
	userRepo       *repositories.UserRepository
	targetRepo     *repositories.TargetRepository
	otpRepo        *repositories.OTPRepository
	otpAttemptRepo *repositories.OTPAttemptRepository
}

func NewService(repos ...interface{}) *Service {
	service := &Service{}
	for _, repo := range repos {
		switch repo.(type) {
		case *repositories.UserRepository:
			service.userRepo = repo.(*repositories.UserRepository)
		case *repositories.TargetRepository:
			service.targetRepo = repo.(*repositories.TargetRepository)
		case *repositories.OTPRepository:
			service.otpRepo = repo.(*repositories.OTPRepository)
		case *repositories.OTPAttemptRepository:
			service.otpAttemptRepo = repo.(*repositories.OTPAttemptRepository)
		default:
			panic("Unknown repository type provided")
		}
	}
	return service
}
