package models

import (
	"encoding/json"
	"fmt"
	"ielts-web-api/common"
)

type QuizPartM struct {
	ID     int     `json:"id"`
	QuizID int     `json:"quiz_id"`
	PartID int     `json:"part_id"`
	Part   *PartV2 `json:"part,omitempty" gorm:"foreignKey:PartID"`
}

type PartV2 struct {
	ID            int `json:"id"`
	QuestionCount int `json:"question_count"`
	Quiz          int `json:"quiz" gorm:"column:quiz"`
}

func (PartV2) TableName() string {
	return common.POSTGRES_TABLE_NAME_PART
}

func (QuizPartM) TableName() string {
	return common.POSTGRES_TABLE_NAME_QUIZ_PART
}

type QuizPart struct {
	Id       int        `json:"id"`
	Quiz     int        `json:"quiz"`
	Question []Question `json:"question,omitempty" gorm:"foreignKey:Part"`
	Type     int        `json:"type"`
	Quizzes  []*Quiz    `json:"-" gorm:"many2many:quiz_part;"`
}

func (r QuizPart) TableName() string {
	return common.POSTGRES_TABLE_NAME_PART

}

type QuestionMultipleChoice struct {
	Text    string `json:"text"`
	Correct bool   `json:"correct"`
	Order   int    `json:"order"`
	Explain string `json:"explain"`
}

func (r Question) CountTotalSubQuestion() int {
	switch r.Type {
	case common.QUESTION_TYPE_SINGLE_RADIO:
		return 1
	case common.QUESTION_TYPE_SINGLE_SELECTION:
		return 1
	case common.QUESTION_TYPE_FILL_IN_THE_BLANK:
		return 1
	case common.QUESTION_TYPE_MULTIPLE:
		var choices []QuestionMultipleChoice
		if r.MultipleChoice == nil {
			return 0
		}
		err := json.Unmarshal(r.MultipleChoice, &choices)
		if err != nil {
			return 0
		}
		count := 0
		for _, c := range choices {
			if c.Correct {
				count++
			}
		}
		fmt.Println("choices", count, choices)
		return count
	}
	return 0
}
