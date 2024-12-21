package services

import (
	"context"
	"errors"
	"fmt"
	"ielts-web-api/common"
	"ielts-web-api/internal/models"
	"ielts-web-api/internal/repositories"
	"time"

	"gorm.io/gorm"
)

func (s *Service) GetQuizzes(ctx context.Context, userID string, request *models.ListQuizzesParamsUri) (*models.BaseListResponse, error) {
	var (
		filters = []repositories.Clause{}
		quizIDs = []int{}
		err     error
	)
	page, pageSize := common.GetPageAndPageSize(request.Page, request.PageSize)

	resData := models.BaseListResponse{
		Total:    0,
		Page:     page,
		PageSize: pageSize,
		Items:    []*models.Quiz{},
	}

	if request.TagPassage != nil ||
		request.TagSection != nil ||
		request.TagQuestionType != nil ||
		request.TagTopic != nil ||
		request.TagBookType != nil {

		tagIDs := []int{}
		if request.TagSection != nil {
			tagIDs = append(tagIDs, *request.TagSection)
		}
		if request.TagPassage != nil {
			tagIDs = append(tagIDs, *request.TagPassage)
		}
		if request.TagQuestionType != nil {
			tagIDs = append(tagIDs, *request.TagQuestionType)
		}

		if request.TagTopic != nil {
			tagIDs = append(tagIDs, *request.TagTopic)
		}

		if request.TagBookType != nil {
			tagIDs = append(tagIDs, *request.TagBookType)
		}

		// get quizIDs have matched tags
		quizIDs, err = s.quizRepo.GetQuizIDsInCludeTagIDs(ctx, tagIDs)
		if err != nil {
			return nil, err
		}

		if len(quizIDs) == 0 {
			return &resData, nil
		}

		filters = append(filters, func(tx *gorm.DB) {
			tx.Where("quiz.id IN ?", quizIDs)
		})
	}

	if request.IsTest != nil {
		filters = append(filters, func(tx *gorm.DB) {
			tx.Where("quiz.is_test = ?", *request.IsTest)
		})
	}

	if request.Type != nil {
		filters = append(filters, func(tx *gorm.DB) {
			tx.Where("quiz.type = ?", *request.Type)
		})
	}

	if request.Status != nil && len(*request.Status) > 0 {
		filters = append(filters, func(tx *gorm.DB) {
			tx.Where("quiz.status = ?", *request.Status)
		})
	}

	if request.Search != nil && len(*request.Search) > 0 {
		var quesQuizIDs []int

		if len(quesQuizIDs) > 0 {
			filters = append(filters, func(tx *gorm.DB) {
				tx.Where("quiz.title ILIKE ? OR id IN (?)", "%"+*request.Search+"%", quesQuizIDs)
			})
		} else {
			filters = append(filters, func(tx *gorm.DB) {
				tx.Where("quiz.title ILIKE ?", "%"+*request.Search+"%")
			})
		}
	}

	// Handle information with log-in student
	var (
		quizSubmittedIDs = []int{}
		quizSubmittedMap = make(map[int]bool)
	)

	if len(userID) > 0 {
		filterQuizIDsSubmitted := func(tx *gorm.DB) {
			if len(quizIDs) > 0 {
				tx.Select("distinct quiz").Where("user_created = ? and quiz IN ?", userID, quizIDs)
			} else {
				tx.Select("distinct quiz").Where("user_created = ?", userID)
			}
		}

		answers, err := s.answerRepo.List(ctx, models.QueryParams{}, filterQuizIDsSubmitted)
		if err != nil {
			return nil, err
		}

		for _, answer := range answers {
			quizSubmittedMap[answer.Quiz] = true
			quizSubmittedIDs = append(quizSubmittedIDs, answer.Quiz)
		}

		if request.SubmittedStatus == common.QuizSubmittedStatusYes {
			filters = append(filters, func(tx *gorm.DB) {
				tx.Where("quiz.id IN ?", quizSubmittedIDs)
			})
		} else if request.SubmittedStatus == common.QuizSubmittedStatusNo {
			if len(quizSubmittedIDs) > 0 {
				filters = append(filters, func(tx *gorm.DB) {
					tx.Where("quiz.id NOT IN ?", quizSubmittedIDs)
				})
			}
		}
	}

	total, err := s.quizRepo.Count(ctx, models.QueryParams{}, filters...)
	if err != nil {
		return nil, err
	}

	if total == 0 {
		return &resData, nil
	}

	resData.Total = int(total)

	// Preload tagSearches
	filters = append(filters, func(tx *gorm.DB) {
		tx.Preload("TagSearches")
	})

	records, err := s.quizRepo.List(
		ctx,
		models.QueryParams{
			Limit:  pageSize,
			Offset: (page - 1) * pageSize,
			QuerySort: models.QuerySort{
				Origin: request.Sort,
			},
		},
		filters...,
	)

	if err != nil {
		return nil, err
	}

	quizIDs = []int{}
	for _, record := range records {
		quizIDs = append(quizIDs, record.ID)
		_, submitted := quizSubmittedMap[record.ID]
		record.IsSubmitted = &submitted
	}

	resData.Items = records
	return &resData, nil
}

