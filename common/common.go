package common

// Define New Table Name here
const (
	POSTGRES_TABLE_NAME_USERS               = "public.users"
	POSTGRES_TABLE_NAME_ROLES               = "public.roles"
	POSTGRES_TABLE_NAME_TARGETS             = "public.student_target"
	POSTGRES_TABLE_NAME_OTPS                = "public.otps"
	POSTGRES_TABLE_NAME_OTP_ATTEMPTS        = "public.otp_attempts"
	POSTGRES_TABLE_NAME_PART                = "public.part"
	POSTGRES_TABLE_NAME_QUESTION            = "public.question"
	POSTGRES_TABLE_NAME_TAG_SEARCH          = "public.tag_search"
	POSTGRES_TABLE_NAME_TAG_POSITION        = "public.tag_position"
	POSTGRES_TABLE_NAME_QUIZ_PART           = "public.quiz_part"
	POSTGRES_TABLE_NAME_QUIZ                = "public.quiz"
	POSTGRES_TABLE_NAME_QUIZ_SKILL          = "public.type"
	POSTGRES_TABLE_NAME_ANSWER              = "public.answers"
	POSTGRES_TABLE_NAME_SUCCESS_QUIZ_LOG    = "public.success_quiz_log"
	POSTGRES_TABLE_NAME_QUIZ_TAG_SEARCH     = "public.quiz_tag_search"
	POSTGRES_TABLE_NAME_USER_VOCAB_CATEGORY = "public.user_vocab_category"
	POSTGRES_TABLE_NAME_USER_VOCAB_BANK     = "public.user_vocab_bank"
	POSTGRES_TABLE_NAME_VOCAB               = "public.vocab"
)

const (
	ROLE_END_USER        = "END_USER"
	ROLE_END_USER_UUID   = "7b524019-4e1f-419d-bd15-30585f8c57ba"
	USER_PROVIDER_GOOGLE = "google"
	DEFAULT_AVATAR       = "https://mdapjazwsbewinkegonp.supabase.co/storage/v1/object/public/portal_attachments/default_avt.jpg?t=2024-12-24T08%3A36%3A53.223Z"
	DEFAULT_VOCAB_COUNT  = 10
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
