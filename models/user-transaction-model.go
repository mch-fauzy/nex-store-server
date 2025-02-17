package models

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/guregu/null"
	"gorm.io/gorm"
)

type UserTransaction struct {
	Id                  int `gorm:"primaryKey"`
	UserId              uuid.UUID
	TransactionStatusId int
	TotalAmount         float32
	InvoiceNumber       string
	CreatedAt           time.Time `gorm:"autoCreateTime"`
	CreatedBy           string
	UpdatedAt           time.Time `gorm:"autoUpdateTime"`
	UpdatedBy           string
	DeletedAt           gorm.DeletedAt `gorm:"index"`
	DeletedBy           null.String
}

func (UserTransaction) TableName() string {
	return "nexmedis_user_transactions"
}
