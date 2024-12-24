package services

import (
	"ielts-web-api/internal/repositories"
)

type Service struct {
	userRepo       *repositories.UserRepository
	targetRepo     *repositories.TargetRepository
	otpRepo        *repositories.OTPRepository
	otpAttemptRepo *repositories.OTPAttemptRepository

	quizRepo              *repositories.QuizRepo
	quizSkillRepo         *repositories.QuizSkillRepo
	tagSearchRepo         *repositories.TagSearchRepository
	tagSearchPositionRepo *repositories.TagSearchPositionRepo
	answerRepo            *repositories.AnswerRepo
	successQuizLogRepo    *repositories.SuccessQuizLogRepo

	vocabCategoriesRepo *repositories.UserVocabCategoryRepository
	userVocabBankRepo   *repositories.UserVocabBankRepository
}

func NewService(repos ...interface{}) *Service {
	service := &Service{}
	for _, repo := range repos {
		switch repo.(type) {
		case *repositories.UserRepository:
			service.userRepo = repo.(*repositories.UserRepository)
		case *repositories.TargetRepository:
			service.targetRepo = repo.(*repositories.TargetRepository)
		case *repositories.OTPRepository:
			service.otpRepo = repo.(*repositories.OTPRepository)
		case *repositories.OTPAttemptRepository:
			service.otpAttemptRepo = repo.(*repositories.OTPAttemptRepository)
		case *repositories.QuizRepo:
			service.quizRepo = repo.(*repositories.QuizRepo)
		case *repositories.QuizSkillRepo:
			service.quizSkillRepo = repo.(*repositories.QuizSkillRepo)
		case *repositories.TagSearchRepository:
			service.tagSearchRepo = repo.(*repositories.TagSearchRepository)
		case *repositories.TagSearchPositionRepo:
			service.tagSearchPositionRepo = repo.(*repositories.TagSearchPositionRepo)
		case *repositories.AnswerRepo:
			service.answerRepo = repo.(*repositories.AnswerRepo)
		case *repositories.SuccessQuizLogRepo:
			service.successQuizLogRepo = repo.(*repositories.SuccessQuizLogRepo)
		case *repositories.UserVocabCategoryRepository:
			service.vocabCategoriesRepo = repo.(*repositories.UserVocabCategoryRepository)
		case *repositories.UserVocabBankRepository:
			service.userVocabBankRepo = repo.(*repositories.UserVocabBankRepository)
		default:
			panic("Unknown repository type provided")
		}
	}
	return service
}
