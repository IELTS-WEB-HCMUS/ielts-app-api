package services

import (
	"context"
	"fmt"
	"ielts-web-api/common"
	"ielts-web-api/internal/models"
	"os"

	"errors"

	"gorm.io/gorm"

	"encoding/json"

	"github.com/go-resty/resty/v2"
)

func (s *Service) LookUpVocabLinear(ctx context.Context, quizId int, sentenceIndex int, word string) (*models.Vocab, error) {
	vocabIdPattern := fmt.Sprintf("%d_%d_%%", quizId, sentenceIndex)

	vocab, err := s.vocabRepo.GetDetailByConditions(ctx, func(tx *gorm.DB) {
		tx.Where("vocab_id LIKE ? AND value = ?", vocabIdPattern, word)
	})
	if err != nil {
		return nil, err
	}

	return vocab, nil
}

func (s *Service) LookUpVocab(ctx context.Context, req models.LookUpVocabRequest, userId string) (*models.Vocab, error) {
	vocabId := fmt.Sprintf("%d_%d_%d", req.QuizId, req.SentenceIndex, req.WordIndex)

	vocab, err := s.vocabRepo.GetDetailByConditions(ctx, func(tx *gorm.DB) {
		tx.Where("vocab_id = ?", vocabId)
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			vocab, err = s.LookUpVocabLinear(ctx, req.QuizId, req.SentenceIndex, req.Word)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		if vocab.Value != req.Word {
			vocab, err = s.LookUpVocabLinear(ctx, req.QuizId, req.SentenceIndex, req.Word)
			if err != nil {
				return nil, err
			}
		}
	}

	user, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	user.VocabUsageCount--
	_, err = s.userRepo.UpdateColumns(ctx, user.ID, map[string]interface{}{
		"vocab_usage_count": user.VocabUsageCount,
	})
	if err != nil {
		return nil, err
	}

	return vocab, nil
}

func (s *Service) CheckVocabUsageCount(ctx context.Context, userId string) error {
	user, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		return err
	}
	if user.VocabUsageCount <= 0 {
		return common.ErrVocabUsageCountExceeded
	}

	return nil
}

// genVocabInfoWithGemini fetches vocabulary information from Gemini
func genVocabInfoWithGemini(vocab string, wordClass string) (*models.GVocabularyResponse, error) {
	APIKey := os.Getenv("GEMINI_API_KEY")
	APIEndpoint := os.Getenv("GEMINI_API_ENDPOINT")

	if APIKey == "" || APIEndpoint == "" {
		return nil, errors.New("API key or endpoint is not set in environment variables")
	}

	prompt := fmt.Sprintf(
		`Provide the following details in JSON format for the given vocabulary, return JSON only without any extra text. Ensure all fields are included and valid JSON is returned:

Vocabulary: '%s'
Word Class: '%s'

The "meaning" field should contain the Vietnamese translation/meaning.

Response JSON example:
{
        "vocabulary": "run",
        "word_class": "verb",
        "meaning": "chạy", // Vietnamese meaning
        "example": "He runs every morning to stay fit.",
        "ipa": "/rʌn/"
}`,
		vocab, wordClass,
	)

	// Create a new HTTP client
	client := resty.New()

	// Prepare the request payload
	request := models.GeminiRequest{
		Contents: []models.GContent{
			{
				Parts: []models.GPart{{Text: prompt}},
			},
		},
	}

	// Make the API call
	var response models.GeminiResponse
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetQueryParam("key", APIKey).
		SetBody(request).
		SetResult(&response).
		Post(APIEndpoint)

	// Handle network or HTTP errors
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err) // Wrap the error
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s, Body: %s", resp.Status(), resp.Body()) // Include body in error
	}

	// Parse JSON response from Gemini
	if len(response.Candidates) > 0 && len(response.Candidates[0].Content.Parts) > 0 {
		rawText := response.Candidates[0].Content.Parts[0].Text

		// Sanitize JSON (important!)
		rawText = sanitizeJSON(rawText)

		var vocabResponse models.GVocabularyResponse
		err := json.Unmarshal([]byte(rawText), &vocabResponse)
		if err != nil {
			return nil, fmt.Errorf("failed to parse JSON response: %v\nRaw Response: %s", err, rawText)
		}

		return &vocabResponse, nil
	}

	return nil, errors.New("no candidates found in the response")
}

func sanitizeJSON(jsonString string) string {
	// Remove any leading/trailing non-JSON characters
	var jsonData interface{}
	json.Unmarshal([]byte(jsonString), &jsonData)       // Try parsing to check if it's valid JSON
	if _, ok := jsonData.(map[string]interface{}); ok { // If it's a map, it's valid JSON
		return jsonString
	}
	for i := 0; i < len(jsonString); i++ {
		if jsonString[i] == '{' {
			for j := len(jsonString) - 1; j >= 0; j-- {
				if jsonString[j] == '}' {
					return jsonString[i : j+1]
				}
			}
		}

	}

	return ""
}
