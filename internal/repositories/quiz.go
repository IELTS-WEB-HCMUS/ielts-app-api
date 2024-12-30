package repositories

import (
	"context"
	"errors"
	"fmt"
	"ielts-web-api/common"
	"ielts-web-api/internal/models"
	"strings"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type QuizRepo struct {
	db *gorm.DB
	BaseRepository[models.Quiz]
}

type QuizSkillRepo struct {
	db *gorm.DB
	BaseRepository[models.QuizSkill]
}

func NewQuizRepository(db *gorm.DB) *QuizRepo {
	baseRepo := NewBaseRepository[models.Quiz](db)
	return &QuizRepo{
		db:             db,
		BaseRepository: baseRepo,
	}
}

func NewQuizSkillRepository(db *gorm.DB) *QuizSkillRepo {
	baseRepo := NewBaseRepository[models.QuizSkill](db)
	return &QuizSkillRepo{
		db:             db,
		BaseRepository: baseRepo,
	}
}

func (r *QuizRepo) GetQuizIDsInCludeTagIDs(ctx context.Context, tagIDs []int) ([]int, error) {
	var qIDs []int

	tagIDsFmt := strings.ReplaceAll(fmt.Sprintf("%+v", tagIDs), " ", ", ")

	tx := r.db.Table(common.POSTGRES_TABLE_NAME_QUIZ_TAG_SEARCH).
		Select("quiz_id").
		Group("quiz_id").
		Where("quiz_id IS NOT NULL").
		Having(fmt.Sprintf("ARRAY%+v <@ ARRAY_AGG(tag_search_id)", tagIDsFmt))
	err := tx.Pluck("quiz_id", &qIDs).Error
	if err != nil {
		return nil, err
	}
	return qIDs, nil
}

func (r *QuizRepo) GetQuizSubmitted(ctx context.Context, quizId int) (*models.Quiz, error) {
	var results []struct {
		QuizID                 int
		QuizTitle              string
		QuizType               int
		PartID                 int
		QuizIDPart             int
		PartPassage            int
		QuestionID             int
		QuestionPartID         int
		QuestionType           string
		QuestionQuestionType   string
		QuestionMultiple       datatypes.JSON
		QuestionGapFillInBlank *string
	}

	// Perform a single query with JOINs
	err := r.db.Raw(`
		SELECT
            q.id AS quiz_id, q.title AS quiz_title, q.type AS quiz_type,
            p.id AS part_id, p.quiz AS quiz_id_part, p.passage AS part_passage,
            qu.id AS question_id, qu.part AS question_part_id, qu.question_type AS question_type,
		    qu.multiple_choice AS question_multiple, qu.gap_fill_in_blank AS question_gap_fill_in_blank,
			qu.type AS question_question_type
        FROM
            public.quiz q
        LEFT JOIN
            public.part p ON q.id = p.quiz
        LEFT JOIN
            public.question qu ON p.id = qu.part
        WHERE
            q.id = ?
	`, quizId).Scan(&results).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrQuizNotFound
		}
		return nil, err
	}

	quiz := &models.Quiz{
		ID:    quizId,
		Title: results[0].QuizTitle,
		Type:  results[0].QuizType,
		Parts: []*models.Part{},
	}

	partMap := make(map[int]*models.Part)
	for _, result := range results {
		if _, exists := partMap[result.PartID]; !exists && result.PartID != 0 {
			part := &models.Part{
				ID:        result.PartID,
				Quiz:      result.QuizIDPart,
				Passage:   result.PartPassage,
				Questions: []models.Question{},
			}
			quiz.Parts = append(quiz.Parts, part)
			partMap[result.PartID] = part
		}

		if result.QuestionID != 0 {
			question := &models.Question{
				ID:             result.QuestionID,
				Part:           &result.QuestionPartID,
				Type:           result.QuestionQuestionType,
				QuestionType:   result.QuestionType,
				MultipleChoice: result.QuestionMultiple,
				GapFillInBlank: func() datatypes.JSON {
					if result.QuestionGapFillInBlank != nil {
						return datatypes.JSON([]byte(*result.QuestionGapFillInBlank))
					}
					return datatypes.JSON([]byte("null"))
				}(),
			}
			if part, exists := partMap[result.PartID]; exists {
				part.Questions = append(part.Questions, *question)
			}
		}
	}

	return quiz, nil
}
