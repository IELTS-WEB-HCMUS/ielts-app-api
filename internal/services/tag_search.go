package services

import (
	"context"
	"ielts-web-api/internal/models"

	"gorm.io/gorm"
)

func (s *Service) GetTagSearches(ctx context.Context) ([]*models.TagSearchPosition, error) {
	var positions []*models.TagSearchPosition
	var err error

	positions, err = s.tagSearchPositionRepo.List(ctx, models.QueryParams{}, func(tx *gorm.DB) {
		tx.Preload("Tags").Order("id ASC").Find(&positions)
	})
	if err != nil {
		return nil, err
	}
	return positions, nil
}
