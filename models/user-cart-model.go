package models

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/guregu/null"
	"gorm.io/gorm"
)

type userCartDbFieldStruct struct {
	Id        string
	UserId    string
	ProductId string
	Quantity  string
	CreatedAt string
	CreatedBy string
	UpdatedAt string
	UpdatedBy string
	DeletedAt string
	DeletedBy string
}

var UserCartDbField = userCartDbFieldStruct{
	Id:        "id",
	UserId:    "user_id",
	ProductId: "product_id",
	Quantity:  "quantity",
	CreatedAt: "created_at",
	CreatedBy: "created_by",
	UpdatedAt: "updated_at",
	UpdatedBy: "updated_by",
	DeletedAt: "deleted_at",
	DeletedBy: "deleted_by",
}

type UserCart struct {
	Id        int `gorm:"primaryKey"`
	UserId    uuid.UUID
	ProductId int
	Quantity  int
	CreatedAt time.Time `gorm:"autoCreateTime"`
	CreatedBy string
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	UpdatedBy string
	DeletedAt gorm.DeletedAt `gorm:"index"`
	DeletedBy null.String
}

func (UserCart) TableName() string {
	return "nexmedis_user_carts"
}

type UserCartPrimaryId struct {
	Id int `gorm:"primaryKey"`
}

func (UserCartPrimaryId) TableName() string {
	return "nexmedis_user_carts"
}
