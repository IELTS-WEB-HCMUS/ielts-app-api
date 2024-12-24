package repositories

import (
	"ielts-web-api/internal/models"

	"gorm.io/gorm"
)

type UserVocabCategoryRepository struct {
	db *gorm.DB
	BaseRepository[models.UserVocabCategory]
}

func NewUserVocabCategoryRepository(db *gorm.DB) *UserVocabCategoryRepository {
	return &UserVocabCategoryRepository{
		db:             db,
		BaseRepository: NewBaseRepository[models.UserVocabCategory](db),
	}
}

type UserVocabBankRepository struct {
	db *gorm.DB
	BaseRepository[models.UserVocabBank]
}

func NewUserVocabBankRepository(db *gorm.DB) *UserVocabBankRepository {
	return &UserVocabBankRepository{
		db:             db,
		BaseRepository: NewBaseRepository[models.UserVocabBank](db),
	}
}