func (s *Service) GetQuiz(ctx context.Context, req *models.QuizParamsUri, userID string) (*models.Quiz, error) {
	var (
		quiz *models.Quiz
		err  error
	)
	filters := []repositories.Clause{}
	filters = append(filters, func(tx *gorm.DB) {
		tx.Preload("Parts", func(db *gorm.DB) *gorm.DB {
			return db.Joins("INNER JOIN quiz_part ON quiz_part.quiz_id = ? AND quiz_part.part_id = part.id", req.QuizID).Order("quiz_part.sort")
		}).Preload("Parts.Questions", func(db *gorm.DB) *gorm.DB {
			return db.Order("question.sort")
		}).Where("id", req.QuizID)
	})

	// Preload vocabs: Milestone 2
	quiz, err = s.quizRepo.GetDetailByConditions(
		ctx,
		filters...,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrRecordNotFound
		}
		return nil, err
	}

	return quiz, nil
}

func (s *Service) SubmitQuizAnswer(ctx context.Context, userId string, quizId int, results models.QuizAnswer) (answer *models.Answer, err error) {
	// get quiz & check type
	q, err := s.quizRepo.GetByID(ctx, quizId)
	if err != nil {
		return nil, err
	}

	t, err := s.quizSkillRepo.GetByID(ctx, q.Type)
	if err != nil {
		return nil, err
	}

	results.Answer.QuizType = q.QuizType

	if t.PublicId == common.QuizSkillReading {
		answer, _, err = s.submitReadingListening(ctx, userId, quizId, results)
		if err != nil {
			return nil, err
		}
	}

	return answer, nil
}

