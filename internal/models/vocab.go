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
