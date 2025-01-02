package services

import (
	"context"
	"ielts-web-api/internal/models"
	"ielts-web-api/internal/repositories"

	"fmt"
	"math"

	"gorm.io/gorm"
)

func (s *Service) GetVocabCategoriresByUserId(ctx context.Context, userId string) ([]*models.UserVocabCategory, error) {
	params := models.QueryParams{
		Offset:   0,
		Limit:    4,
		Selected: []string{"id", "name"},
		Preload:  nil,
	}

	userClause := repositories.Clause(func(tx *gorm.DB) {
		tx.Where("user_id = ?", userId)
		tx.Order("id ASC")
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

	if req.Meaning == "" || req.Example == nil || req.IPA == "" {
		geminiVocab, err := genVocabInfoWithGemini(req.Value, req.WordClass)

		if err != nil {
			return nil, err
		}
		if req.Meaning == "" {
			newVocab.Meaning = geminiVocab.Meaning
		}
		if req.Example == nil {
			newVocab.Example = &geminiVocab.Example
		}
		if req.IPA == "" {
			newVocab.IPA = geminiVocab.IPA
		}
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
	if req.Meaning != nil {
		vocab.Meaning = *req.Meaning
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

func (s *Service) GetVocabs(ctx context.Context, req models.UserVocabBankGetRequest) (*models.UserVocabBankGetResponse, error) {
	params := models.QueryParams{}

	// Set Pagination Parameters
	if req.Limit != nil {
		params.Limit = *req.Limit
	}

	// Calculate Offset from Page
	params.Offset = (req.Page - 1) * params.Limit

	// Define Filter Clause
	userClause := repositories.Clause(func(tx *gorm.DB) {
		// Mandatory Category Filter
		tx.Where("category = ?", req.Category)

		// Optional WordClass Filter
		if req.WordClass != "" {
			tx.Where("word_class = ?", req.WordClass)
		}

		// Optional Status Filter
		if req.Status != "" {
			tx.Where("status = ?", req.Status)
		}

		// Keyword Search (Search in title and description)
		if req.Keyword != "" {
			searchPattern := fmt.Sprintf("%%%s%%", req.Keyword)
			tx.Where("value ILIKE ?", searchPattern)
		}

		// Apply Sorting (Default to `created_at DESC`)
		tx.Order("created_at DESC")
	})

	// Fetch Total Count
	totalItems, err := s.userVocabBankRepo.Count(ctx, models.QueryParams{Offset: 0}, userClause)
	if err != nil {
		return nil, err
	}

	// Fetch Paginated Records
	vocabularies, err := s.userVocabBankRepo.List(ctx, params, userClause)
	if err != nil {
		return nil, err
	}

	// Calculate Total Pages
	totalPages := 0
	if params.Limit > 0 {
		totalPages = int(math.Ceil(float64(totalItems) / float64(params.Limit)))
	}
	if totalPages == 0 {
		totalPages = 1
	}

	// Calculate Current Page
	currentPage := req.Page

	// Build Response
	response := &models.UserVocabBankGetResponse{
		Vocabularies: vocabularies,
		TotalItems:   totalItems,
		TotalPages:   totalPages,
		CurrentPage:  currentPage,
	}

	return response, nil
}
