package models

import (
	"ielts-web-api/common"
	"time"

	"gorm.io/datatypes"
)

type Question struct {
	ID                int            `json:"id" gorm:"id,primaryKey"`
	Type              string         `json:"type" gorm:"type"`
	QuestionType      string         `json:"question_type" gorm:"question_type"`
	Content           string         `json:"content" gorm:"content"`
	ContentWriting    string         `json:"content_writing"`
	Order             *int           `json:"order" gorm:"order"`
	Part              *int           `json:"part_id" gorm:"part_id"`
	Explain           *string        `json:"explain" gorm:"explain"`
	GapFillInBlank    *string        `json:"gap_fill_in_blank" gorm:"gap_fill_in_blank"`
	SingleChoiceRadio datatypes.JSON `json:"single_choice_radio" gorm:"single_choice_radio"`
	Selection         datatypes.JSON `json:"selection" gorm:"selection"`
	MultipleChoice    datatypes.JSON `json:"mutilple_choice" gorm:"column:mutilple_choice"`
	SelectionOption   datatypes.JSON `json:"selection_option" gorm:"selection_option"`
	UserCreated       string         `json:"user_created" gorm:"user_created"`
	UserUpdated       string         `json:"user_updated" gorm:"user_updated"`
	DateCreated       *time.Time     `json:"date_created" gorm:"date_created,autoCreateTime"`
	DateUpdated       *time.Time     `json:"date_updated" gorm:"date_updated,autoUpdateTime"`
}

func (Question) TableName() string {
	return common.POSTGRES_TABLE_NAME_QUESTION
}

type QuestionSuccessCount struct {
	Total   int `json:"total"`
	Success int `json:"success"`
	Failed  int `json:"failed"`
	Skip    int `json:"skip"`
}
