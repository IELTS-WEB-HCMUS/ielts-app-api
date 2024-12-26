package models

import "ielts-web-api/common"

type UserVocabCategory struct {
	ID     int    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Name   string `json:"name" gorm:"column:name;type:text;not null"`
	UserID string `json:"-" gorm:"column:user_id;type:uuid;not null"`
}

// TableName overrides the default table name for GORM
func (UserVocabCategory) TableName() string {
	return common.POSTGRES_TABLE_NAME_USER_VOCAB_CATEGORY
}

type UpdateVocabCategoryRequest struct {
	Id   int    `json:"id" binding:"required"`
	Name string `json:"new_name" binding:"required"`
}

type UserVocabBank struct {
	ID        int     `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Value     string  `json:"value" gorm:"column:value;not null"`
	WordClass string  `json:"word_class" gorm:"column:word_class;not null"`
	Meaning   string  `json:"meaning" gorm:"column:meaning;not null"`
	IPA       string  `json:"ipa" gorm:"column:ipa;not null"`
	Example   *string `json:"example" gorm:"column:example;default:null"`
	Note      *string `json:"note" gorm:"column:note;default:null"`
	Status    string  `json:"status" gorm:"column:status;default:Chưa học"`
	Category  int     `json:"-" gorm:"column:category;not null"`
	CreatedAt string  `json:"-" gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
}

func (UserVocabBank) TableName() string {
	return common.POSTGRES_TABLE_NAME_USER_VOCAB_BANK
}

type UserVocabBankAddRequest struct {
	Value     string  `json:"value" binding:"required"`
	WordClass string  `json:"word_class" binding:"required"`
	Meaning   string  `json:"meaning" binding:"required"`
	IPA       string  `json:"ipa" binding:"required"`
	Example   *string `json:"example"`
	Note      *string `json:"note"`
	Category  int     `json:"category" binding:"required"`
}

type UserVocabBankUpdateRequest struct {
	Id       int     `json:"id" binding:"required"`
	Example  *string `json:"example"`
	Note     *string `json:"note"`
	Status   *string `json:"status"`
	Category *int    `json:"category"`
}

type UserVocabBankGetRequest struct {
	Category  int    `json:"category" binding:"required"`
	WordClass string `json:"word_class"`
	Status    string `json:"status"`
	Keyword   string `json:"keyword"`
	Page      int    `json:"page" binding:"required"`
	Limit     *int   `json:"limit"`
}

type UserVocabBankGetResponse struct {
	Vocabularies []*UserVocabBank `json:"vocabularies"`
	TotalItems   int64            `json:"total_items"`
	TotalPages   int              `json:"total_pages"`
	CurrentPage  int              `json:"current_page"`
}
