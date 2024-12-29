package services

import (
	"context"
	"ielts-web-api/internal/models"
	"ielts-web-api/common"
)

func (s *Service) Vote(ctx context.Context, req models.VoteRequest) (*string, error){

	if req.Type == "quiz" {
		quiz, err := s.quizRepo.GetByID(ctx, req.ID)
		if err != nil {
			return nil, err
		}

		if (req.VoteType == "up") {
			quiz.VoteCount++
		} else if (req.VoteType == "down") {
			quiz.VoteCount--
		} else{
			return nil, common.ErrInvalidInput
		}

		_, err = s.quizRepo.UpdateColumns(ctx, req.ID, map[string]interface{}{
			"vote_count": quiz.VoteCount,
		})
		if err != nil {
			return nil, err
		}
	} else if req.Type == "vocab" {
		vocab, err := s.vocabRepo.GetByID(ctx, req.ID)
		if err != nil {
			return nil, err
		}

		if (req.VoteType == "up") {
			vocab.VoteCount++
		} else if (req.VoteType == "down") {
			vocab.VoteCount--
		} else{
			return nil, common.ErrInvalidInput
		}

		_, err = s.vocabRepo.UpdateColumns(ctx, req.ID, map[string]interface{}{
			"vote_count": vocab.VoteCount,
		}) 
		if err != nil {
			return nil, err
		}
	} else {
		return nil, common.ErrInvalidInput
	}

	successMessage := "Vote successfully"
	return &successMessage, nil
}
