package services

import (
	"context"
	"ielts-web-api/internal/models"
)

func (s *Service) GetUserProfileById(ctx context.Context, id string) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
