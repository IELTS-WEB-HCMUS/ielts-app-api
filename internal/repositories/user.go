package repositories

import (
	"ielts-web-api/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
	BaseRepository[models.User]
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db:             db,
		BaseRepository: NewBaseRepository[models.User](db),
	}
}
