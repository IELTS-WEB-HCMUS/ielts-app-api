package services

import (
	"context"
	"ielts-web-api/common"
	"ielts-web-api/internal/models"
	"ielts-web-api/internal/repositories"

	"gorm.io/gorm"
)

func (s *Service) GetVocabCategoriresByUserId(ctx context.Context, userId string) ([]*models.UserVocabCategory, error) {
	params := models.QueryParams{
		Offset:    0,
		Limit:     4,
		QuerySort: models.QuerySort{},
		Selected:  []string{"id", "name"},
		Preload:   nil,
	}

	userClause := repositories.Clause(func(tx *gorm.DB) {
		tx.Where("user_id = ?", userId)
	})

	categories, err := s.vocabCategoriesRepo.List(ctx, params, userClause)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (s *Service) UpdateVocabCategory(ctx context.Context, req models.UpdateVocabCategoryRequest) (*string, error) {
	category, err := s.vocabCategoriesRepo.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	category.Name = req.Name
	_, err = s.vocabCategoriesRepo.Update(ctx, category.ID, category)

	if err != nil {
		return nil, err
	}
	successMessage := "Category updated successfully"
	return &successMessage, nil
}

func (s *Service) AddVocab(ctx context.Context, req models.UserVocabBankAddRequest) (*string, error) {
	newVocab := models.UserVocabBank{
		Value:     req.Value,
		WordClass: req.WordClass,
		Meaning:   req.Meaning,
		IPA:       req.IPA,
		Example:   req.Example,
		Note:      req.Note,
		Category:  req.Category,
	}

	_, err := s.userVocabBankRepo.Create(ctx, &newVocab)
	if err != nil {
		return nil, err
	}

	successMessage := "Add vocab successfully"
	return &successMessage, nil
}

func (s *Service) UpdateVocab(ctx context.Context, req models.UserVocabBankUpdateRequest) (*string, error) {
	vocab, err := s.userVocabBankRepo.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	// Update fields only if they are not null
	if req.Example != nil {
		vocab.Example = req.Example
	}
	if req.Note != nil {
		vocab.Note = req.Note
	}
	if req.Status != nil {
		vocab.Status = *req.Status
	}
	if req.Category != nil {
		vocab.Category = *req.Category
	}

	// Save updated vocab
	_, err = s.userVocabBankRepo.Update(ctx, vocab.ID, vocab)
	if err != nil {
		return nil, err
	}

	successMessage := "Update vocab successfully"
	return &successMessage, nil
}

func (s *Service) DeleteVocab(ctx context.Context, vocabID int) (*string, error) {
	// Check if the vocab exists using GetByID
	_, err := s.userVocabBankRepo.GetByID(ctx, vocabID)
	if err != nil {
		return nil, err
	}

	filter := repositories.Clause(func(tx *gorm.DB) {
		tx.Where("id = ?", vocabID)
	})

	err = s.userVocabBankRepo.Delete(ctx, filter)
	if err != nil {
		return nil, err
	}

	successMessage := "Delete vocab successfully"
	return &successMessage, nil
}

func (s *Service) GetVocabs(ctx context.Context, categoryID int) ([]*models.UserVocabBank, error) {
	params := models.QueryParams{
		Offset:    0,
		QuerySort: models.QuerySort{},
		Preload:   nil,
	}

	userClause := repositories.Clause(func(tx *gorm.DB) {
		tx.Where("category = ?", categoryID)
	})

	vocabularies, err := s.userVocabBankRepo.List(ctx, params, userClause)
	if len(vocabularies) == 0 {
		return nil, common.ErrCategoryNotFound
	}
	if err != nil {
		return nil, err
	}

	return vocabularies, nil
}
