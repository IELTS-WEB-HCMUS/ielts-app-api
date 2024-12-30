package services

import (
	"context"
	"errors"
	"ielts-web-api/common"
	"ielts-web-api/internal/models"
	"ielts-web-api/internal/repositories"
	"sort"

	"gorm.io/gorm"
)

func (s *Service) GetAnswer(ctx context.Context, userID string, answerID int) (*models.Answer, error) {
	// Get detail answer
	conds := []repositories.Clause{
		func(tx *gorm.DB) {
			tx.Where("id", answerID)
		},
		func(tx *gorm.DB) {
			ps := []common.Preload{
				{
					Model:    "QuizDetail",
					Selected: []string{"id", "title"},
				},
			}

			for _, p := range ps {
				common.ApplyPreload(tx, p)
			}
		},
	}

	answer, err := s.answerRepo.GetDetailByConditions(ctx, conds...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrRecordNotFound
		}
		return nil, err
	}

	// get student info
	student, err := s.userRepo.GetDetailByConditions(ctx, func(tx *gorm.DB) {
		tx.Select("id, first_name, last_name, avatar").Where("id = ?", answer.UserCreated)
	})
	if err != nil {
		return nil, err
	}
	answer.Student = student

	return answer, nil
}

func (s *Service) GetAnswerStatistic(ctx context.Context, studentID string, request *models.AnswerStatisticsQuery) (interface{}, error) {
	filters := []repositories.Clause{}
	var (
		statisticsByQuiz            models.AnswerStatistics
		statisticsByPassageOrQsType models.SuccessCounts
		err                         error
	)

	if request.Type == nil {
		return nil, nil
	}

	var mappingType = map[int]string{
		common.AnswerStatisticByPassage:    "passage",
		common.AnswerStatisticQuestionType: "question_type",
		common.AnswerStatisticByQuiz:       "answer_id",
	}

	filterSuccessQuizLog := func(tx *gorm.DB) *gorm.DB {
		return tx.Select(
			mappingType[*request.Type],
			"sum(total) as total",
			"sum(success) as success",
			"sum(failed) as failed",
			"sum(skipped) as skipped",
		).Where(mappingType[*request.Type] + " IS NOT NULL").Group(mappingType[*request.Type])
	}

	if *request.Type == common.AnswerStatisticByQuiz {
		filters = append(filters, func(tx *gorm.DB) {
			tx.Where("user_created = ?", studentID)
		})
		if request.SkillId > 0 {
			filters = append(filters, func(tx *gorm.DB) {
				tx.Where("type = ?", request.SkillId)
			})
		}
		page, pageSize := common.GetPageAndPageSize(request.Page, request.PageSize)
		total, err := s.answerRepo.Count(ctx, models.QueryParams{}, filters...)
		if err != nil {
			return nil, err
		}
		resData := models.BaseListResponse{
			Total:    int(total),
			Page:     page,
			PageSize: pageSize,
			Items:    []*models.AnswerStatistic{},
		}

		if total == 0 {
			return &resData, nil
		}

		if len(request.Sort) == 0 {
			request.Sort = "date_created.desc"
		}

		filters = append(filters, func(tx *gorm.DB) {
			tx.Preload("SuccessQuizLog", func(tx *gorm.DB) *gorm.DB {
				return tx.Select(
					mappingType[*request.Type],
					"sum(total) as total",
					"sum(success) as success",
					"sum(failed) as failed",
					"sum(skipped) as skipped",
				).Where(mappingType[*request.Type] + " IS NOT NULL AND question_type !=''").Group(mappingType[*request.Type])
			}).Preload("QuizDetail")
		})

		statisticsByQuiz, err = s.answerRepo.Statistic.List(ctx, models.QueryParams{
			Limit:  pageSize,
			Offset: (page - 1) * pageSize,
			QuerySort: models.QuerySort{
				Origin: request.Sort,
			},
		}, filters...)

		if err != nil {
			return nil, err
		}
		for _, item := range statisticsByQuiz {
			if item.QuizDetail != nil {
				item.QuizTitle = item.QuizDetail.Title
			}
		}
		resData.Items = statisticsByQuiz.Parse()
		return &resData, nil
	} else if *request.Type == common.AnswerStatisticByPassage || *request.Type == common.AnswerStatisticQuestionType {
		if request.SkillId > 0 {
			filters = append(filters, func(tx *gorm.DB) {
				tx.Where("skill = ?", request.SkillId)
			})
		}
		filters = append(filters, func(tx *gorm.DB) {
			tx.Where("user_id = ? ", studentID)
		})
		filters = append(filters, func(tx *gorm.DB) {
			filterSuccessQuizLog(tx)
		})
		if *request.Type == common.AnswerStatisticByPassage {
			filters = append(filters, func(tx *gorm.DB) {
				tx.Where(mappingType[*request.Type] + " != 0") // passage != 0
			})
		} else {
			filters = append(filters, func(tx *gorm.DB) {
				tx.Where(mappingType[*request.Type] + " != ''") // question_type != ''
			})
		}
		statisticsByPassageOrQsType, err = s.successQuizLogRepo.Statistic.List(ctx, models.QueryParams{
			QuerySort: models.QuerySort{
				Origin: request.Sort,
			},
		}, filters...)
		if err != nil {
			return nil, err
		}
		statisticsByPassageOrQsType = statisticsByPassageOrQsType.Parse()

		// default sort: correct percent asc
		if len(request.Sort) == 0 {
			sort.Slice(statisticsByPassageOrQsType, func(i, j int) bool {
				return statisticsByPassageOrQsType[i].CorrectPercent < statisticsByPassageOrQsType[j].CorrectPercent
			})
		}
		return &models.BaseListResponse{
			Total:    len(statisticsByPassageOrQsType),
			Page:     1,
			PageSize: len(statisticsByPassageOrQsType),
			Items:    statisticsByPassageOrQsType,
		}, nil
	}
	return nil, nil
}
