package services

import (
	"context"
	"fmt"
	"ielts-web-api/internal/models"
)

func (s *Service) GetTargetById(ctx context.Context, id string) (*models.Target, error) {
	fmt.Println(id)
	target, err := s.TargetRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return target, nil
}
