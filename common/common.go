package common

// Define New Table Name here
const (
	POSTGRES_TABLE_NAME_USERS        = "public.users"
	POSTGRES_TABLE_NAME_ROLES        = "public.roles"
	POSTGRES_TABLE_NAME_TARGETS      = "public.student_target"
	POSTGRES_TABLE_NAME_OTPS         = "public.otps"
	POSTGRES_TABLE_NAME_OTP_ATTEMPTS = "public.otp_attempts"
)

const (
	ROLE_END_USER        = "end_user"
	ROLE_END_USER_UUID   = "da0e07d4-ce51-4784-a5a9-a018434adf8e"
	USER_PROVIDER_GOOGLE = "google"
)

const (
	FromEmail         = "mainhatnam01@gmail.com"
	TypeResetPassword = "reset_password"
)
