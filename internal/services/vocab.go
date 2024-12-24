package services

import (
	"context"
	"ielts-web-api/internal/models"
	"ielts-web-api/internal/repositories"

	"gorm.io/gorm"
)

func (s *Service) GetVocabCategoriresByUserId(ctx context.Context, userId string) ([]*models.UserVocabCategory, error) {
	params := models.QueryParams{
		Offset:    0,
		Limit:     4,
		QuerySort: models.QuerySort{},
		Selected:  []string{"id", "name"},
		Preload:   nil,
	}

	userClause := repositories.Clause(func(tx *gorm.DB) {
		tx.Where("user_id = ?", userId)
	})

	categories, err := s.vocabCategoriesRepo.List(ctx, params, userClause)
	if err != nil {
		return nil, err
	}

	return categories, nil
}
