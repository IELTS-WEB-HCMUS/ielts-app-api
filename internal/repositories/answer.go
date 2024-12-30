package repositories

import (
	"ielts-web-api/internal/models"

	"gorm.io/gorm"
)

type AnswerRepo struct {
	db *gorm.DB
	BaseRepository[models.Answer]
	Statistic BaseRepository[models.AnswerStatistic]
}

func NewAnswerRepository(db *gorm.DB) *AnswerRepo {
	baseRepo := NewBaseRepository[models.Answer](db)
	return &AnswerRepo{
		db:             db,
		BaseRepository: baseRepo,
		Statistic:      NewBaseRepository[models.AnswerStatistic](db),
	}
}
