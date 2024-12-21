package common

// Define New Table Name here
const (
	POSTGRES_TABLE_NAME_USERS            = "public.users"
	POSTGRES_TABLE_NAME_ROLES            = "public.roles"
	POSTGRES_TABLE_NAME_TARGETS          = "public.student_target"
	POSTGRES_TABLE_NAME_OTPS             = "public.otps"
	POSTGRES_TABLE_NAME_OTP_ATTEMPTS     = "public.otp_attempts"
	POSTGRES_TABLE_NAME_PART             = "public.part"
	POSTGRES_TABLE_NAME_QUESTION         = "public.question"
	POSTGRES_TABLE_NAME_TAG_SEARCH       = "public.tag_search"
	POSTGRES_TABLE_NAME_TAG_POSITION     = "public.tag_position"
	POSTGRES_TABLE_NAME_QUIZ_PART        = "public.quiz_part"
	POSTGRES_TABLE_NAME_QUIZ             = "public.quiz"
	POSTGRES_TABLE_NAME_QUIZ_SKILL       = "public.type"
	POSTGRES_TABLE_NAME_ANSWER           = "public.answers"
	POSTGRES_TABLE_NAME_SUCCESS_QUIZ_LOG = "public.success_quiz_log"
	POSTGRES_TABLE_NAME_QUIZ_TAG_SEARCH  = "public.quiz_tag_search"
)

const (
	ROLE_END_USER        = "end_user"
	ROLE_END_USER_UUID   = "da0e07d4-ce51-4784-a5a9-a018434adf8e"
	USER_PROVIDER_GOOGLE = "google"
)

const (
	FromEmail = "mainhatnam01@gmail.com"
)

const (
	RESET_PASSSWORD_TYPE = "reset_password"
	VERIFY_EMAIL_TYPE    = "verify_email"
)

// Define other common variable here
const (
	QUESTION_TYPE_SINGLE_RADIO      = "SINGLE-RADIO"
	QUESTION_TYPE_SINGLE_SELECTION  = "SINGLE-SELECTION"
	QUESTION_TYPE_FILL_IN_THE_BLANK = "FILL-IN-THE-BLANK"
	QUESTION_TYPE_MULTIPLE          = "MULTIPLE"
	QUESTION_TYPE_CATEGORY_OTHERS   = "MULTIPLE"
)

const (
	QuizSkillReading      = "READING"
	QuizSkillListening    = "LISTENING"
	QuizSkillSpeaking     = "SPEAKING"
	QuizSkillWriting      = "WRITING"
	QUIZ_STATUS_PUBLISHED = "published"
)

const (
	AnswerStatisticByQuiz       = 1
	AnswerStatisticQuestionType = 2
	AnswerStatisticByPassage    = 3

	QuizSubmittedStatusUnknown = 0
	QuizSubmittedStatusYes     = 1
	QuizSubmittedStatusNo      = 2
	QuizTypeUnknown            = 0
	QuizTypeExercise           = 1
	QuizTypeAssignment         = 2
	QuizTypeTest               = 3
)
