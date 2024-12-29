package services

import (
	"context"
	"errors"
	"fmt"
	"ielts-web-api/common"
	"ielts-web-api/internal/models"
	"ielts-web-api/internal/repositories"
	"log"
	"sync"
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
		request.TagQuestionType != nil {

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

	if request.Mode != nil {
		filters = append(filters, func(tx *gorm.DB) {
			tx.Where("quiz.mode = ?", *request.Mode)
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
			Selected: []string{
				"quiz.id",
				"quiz.status",
				"quiz.type",
				"quiz.title",
				"quiz.time",
				"quiz.listening_file AS listening",
				"quiz.level",
				"quiz.vote_count",
				"quiz.description",
				"quiz.thumbnail",
				"quiz.mode",
				"quiz.is_public",
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
		fmt.Println("record.IsSubmitted: ", record.IsSubmitted)
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
			return db.Joins("INNER JOIN quiz_part ON quiz_part.quiz_id = ? AND quiz_part.part_id = part.quiz", req.QuizID).Order("quiz_part.sort")
		}).Preload("Parts.Questions", func(db *gorm.DB) *gorm.DB {
			return db.Order("question.order")
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
	// Fetch quiz details with timeout
	fmt.Println("SubmitQuizAnswer")

	// Get quiz and validate
	q, err := s.quizRepo.GetByID(ctx, quizId)
	if err != nil {
		log.Printf("Failed to fetch quiz by ID: %v, Error: %v", quizId, err)
		return nil, err
	}
	fmt.Println("SubmitQuizAnswer - Fetched Quiz:", q)

	// Get quiz skill details
	t, err := s.quizSkillRepo.GetByID(ctx, q.Type)
	if err != nil {
		log.Printf("Failed to fetch quiz skill by type: %v, Error: %v", q.Type, err)
		return nil, err
	}

	fmt.Println("SubmitQuizAnswer - Fetched Quiz Skill:", t)

	// Process based on quiz type
	if t.ID == common.QuizSkillReading {
		answer, _, err = s.submitReadingListening(ctx, userId, quizId, results)
		if err != nil {
			return nil, err
		}
	}

	return answer, nil
}

func (s *Service) submitReadingListening(ctx context.Context, userId string, quizId int, results models.QuizAnswer) (*models.Answer, map[int]models.QuestionSuccessCount, error) {
	// Save initial answer data
	results.Answer.UserCreated = userId
	results.Answer.DateCreated = time.Now()

	log.Printf("[DEBUG] Submitting Quiz Answer - Quiz ID: %d, User ID: %s", quizId, userId)

	// Insert answer into DB
	answer, err := s.answerRepo.Create(ctx, results.Answer)
	if err != nil {
		log.Printf("[ERROR] Failed to create answer entry: %v", err)
		return nil, nil, err
	}
	log.Printf("[DEBUG] Answer Created - Answer ID: %d", answer.ID)

	// Fetch submitted quiz with timeout

	quiz, err := s.quizRepo.GetQuizSubmitted(ctx, quizId)
	if err != nil {
		log.Printf("[ERROR] Failed to fetch quiz submitted: %v", err)
		return nil, nil, err
	}

	// Fetch quiz configuration
	quizCfg, err := s.quizRepo.List(ctx, models.QueryParams{}, func(tx *gorm.DB) {
		tx.Select("id", "type").Where("id", quizId).Where("status", common.QUIZ_STATUS_PUBLISHED)
	})
	if err != nil || len(quizCfg) == 0 {
		log.Printf("[ERROR] Failed to fetch quiz configuration for Quiz ID: %d", quizId)
		return nil, nil, errors.New("failed to fetch quiz configuration")
	}
	log.Printf("[DEBUG] Quiz Configuration Loaded - Type: %d", quizCfg[0].Type)

	// Calculate statistics
	questionStats, passageStats, err := s.countAnswerStatistic(ctx, quiz, results.QuestionResult)
	if err != nil {
		log.Printf("[ERROR] Failed to calculate statistics: %v", err)
		return nil, nil, err
	}

	// Prepare logs for batch insertion
	var logs []*models.SuccessQuizLog
	for qType, stat := range questionStats {
		log.Printf("[DEBUG] Adding Question Stat - Type: %s, Total: %d, Success: %d", qType, stat.Total, stat.Success)
		logs = append(logs, &models.SuccessQuizLog{
			Total:        stat.Total,
			Success:      stat.Success,
			Skill:        quizCfg[0].Type,
			QuestionType: qType,
			UserId:       userId,
			Skipped:      stat.Skip,
			Failed:       stat.Failed,
			AnswerId:     answer.ID,
		})
	}
	for passageID, stat := range passageStats {
		log.Printf("[DEBUG] Adding Passage Stat - Passage ID: %d, Total: %d, Success: %d", passageID, stat.Total, stat.Success)
		logs = append(logs, &models.SuccessQuizLog{
			Total:    stat.Total,
			Success:  stat.Success,
			Skill:    quizCfg[0].Type,
			Passage:  passageID,
			UserId:   userId,
			Skipped:  stat.Skip,
			Failed:   stat.Failed,
			AnswerId: answer.ID,
		})
	}

	// Check for empty logs
	if len(logs) == 0 {
		log.Printf("[ERROR] createsQuizLog is empty! Quiz ID: %d, Answer ID: %d, User ID: %s", quizId, answer.ID, userId)
		for i, v := range questionStats {
			log.Printf("[DEBUG] QuestionTypeStat[%s]: %+v", i, v)
		}
		for i, v := range passageStats {
			log.Printf("[DEBUG] PassageStat[%d]: %+v", i, v)
		}
		return answer, nil, errors.New("no logs to insert (createsQuizLog is empty)")
	}

	// Batch insert logs
	if err := s.successQuizLogRepo.CreatesMultiple(ctx, logs); err != nil {
		log.Printf("[ERROR] Failed to insert logs: %v", err)
		return answer, nil, err
	}

	log.Printf("[DEBUG] Logs successfully inserted. Logs Count: %d", len(logs))
	return answer, passageStats, nil
}

func (s *Service) countAnswerStatistic(ctx context.Context, quiz *models.Quiz, corrections []models.QuestionResult) (map[string]models.QuestionSuccessCount, map[int]models.QuestionSuccessCount, error) {
	var wg sync.WaitGroup
	questionStats := sync.Map{}
	passageStats := sync.Map{}

	log.Printf("[DEBUG] Counting Answer Statistics - Quiz ID: %d, Corrections Count: %d", quiz.ID, len(corrections))

	// Prepare result map
	resultMap := make(map[int]models.QuizResult)
	for _, correction := range corrections {
		resultMap[correction.Id] = correction.QuizResult
	}

	//Process quiz parts concurrently
	for _, part := range quiz.Parts {
		log.Printf("[DEBUG] Processing Part ID: %d, Questions Count: %d", part.ID, len(part.Questions))
		wg.Add(1)
		go func(part models.Part) {
			defer wg.Done()
			for _, question := range part.Questions {
				subQuestionCount := question.CountTotalSubQuestion()
				processPassageStats(&passageStats, part, resultMap, question, subQuestionCount)
				processQuestionStats(&questionStats, resultMap, question, subQuestionCount)
			}
		}(*part)
	}

	wg.Wait()

	// Convert sync.Map to regular maps
	finalQuestionStats := make(map[string]models.QuestionSuccessCount)
	finalPassageStats := make(map[int]models.QuestionSuccessCount)

	questionStats.Range(func(key, value interface{}) bool {
		finalQuestionStats[key.(string)] = value.(models.QuestionSuccessCount)
		return true
	})
	passageStats.Range(func(key, value interface{}) bool {
		finalPassageStats[key.(int)] = value.(models.QuestionSuccessCount)
		return true
	})

	log.Printf("[DEBUG] Statistics Calculated - Questions: %d, Passages: %d", len(finalQuestionStats), len(finalPassageStats))
	return finalQuestionStats, finalPassageStats, nil
}

func processPassageStats(passageStats *sync.Map, part models.Part, resultMap map[int]models.QuizResult, question models.Question, subQuestionCount int) {
	if part.Passage == 0 {
		return
	}
	val, _ := passageStats.LoadOrStore(part.Passage, models.QuestionSuccessCount{})
	stats := val.(models.QuestionSuccessCount)

	result, ok := resultMap[question.ID]
	if ok {
		stats.Total += subQuestionCount
		stats.Success += result.SuccessCount
		stats.Failed += result.Total - result.SuccessCount
		stats.Skip += subQuestionCount - result.Total
	} else {
		stats.Skip += subQuestionCount
		stats.Total += subQuestionCount
	}
	passageStats.Store(part.Passage, stats)

	log.Printf("[DEBUG] Updated Passage Stat - Passage ID: %d, Total: %d, Success: %d", part.Passage, stats.Total, stats.Success)
}

func processQuestionStats(questionStats *sync.Map, resultMap map[int]models.QuizResult, question models.Question, subQuestionCount int) {
	qType := question.QuestionType
	if qType == "" {
		qType = common.QUESTION_TYPE_CATEGORY_OTHERS
	}
	log.Printf("[DEBUG] Processing Question - ID: %d, Type: %s, SubQuestions: %d", question.ID, qType, subQuestionCount)
	val, _ := questionStats.LoadOrStore(qType, models.QuestionSuccessCount{})
	stats := val.(models.QuestionSuccessCount)

	result, ok := resultMap[question.ID]
	if ok {
		stats.Total += subQuestionCount
		stats.Success += result.SuccessCount
		stats.Failed += result.Total - result.SuccessCount
		stats.Skip += subQuestionCount - result.Total
		log.Printf("[DEBUG] Question Result Found - ID: %d, Total: %d, Success: %d, Failed: %d, Skipped: %d",
			question.ID, stats.Total, stats.Success, stats.Failed, stats.Skip)
	} else {
		stats.Skip += subQuestionCount
		stats.Total += subQuestionCount
		log.Printf("[WARNING] No Result Found for Question - ID: %d. Marked All as Skipped. Total: %d, Skipped: %d",
			question.ID, stats.Total, stats.Skip)
	}
	log.Printf("[DEBUG] Updated Question Stats - Type: %s, Total: %d, Success: %d, Failed: %d, Skipped: %d",
		qType, stats.Total, stats.Success, stats.Failed, stats.Skip)
	questionStats.Store(qType, stats)
}
