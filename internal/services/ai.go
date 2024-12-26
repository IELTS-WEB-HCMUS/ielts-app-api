package services

import (
	"context"
	"fmt"
	"ielts-web-api/common"
	"ielts-web-api/internal/models"

	"gorm.io/gorm"
)

func (s *Service) LookUpVocab(ctx context.Context, req models.LookUpVocabRequest, userId string) (*models.Vocab, error) {
	vocabId := fmt.Sprintf("%d_%d", req.QuizId, req.VocabIndex)

	vocab, err := s.vocabRepo.GetDetailByConditions(ctx, func(tx *gorm.DB) {
		tx.Where("vocab_id = ?", vocabId)
	})
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	user.VocabUsageCount--
	_, err = s.userRepo.UpdateColumns(ctx, user.ID, map[string]interface{}{
		"vocab_usage_count": user.VocabUsageCount,
	})
	if err != nil {
		return nil, err
	}

	return vocab, nil
}

func (s *Service) CheckVocabUsageCount(ctx context.Context, userId string) error {
	user, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		return err
	}
	if user.VocabUsageCount <= 0 {
		return common.ErrVocabUsageCountExceeded
	}

	return nil
}
