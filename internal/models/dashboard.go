package models

import (
	"ielts-web-api/common"
	"math"
	"time"
)

type SuccessQuizLog struct {
	Id           int    `json:"id"`
	Passage      int    `json:"passage"`
	Total        int    `json:"total"`
	Success      int    `json:"success"`
	Skill        int    `json:"skill"`
	QuestionType string `json:"question_type"`
	UserId       string `json:"user_id"`
	Failed       int    `json:"failed"`
	Skipped      int    `json:"skipped"`
	AnswerId     int    `json:"answer_id"`
}

type SuccessDashboardData struct {
	DateCount         []SuccessCount `json:"date_count"`
	PassageCount      []SuccessCount `json:"passage_count"`
	QuestionTypeCount []SuccessCount `json:"question_type_count"`
}

type SuccessCount struct {
	Total          int     `json:"total"`
	Failed         int     `json:"failed"`
	Success        int     `json:"success"`
	Skipped        int     `json:"skipped"`
	CorrectPercent float64 `json:"correct_percent"`
	Date           string  `json:"date,omitempty"`
	Passage        int     `json:"passage,omitempty"`
	QuestionType   string  `json:"question_type,omitempty"`
}

type SuccessCounts []*SuccessCount

type AnswerStatisticsQuery struct {
	BaseRequestParamsUri
	SkillId         int       `form:"skill_id"` // Sort by skill
	Type            *int      `form:"type"`
	StartedAt       time.Time `form:"started_at"`
	EndedAt         time.Time `form:"ended_at"`
	Sort            string    `form:"sort"`
	QuizTypes       *[]int    `form:"quiz_types"`
	MockTestIDs     *[]int    `form:"mock_test_ids"`
	WritingTaskType *int      `form:"writing_task_type"`
}

func (SuccessQuizLog) TableName() string {
	return common.POSTGRES_TABLE_NAME_SUCCESS_QUIZ_LOG
}

func (SuccessCount) TableName() string {
	return common.POSTGRES_TABLE_NAME_SUCCESS_QUIZ_LOG
}

func (s *SuccessCount) Parse() *SuccessCount {
	if s.Total == 0 {
		s.CorrectPercent = 0
	} else {
		s.CorrectPercent = math.Round(((float64(s.Success)/float64(s.Total))*100)*100) / 100
	}

	return s
}

func (sl SuccessCounts) Parse() SuccessCounts {
	if len(sl) == 0 {
		return sl
	}

	for i, s := range sl {
		sl[i] = s.Parse()
	}

	return sl
}
