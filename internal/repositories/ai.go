package repositories

import (
	"ielts-web-api/internal/models"

	"gorm.io/gorm"
)

type VocabRepository struct {
	db *gorm.DB
	BaseRepository[models.Vocab]
}

func NewVocabRepository(db *gorm.DB) *VocabRepository {
	return &VocabRepository{
		db:             db,
		BaseRepository: NewBaseRepository[models.Vocab](db),
	}
}
