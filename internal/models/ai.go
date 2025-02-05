package models

import (
	"ielts-web-api/common"

	"gorm.io/datatypes"
)

type Vocab struct {
	ID          int            `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	VocabID     string         `json:"-" gorm:"column:vocab_id;not null"`
	Value       string         `json:"-" gorm:"column:value;not null"`
	WordDisplay string         `json:"word_display" gorm:"column:word_display;not null"`
	WordClass   string         `json:"word_class" gorm:"column:word_class;not null"`
	Meaning     string         `json:"meaning" gorm:"column:meaning;not null"`
	IPA         string         `json:"ipa" gorm:"column:ipa;not null"`
	Explanation string         `json:"explanation" gorm:"column:explanation;not null"`
	Collocation string         `json:"collocation" gorm:"column:collocation;not null"`
	Example     datatypes.JSON `json:"example" gorm:"column:example;not null"`
	VoteCount   int            `json:"-" gorm:"column:vote_count;not null"`
}

// TableName overrides the default table name for GORM
func (Vocab) TableName() string {
	return common.POSTGRES_TABLE_NAME_VOCAB
}

type LookUpVocabRequest struct {
	QuizId        int    `json:"quiz_id" binding:"required"`
	SentenceIndex int    `json:"sentence_index" binding:"required"`
	WordIndex     int    `json:"vocab_index" binding:"required"`
	Word          string `json:"word" binding:"required"`
}

// Define your models
type GPart struct {
        Text string `json:"text"`
}

type GContent struct {
        Parts []GPart `json:"parts"`
}

type GGenerationConfig struct {
        ResponseMimeType string `json:"responseMimeType,omitempty"`
}

type GeminiRequest struct {
        Contents         []GContent        `json:"contents"`
        GenerationConfig GGenerationConfig `json:"generationConfig,omitempty"`
}

type GVocabularyResponse struct {
        Vocabulary string `json:"vocabulary"`
        WordClass  string `json:"word_class"`
        Meaning    string `json:"meaning"`
        Example    string `json:"example"`
        IPA        string `json:"ipa"`
}

type GeminiResponse struct {
        Candidates []struct {
                Content struct {
                        Parts []GPart `json:"parts"`
                } `json:"content"`
        } `json:"candidates"`
}