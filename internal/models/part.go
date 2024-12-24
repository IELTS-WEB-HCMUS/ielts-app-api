package models

import (
	"ielts-web-api/common"
	"time"
)

type Part struct {
	ID            int         `json:"id" gorm:"id,primaryKey"`
	Quiz          int         `json:"quiz_id" gorm:"quiz_id"`
	Title         string      `json:"title" gorm:"title"`
	Content       string      `json:"content" gorm:"content"`
	UserCreated   string      `json:"user_created" gorm:"user_created"`
	UserUpdated   string      `json:"user_updated" gorm:"user_updated"`
	DateCreated   *time.Time  `json:"date_created" gorm:"date_created"`
	DateUpdated   *time.Time  `json:"date_updated" gorm:"date_updated"`
	Description   string      `json:"description" gorm:"description"`
	QuestionCount int         `json:"question_count" gorm:"question_count"`
	Type          string      `json:"type" gorm:"type"`
	Level         string      `json:"level" gorm:"level"`
	Questions     []*Question `json:"questions" gorm:"foreignKey:Part"`
	Quizzes       []*Quiz     `json:"-" gorm:"many2many:quiz_part;"`
	ListenFrom    *int        `json:"listen_from"`
	ListenTo      *int        `json:"listen_to"`
}

func (Part) TableName() string {
	return common.POSTGRES_TABLE_NAME_PART
}
