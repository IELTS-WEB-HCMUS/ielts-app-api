package models

import (
	"ielts-web-api/common"
)

type User struct {
	ID                 string  `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Email              string  `json:"email" gorm:"uniqueIndex"`
	FirstName          *string `json:"first_name" gorm:"column:first_name"`
	LastName           *string `json:"last_name" gorm:"column:last_name"`
	Password           string  `json:"-" gorm:"password"`
	RoleID             string  `json:"role_id" gorm:"column:role"`
	IsBanned           bool    `json:"is_banned" gorm:"is_banned"`
	EmailNotifications bool    `json:"email_notifications" gorm:"email_notifications"`
	Provider           string  `gorm:"nullable"`
	VocabUsageCount    int64   `json:"vocab_usage_count" gorm:"vocab_usage_count"`
	Avatar             string  `json:"avatar" gorm:"avatar"`
}

func (User) TableName() string {
	// define common variable here
	return common.POSTGRES_TABLE_NAME_USERS
}

type GoogleUser struct {
	Iss           string `json:"iss"`
	Azp           string `json:"azp"`
	Aud           string `json:"aud"`
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	Exp           string `json:"exp"`
	AtHash        string `json:"at_hash"`
	Alg           string `json:"alg"`
	Kid           string `json:"kid"`
	Typ           string `json:"typ"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
}
