package models

import (
	"ielts-web-api/common"
	"time"
)

type SignupRequest struct {
	Email     string `json:"email" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Role      string `json:"role" binding:"required"`
}

type LoginRequest struct {
	Email    *string `json:"email,omitempty"`
	Password *string `json:"password,omitempty"`
	IdToken  *string `json:"id_token,omitempty"`
}

type OTP struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Target     string    `gorm:"size:255;not null" json:"target"`
	Type       string    `gorm:"size:50;not null" json:"type"`
	OTPCode    string    `gorm:"size:6;not null" json:"otp_code"`
	ExpiredAt  time.Time `gorm:"not null" json:"expired_at"`
	IsVerified bool      `gorm:"default:false" json:"is_verified"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (OTP) TableName() string {
	return common.POSTGRES_TABLE_NAME_OTPS
}

type ResetPasswordRequest struct {
	Email       string `json:"email" binding:"required,email"`
	NewPassword string `json:"new_password" binding:"required"`
}

type OTPAttempt struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	OTPID     uint      `gorm:"not null" json:"otp_id"`
	Value     string    `gorm:"size:6;not null" json:"value"`
	IsSuccess bool      `gorm:"default:false" json:"is_success"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (OTPAttempt) TableName() string {
	return common.POSTGRES_TABLE_NAME_OTP_ATTEMPTS
}

type OTPValidateRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required"`
}
