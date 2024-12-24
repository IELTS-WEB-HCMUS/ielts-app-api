package models

import (
	"ielts-web-api/common"
	"time"
)

type Part struct {
	ID                int         `json:"id" gorm:"id,primaryKey"`
	Quiz              int         `json:"quiz_id" gorm:"quiz"`
	Passage           int         `json:"passage" gorm:"passage"`
	Title             string      `json:"title" gorm:"title"`
	Sort              *int        `json:"sort" gorm:"sort"`
	Time              *int        `json:"time" gorm:"time"`
	Content           string      `json:"content" gorm:"content"`
	SimplifiedContent *string     `json:"simplified_content" gorm:"simplified_content"`
	UserCreated       string      `json:"user_created" gorm:"user_created"`
	DateCreated       *time.Time  `json:"date_created" gorm:"date_created"`
	DateUpdated       *time.Time  `json:"date_updated" gorm:"date_updated"`
	Questions         []*Question `json:"questions" gorm:"foreignKey:Part"`
	Quizzes           []*Quiz     `json:"-" gorm:"many2many:quiz_part;"`
	ListenFrom        *int        `json:"listen_from"`
	ListenTo          *int        `json:"listen_to"`
}

func (Part) TableName() string {
	return common.POSTGRES_TABLE_NAME_PART
}
