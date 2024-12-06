package services

import (
	"context"
	"errors"
	"fmt"
	"ielts-web-api/common"
	"ielts-web-api/internal/models"
	"time"

	"gorm.io/gorm"
)

func (s *Service) GetTargetById(ctx context.Context, id string) (*models.Target, error) {
	fmt.Println(id)
	target, err := s.TargetRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return target, nil
}

func (s *Service) CreateTarget(ctx context.Context, userId string, req models.TargetRequest) (*models.Target, error) {
	// Check if a target for the user already exists
	_, err := s.TargetRepo.GetByID(ctx, userId)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	} else {
		err = common.ErrTargetAlreadyExists
		return nil, err
	}

	// Parse NextExamDate if provided, else use a default date or handle it as needed
	var parsedNextExamDate time.Time
	if req.NextExamDate != nil {
		layout := "2006-01-02 15:04:05"
		parsedNextExamDate, err = time.Parse(layout, *req.NextExamDate)
		if err != nil {
			return nil, err
		}
	} else {
		// Use a zero time value or another default value if NextExamDate is not provided
		parsedNextExamDate = time.Time{}
	}

	// Assign default values for other fields if they are nil
	newTarget := models.Target{
		ID:                  userId,
		TargetStudyDuration: getOrDefaultInt(req.TargetStudyDuration, 0),
		TargetReading:       getOrDefaultFloat32(req.TargetReading, 0.0),
		TargetListening:     getOrDefaultFloat32(req.TargetListening, 0.0),
		TargetSpeaking:      getOrDefaultFloat32(req.TargetSpeaking, 0.0),
		TargetWriting:       getOrDefaultFloat32(req.TargetWriting, 0.0),
		NextExamDate:        parsedNextExamDate,
	}

	// Create the new target record
	createdTarget, err := s.TargetRepo.Create(ctx, &newTarget)
	if err != nil {
		return nil, err
	}
	return createdTarget, nil
}

func getOrDefaultInt(val *int, defaultVal int) int {
	if val != nil {
		return *val
	}
	return defaultVal
}

func getOrDefaultFloat32(val *float32, defaultVal float32) float32 {
	if val != nil {
		return *val
	}
	return defaultVal
}

func (s *Service) UpdateTarget(ctx context.Context, userId string, req models.TargetRequest) (*models.Target, error) {
	updateTarget, err := s.TargetRepo.GetByID(ctx, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}

	if req.TargetStudyDuration != nil {
		updateTarget.TargetStudyDuration = *req.TargetStudyDuration
	}
	if req.TargetReading != nil {
		updateTarget.TargetReading = *req.TargetReading
	}
	if req.TargetListening != nil {
		updateTarget.TargetListening = *req.TargetListening
	}
	if req.TargetSpeaking != nil {
		updateTarget.TargetSpeaking = *req.TargetSpeaking
	}
	if req.TargetWriting != nil {
		updateTarget.TargetWriting = *req.TargetWriting
	}
	if req.NextExamDate != nil {
		layout := "2006-01-02"
		parsedNextExamDate, err := time.Parse(layout, *req.NextExamDate)
		if err != nil {
			return nil, err
		}
		updateTarget.NextExamDate = parsedNextExamDate
	}
	updatedTarget, err := s.TargetRepo.Update(ctx, userId, updateTarget)
	if err != nil {
		return nil, err
	}
	return updatedTarget, nil
}
