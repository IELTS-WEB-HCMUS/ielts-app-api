package models

import (
	"ielts-web-api/common"
	"time"
)

type Quiz struct {
	ID          int          `json:"id" gorm:"id,primaryKey"`
	Status      string       `json:"status" gorm:"status"`
	UserCreated string       `json:"user_created" gorm:"user_created"`
	UserUpdated string       `json:"user_updated" gorm:"user_updated"`
	DateCreated *time.Time   `json:"date_created" gorm:"date_created"`
	DateUpdated *time.Time   `json:"date_updated" gorm:"date_updated"`
	Type        int          `json:"type" gorm:"type"`
	Content     *string      `json:"content" gorm:"content"`
	Title       string       `json:"title" gorm:"title"`
	Time        *int         `json:"time" gorm:"time"`
	Listening   *string      `json:"listening" gorm:"listening_file"`
	Level       *int         `json:"level,omitempty" gorm:"level"`
	VoteCount   int          `json:"vote_count" gorm:"vote_count"`
	Description *string      `json:"description,omitempty" gorm:"description"`
	Thumbnail   *string      `json:"thumbnail" gorm:"thumbnail"`
	Mode        int          `json:"mode" gorm:"mode"`
	IsPublic    bool         `json:"is_public" gorm:"is_public"`
	Parts       []*Part      `json:"parts" gorm:"many2many:quiz_part;"`
	TagSearches []*TagSearch `json:"tags" gorm:"many2many:quiz_tag_search;"`
	QuizPart    []QuizPart   `json:"quiz_part" gorm:"foreignKey:Quiz"`
	QuizPartMs  []QuizPartM  `json:"quiz_part_ms,omitempty" gorm:"foreignKey:QuizID"`
	IsSubmitted *bool        `json:"is_submitted,omitempty" gorm:"-"`
}

func (r Quiz) TableName() string {
	return common.POSTGRES_TABLE_NAME_QUIZ
}

type ListQuizzesParamsUri struct {
	BaseRequestParamsUri
	Status          *string `form:"status"`
	Type            *int    `form:"type" validate:"omitempty,min=1"`
	TagPassage      *int    `form:"tag_passage" validate:"omitempty,min=1"`
	TagSection      *int    `form:"tag_section" validate:"omitempty,min=1"`
	TagQuestionType *int    `form:"tag_question_type" validate:"omitempty,min=1"`
	Search          *string `form:"search"`
	Mode            *int    `form:"mode" validate:"omitempty,min=0"`
	SubmittedStatus int     `form:"submitted_status"`
}

type QuizSkill struct {
	ID       int    `json:"id" gorm:"id,primaryKey"`
	PublicId string `json:"public_id"`
}

func (r QuizSkill) TableName() string {
	return common.POSTGRES_TABLE_NAME_QUIZ_SKILL
}

type QuizParamsUri struct {
	QuizID int `uri:"quiz_id" binding:"required"`
}

type QuizAnswer struct {
	QuestionResult []QuestionResult `json:"question"`
	Answer         *Answer          `json:"answer"`
}

type QuizResult struct {
	SuccessCount int `json:"success_count"`
	Total        int `json:"total"`
}

type QuestionResult struct {
	Id int `json:"id"`
	QuizResult
}

type QuizDetail struct {
	QuizID             int    `json:"quiz_id"`
	QuizName           string `json:"quiz_name"`
	QuizCreatedAt      string `json:"quiz_created_at"`
	PartID             int    `json:"part_id"`
	PartQuiz           int    `json:"part_quiz"`
	PartPassage        string `json:"part_passage"`
	QuestionID         int    `json:"question_id"`
	QuestionPart       int    `json:"question_part"`
	QuestionType       string `json:"question_type"`
	QuestionTypeDetail string `json:"question_type_detail"`
	IsMultipleChoice   bool   `json:"is_multiple_choice"`
	IsGapFillInBlank   bool   `json:"is_gap_fill_in_blank"`
}
