package dto

import (
	"strconv"

	"github.com/guregu/null"
	"github.com/nexmedis-be-technical-test/models"
	"github.com/nexmedis-be-technical-test/utils/constant"
	"github.com/nexmedis-be-technical-test/utils/failure"
)

type ProductGetByFilterRequest struct {
	Page           string `json:"page"`
	PageSize       string `json:"page_size"`
	Name           string `json:"name"`
	ParsedPage     int
	ParsedPageSize int
}

func (r *ProductGetByFilterRequest) Validate() error {
	if r.Page == "" {
		r.ParsedPage = constant.DefaultPage
	} else {
		page, err := strconv.Atoi(r.Page)
		if err != nil {
			return failure.BadRequest("Page must be numeric")
		}

		if page <= 0 {
			page = constant.DefaultPage
		}
		r.ParsedPage = page
	}

	if r.PageSize == "" {
		r.ParsedPageSize = constant.DefaultPageSize
	} else {
		pageSize, err := strconv.Atoi(r.PageSize)
		if err != nil {
			return failure.BadRequest("Page size must be numeric")
		}

		if pageSize <= 0 {
			pageSize = constant.DefaultPageSize
		}
		r.ParsedPageSize = pageSize
	}

	return nil
}

type ProductGetByFilterResponse struct {
	Sku         string      `json:"sku"`
	Name        string      `json:"name"`
	Slug        string      `json:"slug"`
	Description null.String `json:"description"`
	Color       null.String `json:"color"`
	Size        null.String `json:"size"`
	Price       float32     `json:"price"`
	Stock       int         `json:"stock"`
	CreatedAt   string      `json:"createdAt"`
	UpdatedAt   string      `json:"updatedAt"`
}

func ProductBuildGetByFilterResponse(products []models.Product) []ProductGetByFilterResponse {
	var responses []ProductGetByFilterResponse
	for _, product := range products {
		responses = append(responses, ProductGetByFilterResponse{
			Sku:         product.Sku,
			Name:        product.Name,
			Slug:        product.Slug,
			Description: product.Description,
			Color:       product.Color,
			Size:        product.Size,
			Price:       product.Price,
			Stock:       product.Stock,
			CreatedAt:   product.CreatedAt.Format(constant.DateTimeUTCFormat),
			UpdatedAt:   product.UpdatedAt.Format(constant.DateTimeUTCFormat),
		})
	}
	return responses
}
