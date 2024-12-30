package repositories

import (
	"ielts-web-api/internal/models"

	"gorm.io/gorm"
)

type TagSearchRepository struct {
	db *gorm.DB
	BaseRepository[models.TagSearch]
}

func NewTagSearchRepository(db *gorm.DB) *TagSearchRepository {
	baseRepo := NewBaseRepository[models.TagSearch](db)
	return &TagSearchRepository{
		db:             db,
		BaseRepository: baseRepo,
	}
}

type TagSearchPositionRepo struct {
	db *gorm.DB
	BaseRepository[models.TagSearchPosition]
}

func NewTagSearchPositionRepo(db *gorm.DB) *TagSearchPositionRepo {
	baseRepo := NewBaseRepository[models.TagSearchPosition](db)
	return &TagSearchPositionRepo{
		db:             db,
		BaseRepository: baseRepo,
	}
}
