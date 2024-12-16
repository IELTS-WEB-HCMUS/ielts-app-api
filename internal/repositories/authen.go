package repositories

import (
	"ielts-web-api/internal/models"

	"gorm.io/gorm"
)

type OTPRepository struct {
	BaseRepository[models.OTP]
}

func NewOTPRepository(db *gorm.DB) *OTPRepository {
	return &OTPRepository{
		BaseRepository: NewBaseRepository[models.OTP](db),
	}
}

type OTPAttemptRepository struct {
	BaseRepository[models.OTPAttempt]
}

func NewOTPAttemptRepository(db *gorm.DB) *OTPAttemptRepository {
	return &OTPAttemptRepository{
		BaseRepository: NewBaseRepository[models.OTPAttempt](db),
	}
}
