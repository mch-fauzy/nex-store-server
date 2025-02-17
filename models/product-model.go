package models

import (
	"time"

	"github.com/guregu/null"
	"gorm.io/gorm"
)

type productDbFieldStruct struct {
	Id          string
	Sku         string
	Name        string
	Slug        string
	Description string
	Color       string
	Size        string
	Price       string
	Stock       string
	CreatedAt   string
	CreatedBy   string
	UpdatedAt   string
	UpdatedBy   string
	DeletedAt   string
	DeletedBy   string
}

var ProductDbField = productDbFieldStruct{
	Id:          "id",
	Sku:         "sku",
	Name:        "name",
	Slug:        "slug",
	Description: "description",
	Color:       "color",
	Size:        "size",
	Price:       "price",
	Stock:       "stock",
	CreatedAt:   "created_at",
	CreatedBy:   "created_by",
	UpdatedAt:   "updated_at",
	UpdatedBy:   "updated_by",
	DeletedAt:   "deleted_at",
	DeletedBy:   "deleted_by",
}

type Product struct {
	Id          int    `gorm:"primaryKey"`
	Sku         string `gorm:"uniqueIndex"`
	Name        string `gorm:"uniqueIndex"`
	Slug        string `gorm:"uniqueIndex"`
	Description null.String
	Color       null.String
	Size        null.String
	Price       float32
	Stock       int
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	CreatedBy   string
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	UpdatedBy   string
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	DeletedBy   null.String
}

func (Product) TableName() string {
	return "nexmedis_products"
}

type ProductPrimaryId struct {
	Id int `gorm:"primaryKey"`
}

func (ProductPrimaryId) TableName() string {
	return "nexmedis_products"
}
