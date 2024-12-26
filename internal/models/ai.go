package models

import (
	"ielts-web-api/common"
)

type Vocab struct {
	ID            int    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	VocabID       string `json:"-" gorm:"column:vocab_id;not null"`
	Value         string `json:"value" gorm:"column:value;not null"`
	WordClass     string `json:"word_class" gorm:"column:word_class;not null"`
	Meaning       string `json:"meaning" gorm:"column:meaning;not null"`
	IPA           string `json:"ipa" gorm:"column:ipa;not null"`
	Explanation   string `json:"explanation" gorm:"column:explanation;not null"`
	Collocation   string `json:"collocation" gorm:"column:collocation;not null"`
	VerbStructure string `json:"verb_structure" gorm:"column:verb_structure;not null"`
}

// TableName overrides the default table name for GORM
func (Vocab) TableName() string {
	return common.POSTGRES_TABLE_NAME_VOCAB
}

type LookUpVocabRequest struct {
	QuizId     int `json:"quiz_id" binding:"required"`
	VocabIndex int `json:"vocab_index" binding:"required"`
}