func (s *Service) submitReadingListening(ctx context.Context, userId string, quizId int, results models.QuizAnswer) (*models.Answer, map[int]models.QuestionSuccessCount, error) {
	results.Answer.UserCreated = userId
	results.Answer.DateCreated = time.Now()
	answer, err := s.answerRepo.Create(ctx, results.Answer)
	if err != nil {
		return nil, nil, err
	}

	quiz, err := s.quizRepo.GetQuizSubmitted(ctx, quizId)

	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, err
		} else {
			return nil, nil, common.ErrQuizNotFound
		}
	}

	quizCfg, err := s.quizRepo.List(ctx, models.QueryParams{}, func(tx *gorm.DB) {
		tx.Select("id", "type").Where("id", quizId).Where("status", common.QUIZ_STATUS_PUBLISHED)
	})
	if err != nil || len(quizCfg) == 0 {
		return nil, nil, err
	}

	var questionTypeSuccessCount = make(map[string]models.QuestionSuccessCount)
	var passageSuccessCount = make(map[int]models.QuestionSuccessCount)

	questionTypeSuccessCount, passageSuccessCount, err = s.countAnswerStatistic(ctx, quiz, results.QuestionResult)
	if err != nil {
		return nil, nil, err
	}

	var createsQuizLog []*models.SuccessQuizLog

	for i, v := range questionTypeSuccessCount {
		createsQuizLog = append(createsQuizLog, &models.SuccessQuizLog{
			Total:        v.Total,
			Success:      v.Success,
			Date:         time.Now(),
			Status:       1,
			Skill:        quizCfg[0].Type,
			QuestionType: i,
			UserId:       userId,
			Skipped:      v.Skip,
			Failed:       v.Failed,
			AnswerId:     answer.ID,
			QuizType:     results.Answer.QuizType,
		})
	}

	for i, v := range passageSuccessCount {
		createsQuizLog = append(createsQuizLog, &models.SuccessQuizLog{
			Total:    v.Total,
			Success:  v.Success,
			Date:     time.Now(),
			Status:   1,
			Skill:    quizCfg[0].Type,
			Passage:  i,
			UserId:   userId,
			Skipped:  v.Skip,
			Failed:   v.Failed,
			AnswerId: answer.ID,
			QuizType: results.Answer.QuizType,
		})
	}

	if len(createsQuizLog) == 0 {
		return answer, nil, err
	}

	err = s.successQuizLogRepo.CreatesMultiple(ctx, createsQuizLog)
	if err != nil {
		return answer, nil, err
	}
	return answer, passageSuccessCount, nil
}

func (s *Service) countAnswerStatistic(ctx context.Context, quiz *models.Quiz, correction []models.QuestionResult) (questionTypeSuccessCount map[string]models.QuestionSuccessCount, passageSuccessCount map[int]models.QuestionSuccessCount, err error) {
	for i, v := range correction {
		fmt.Println(i, v)
	}

	quizCfg, err := s.quizRepo.List(ctx, models.QueryParams{}, func(tx *gorm.DB) {
		tx.Select("id", "type").Where("id", quiz.ID).Where("status", common.QUIZ_STATUS_PUBLISHED)
	})
	if err != nil || len(quizCfg) == 0 {
		err = errors.New("fetch quiz error")
		return
	}

	resultObject := make(map[int]models.QuizResult)
	for _, r := range correction {
		resultObject[r.Id] = r.QuizResult
	}
	questionTypeSuccessCount = make(map[string]models.QuestionSuccessCount)
	passageSuccessCount = make(map[int]models.QuestionSuccessCount)

	for _, part := range quiz.Parts {
		for _, question := range part.Questions {
			subQuestionCount := question.CountTotalSubQuestion()
			if part.Passage != 0 {
				val, ok := passageSuccessCount[part.Passage]
				if !ok {
					val = models.QuestionSuccessCount{
						Total:   0,
						Success: 0,
					}
				}

				result, ok := resultObject[question.ID]
				if ok {
					val.Total += subQuestionCount
					val.Success += result.SuccessCount
					val.Failed += result.Total - result.SuccessCount
					val.Skip += subQuestionCount - result.Total
				} else {
					val.Skip += subQuestionCount
					val.Total += subQuestionCount
				}
				passageSuccessCount[part.Passage] = val
			}

			var logQuestionType string
			if question.QuestionType != "" {
				logQuestionType = question.QuestionType
			} else {
				logQuestionType = common.QUESTION_TYPE_CATEGORY_OTHERS
			}
			val, ok := questionTypeSuccessCount[logQuestionType]
			if !ok {
				val = models.QuestionSuccessCount{
					Total:   0,
					Success: 0,
				}
			}

			result, ok := resultObject[question.ID]
			if ok {
				val.Total += subQuestionCount
				val.Success += result.SuccessCount
				val.Failed += result.Total - result.SuccessCount
				val.Skip += subQuestionCount - result.Total
			} else {
				val.Skip += subQuestionCount
				val.Total += subQuestionCount
			}
			questionTypeSuccessCount[logQuestionType] = val
		}
	}

	return
}
