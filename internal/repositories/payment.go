package repositories

import (
	"ielts-web-api/internal/models"

	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
	BaseRepository[models.Payment]
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{
		db:             db,
		BaseRepository: NewBaseRepository[models.Payment](db),
	}
}
