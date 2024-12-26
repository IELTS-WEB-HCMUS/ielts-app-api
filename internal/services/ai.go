package services

import (
	"context"
	"fmt"
	"ielts-web-api/internal/models"

	"gorm.io/gorm"
)

func (s *Service) LookUpVocab(ctx context.Context, req models.LookUpVocabRequest) (*models.Vocab, error) {
	vocabId := fmt.Sprintf("%d_%d", req.QuizId, req.VocabIndex)

	vocab, err := s.vocabRepo.GetDetailByConditions(ctx, func(tx *gorm.DB) {
		tx.Where("vocab_id = ?", vocabId)
	})
	if err != nil {
		return nil, err
	}

	return vocab, nil
}
