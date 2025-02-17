package repositories

import (
	"errors"

	"github.com/nexmedis-be-technical-test/models"
	"github.com/nexmedis-be-technical-test/utils/failure"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func (r *Repository) ProductFindById(primaryId models.ProductPrimaryId) (models.Product, error) {
	var product models.Product
	productData := r.PostgreSqlConn.Db.First(&product, primaryId.Id)
	err := productData.Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Product{}, failure.BadRequest("Product is not found")
		}
		log.Error().Err(err).Msg("[ProductFindById] Repository error retrieving product by id")
		return models.Product{}, err
	}

	return product, nil
}

func (r *Repository) ProductFindManyAndCountByFilter(filter models.Filter) ([]models.Product, int64, error) {
	var products []models.Product
	var totalProducts int64
	var err error

	// Start a transaction
	err = r.PostgreSqlConn.Db.Transaction(func(tx *gorm.DB) error {
		// Start building the query with the transaction context (tx)
		query := tx.Model(&models.Product{})
		countQuery := tx.Model(&models.Product{})

		// Handle select specific fields
		if len(filter.SelectFields) > 0 {
			query = query.Select(filter.SelectFields)
			countQuery = countQuery.Select(filter.SelectFields)
		}

		// Handle filter fields (where conditions)
		if len(filter.FilterFields) > 0 {
			for _, filterField := range filter.FilterFields {
				switch filterField.Operator {
				case models.OperatorEqual:
					query = query.Where(filterField.Field+" = ?", filterField.Value)
					countQuery = countQuery.Where(filterField.Field+" = ?", filterField.Value)
				case models.OperatorBetween:
					values, ok := filterField.Value.([]interface{})
					if ok && len(values) == 2 {
						query = query.Where(filterField.Field+" BETWEEN ? AND ?", values[0], values[1])
						countQuery = countQuery.Where(filterField.Field+" BETWEEN ? AND ?", values[0], values[1])
					}
				case models.OperatorIn:
					query = query.Where(filterField.Field+" IN ?", filterField.Value)
					countQuery = countQuery.Where(filterField.Field+" IN ?", filterField.Value)
				case models.OperatorIsNull:
					query = query.Where(filterField.Field + " IS NULL")
					countQuery = countQuery.Where(filterField.Field + " IS NULL")
				case models.OperatorNot:
					query = query.Where(filterField.Field+" != ?", filterField.Value)
					countQuery = countQuery.Where(filterField.Field+" != ?", filterField.Value)
				case models.OperatorContains:
					// Ensure that the filter value is a string
					strValue, ok := filterField.Value.(string)
					if ok {
						// ILIKE for case-insensitive matching in PostgreSQL
						query = query.Where(filterField.Field+" ILIKE ?", "%"+strValue+"%")
						countQuery = countQuery.Where(filterField.Field+" ILIKE ?", "%"+strValue+"%")
					}
				default:
					log.Warn().Msgf("[ProductFindManyAndCountByFilter] Unsupported filter operator: %s for field: %s", filterField.Operator, filterField.Field)
				}
			}
		}

		// Handle pagination
		if filter.Pagination.Page > 0 && filter.Pagination.PageSize > 0 {
			offset := (filter.Pagination.Page - 1) * filter.Pagination.PageSize
			query = query.Offset(offset).Limit(filter.Pagination.PageSize)
		}

		// Handle sorting
		if len(filter.Sorts) > 0 {
			for _, sort := range filter.Sorts {
				switch sort.Order {
				case models.SortAsc:
					query = query.Order(sort.Field + " asc")
				case models.SortDesc:
					query = query.Order(sort.Field + " desc")
				default:
					log.Warn().Msgf("[ProductFindManyAndCountByFilter] Unknown sort order: %s for field: %s", sort.Order, sort.Field)
				}
			}
		}

		// Finds all records matching given conditions
		err = query.Find(&products).Error
		if err != nil {
			log.Error().Err(err).Msg("[ProductFindManyAndCountByFilter] Repository error retrieving products by filter")
			return err
		}

		// Count total products based on the filtered conditions
		err = countQuery.Count(&totalProducts).Error
		if err != nil {
			log.Error().Err(err).Msg("[ProductFindManyAndCountByFilter] Repository error counting total products")
			return err
		}

		return nil
	})

	if err != nil {
		return nil, 0, err
	}

	return products, totalProducts, nil
}
