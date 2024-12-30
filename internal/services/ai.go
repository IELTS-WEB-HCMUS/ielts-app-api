package services

import (
	"context"
	"fmt"
	"ielts-web-api/common"
	"ielts-web-api/internal/models"

	"errors"

	"gorm.io/gorm"
)

func (s *Service) LookUpVocabLinear(ctx context.Context, quizId int, sentenceIndex int, word string) (*models.Vocab, error) {
	vocabIdPattern := fmt.Sprintf("%d_%d_%%", quizId, sentenceIndex)

	vocab, err := s.vocabRepo.GetDetailByConditions(ctx, func(tx *gorm.DB) {
		tx.Where("vocab_id LIKE ? AND value = ?", vocabIdPattern, word)
	})
	if err != nil {
		return nil, err
	}

	return vocab, nil
}

func (s *Service) LookUpVocab(ctx context.Context, req models.LookUpVocabRequest, userId string) (*models.Vocab, error) {
	vocabId := fmt.Sprintf("%d_%d_%d", req.QuizId, req.SentenceIndex, req.WordIndex)

	vocab, err := s.vocabRepo.GetDetailByConditions(ctx, func(tx *gorm.DB) {
		tx.Where("vocab_id = ?", vocabId)
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			vocab, err = s.LookUpVocabLinear(ctx, req.QuizId, req.SentenceIndex, req.Word)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		if vocab.Value != req.Word {
			vocab, err = s.LookUpVocabLinear(ctx, req.QuizId, req.SentenceIndex, req.Word)
			if err != nil {
				return nil, err
			}
		}
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
