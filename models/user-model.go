package models

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/guregu/null"
	"gorm.io/gorm"
)

type userDbFieldStruct struct {
	Id        string
	RoleId    string
	Email     string
	Password  string
	LastLogin string
	Balance   string
	CreatedAt string
	CreatedBy string
	UpdatedAt string
	UpdatedBy string
	DeletedAt string
	DeletedBy string
}

var UserDbField = userDbFieldStruct{
	Id:        "id",
	RoleId:    "role_id",
	Email:     "email",
	Password:  "password",
	LastLogin: "last_login",
	Balance:   "balance",
	CreatedAt: "created_at",
	CreatedBy: "created_by",
	UpdatedAt: "updated_at",
	UpdatedBy: "updated_by",
	DeletedAt: "deleted_at",
	DeletedBy: "deleted_by",
}

type User struct {
	Id        uuid.UUID `gorm:"primaryKey"`
	RoleId    int
	Email     string `gorm:"uniqueIndex"`
	Password  string
	LastLogin *time.Time
	Balance   float32
	CreatedAt time.Time `gorm:"autoCreateTime"`
	CreatedBy string
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	UpdatedBy string
	DeletedAt gorm.DeletedAt `gorm:"index"`
	DeletedBy null.String
}

/* Set table name for User struct*/
func (User) TableName() string {
	return "nexmedis_users"
}

type UserPrimaryId struct {
	Id uuid.UUID `gorm:"primaryKey"`
}

func (UserPrimaryId) TableName() string {
	return "nexmedis_users"
}
