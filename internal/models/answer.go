package models

import (
	"ielts-web-api/common"
	"math"
	"time"

	"gorm.io/datatypes"
)

type AnswerParamsUri struct {
	AnswerID int `uri:"answer_id" binding:"required"`
}

type AnswerSummary struct {
	Correct *int `json:"correct" gorm:"correct"`
	Total   *int `json:"total" gorm:"total"`
}

type AnswerSubmittedByQuiz struct {
	Quiz           int `json:"quiz"`
	TotalSubmitted int `json:"total_submitted"`
}

type Answer struct {
	ID                int               `json:"id" gorm:"id,primaryKey"`
	UserCreated       string            `json:"user_created" gorm:"user_created"`
	DateCreated       time.Time         `json:"date_created" gorm:"date_created"`
	Quiz              int               `json:"quiz" gorm:"quiz"`
	Detail            datatypes.JSON    `json:"detail" gorm:"type:jsonb"`
	Summary           datatypes.JSON    `json:"summary" gorm:"type:jsonb"`
	CompletedDuration int               `json:"completed_duration"`
	Questions         []*Question       `json:"questions" gorm:"-"`
	Student           *User             `json:"student" gorm:"-"`
	QuizDetail        *AnswerQuizDetail `json:"quiz_detail" gorm:"foreignKey:Quiz"`
}

type AnswerQuizDetail struct {
	ID    int    `json:"id" gorm:"id,primaryKey"`
	Title string `json:"title,omitempty"`
}

func (receiver AnswerQuizDetail) TableName() string {
	return common.POSTGRES_TABLE_NAME_QUIZ
}

type AnswerStatistic struct {
	ID                int       `json:"id" gorm:"id,primaryKey"`
	QuizID            int       `json:"quiz_id" gorm:"column:quiz"`
	UserCreated       string    `json:"user_created"`
	DateCreated       time.Time `json:"date_created"`
	QuizType          int       `json:"quiz_type"`
	QuizTitle         string    `json:"quiz_title,omitempty" gorm:"-"`
	QuizTopicFormat   string    `json:"topic_format,omitempty" gorm:"-"`
	SuccessCount      `json:"" gorm:"-"`
	SuccessQuizLog    *SuccessQuizLog `json:"-" gorm:"foreignKey:AnswerId"`
	CompletedDuration *int            `json:"completed_duration"`
	QuizDetail        *Quiz           `json:"-" gorm:"foreignKey:QuizID"`
	Type              int             `json:"type" gorm:"type"`
	BandScore         *float32        `json:"band_score"`
}

type AnswerStatistics []*AnswerStatistic

func (a *AnswerStatistic) Parse() *AnswerStatistic {
	if a.SuccessQuizLog != nil {
		a.SuccessCount = SuccessCount{
			Total:   a.SuccessQuizLog.Total,
			Success: a.SuccessQuizLog.Success,
			Skipped: a.SuccessQuizLog.Skipped,
			Failed:  a.SuccessQuizLog.Failed,
		}
		if a.SuccessQuizLog.Total > 0 {
			a.CorrectPercent = math.Round(((float64(a.Success)/float64(a.Total))*100)*100) / 100
		}
	}

	return a
}

func (l AnswerStatistics) Parse() AnswerStatistics {
	if len(l) == 0 {
		return l
	}
	for i := range l {
		l[i].Parse()
	}
	return l
}

func (a *Answer) TableName() string {
	return common.POSTGRES_TABLE_NAME_ANSWER
}

func (AnswerStatistic) TableName() string {
	return common.POSTGRES_TABLE_NAME_ANSWER
}
