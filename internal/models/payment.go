package models

import (
	"ielts-web-api/common"
	"time"
)

type Payment struct {
	Amount          int       `json:"amout" gorm:"column:amount;not null"`
	Status          string    `json:"status" gorm:"column:status;not null"`
	Type            string    `json:"type" gorm:"column:type;not null"`
	UserId          string    `json:"user_id" gorm:"column:user_id;not null"`
	TransactionTime time.Time `json:"transaction_time" gorm:"column:transaction_time;not null"`
}

// TableName overrides the default table name for GORM
func (Payment) TableName() string {
	return common.POSTGRES_TABLE_NAME_PAYMENT
}
