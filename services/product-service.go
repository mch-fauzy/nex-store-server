package services

import (
	"github.com/nexmedis-be-technical-test/models"
	"github.com/nexmedis-be-technical-test/models/dto"
	"github.com/nexmedis-be-technical-test/utils/pagination"
	"github.com/rs/zerolog/log"
)

func (s *Service) ProductGetListByFilter(req dto.ProductGetByFilterRequest) ([]dto.ProductGetByFilterResponse, dto.PaginationResponse, error) {
	var responses []dto.ProductGetByFilterResponse
	products, totalProducts, err := s.Repository.ProductFindManyAndCountByFilter(models.Filter{
		Pagination: models.Pagination{
			Page:     req.ParsedPage,
			PageSize: req.ParsedPageSize,
		},
		FilterFields: []models.FilterField{
			{
				Field:    models.ProductDbField.Name,
				Operator: models.OperatorContains,
				Value:    req.Name,
			},
			{
				Field:    models.ProductDbField.DeletedAt,
				Operator: models.OperatorIsNull,
				Value:    true,
			},
		},
		Sorts: []models.Sort{
			{
				Field: models.ProductDbField.CreatedAt,
				Order: models.SortDesc,
			},
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("[ProductGetListByFilter] Service error getting products")
		return responses, dto.PaginationResponse{}, err
	}

	responses = dto.ProductBuildGetByFilterResponse(products)
	metadata := pagination.CalculatePaginationMetadata(int64(req.ParsedPage), int64(req.ParsedPageSize), totalProducts)
	return responses, metadata, nil
}
