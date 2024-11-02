package services

import (
	"ielts-web-api/internal/repositories"
)

type Service struct {
	UserRepo   *repositories.UserRepository
	TargetRepo *repositories.TargetRepository
}

func NewService(repos ...interface{}) *Service {
	service := &Service{}
	for _, repo := range repos {
		switch repo.(type) {
		case *repositories.UserRepository:
			service.UserRepo = repo.(*repositories.UserRepository)
		case *repositories.TargetRepository:
			service.TargetRepo = repo.(*repositories.TargetRepository)
		default:
			panic("Unknown repository type provided")
		}
	}
	return service
}
