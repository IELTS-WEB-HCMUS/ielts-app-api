package models

import (
	"ielts-web-api/common"
	"time"
)

type TargetRequest struct {
	TargetStudyDuration *int     `json:"target_study_duration"`
	TargetReading       *float32 `json:"target_reading"`
	TargetListening     *float32 `json:"target_listening"`
	TargetSpeaking      *float32 `json:"target_speaking"`
	TargetWriting       *float32 `json:"target_writing"`
	NextExamDate        *string  `json:"next_exam_date"`
}

type Target struct {
	ID                  string    `json:"id" gorm:"type:uuid;primaryKey"`
	TargetStudyDuration int       `json:"TargetStudyDuration" gorm:"column:target_study_duration"`
	TargetReading       float32   `json:"TargetReading" gorm:"column:target_reading"`
	TargetListening     float32   `json:"TargetListening" gorm:"column:target_listening"`
	TargetSpeaking      float32   `json:"TargetSpeaking" gorm:"column:target_speaking"`
	TargetWriting       float32   `json:"TargetWriting" gorm:"column:target_writing"`
	NextExamDate        time.Time `json:"NextExamDate" gorm:"column:next_exam_date"`
}

func (Target) TableName() string {
	// define common variable here
	return common.POSTGRES_TABLE_NAME_TARGETS
}
