package repositories

import (
	"ielts-web-api/internal/models"

	"gorm.io/gorm"
)

type TargetRepository struct {
	db *gorm.DB
	BaseRepository[models.Target]
}

func NewTargetRepository(db *gorm.DB) *TargetRepository {
	return &TargetRepository{
		db:             db,
		BaseRepository: NewBaseRepository[models.Target](db),
	}
}
