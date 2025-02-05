package models

import (
	"ielts-web-api/common"
	"time"
)

type Part struct {
	ID            int        `json:"id" gorm:"id,primaryKey"`
	UserCreated   string     `json:"user_created" gorm:"user_created"`
	DateCreated   *time.Time `json:"date_created" gorm:"date_created"`
	DateUpdated   *time.Time `json:"date_updated" gorm:"date_updated"`
	Title         string     `json:"title" gorm:"title"`
	Content       string     `json:"content" gorm:"content"`
	Description   string     `json:"description" gorm:"description"`
	QuestionCount int        `json:"question_count" gorm:"question_count"`
	Type          int        `json:"type" gorm:"type"`
	Level         int        `json:"level" gorm:"level"`
	Quiz          int        `json:"quiz_id" gorm:"quiz"`
	Questions     []Question `json:"questions" gorm:"foreignKey:Part;references:ID"`
	Passage       int        `json:"passage" gorm:"passage"`
}

func (Part) TableName() string {
	return common.POSTGRES_TABLE_NAME_PART
}
