package repositories

import (
	"ielts-web-api/internal/models"

	"gorm.io/gorm"
)

type SuccessQuizLogRepo struct {
	db *gorm.DB
	BaseRepository[models.SuccessQuizLog]
	Statistic BaseRepository[models.SuccessCount]
}

func NewSuccessQuizLogRepository(db *gorm.DB) *SuccessQuizLogRepo {
	baseRepo := NewBaseRepository[models.SuccessQuizLog](db)
	return &SuccessQuizLogRepo{
		db:             db,
		BaseRepository: baseRepo,
		Statistic:      NewBaseRepository[models.SuccessCount](db),
	}
}
