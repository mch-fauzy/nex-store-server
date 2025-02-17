package repositories

import (
	"github.com/nexmedis-be-technical-test/models"
	"github.com/nexmedis-be-technical-test/utils/failure"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func (r *Repository) UserCartCreate(data *models.UserCart) error {
	createdUserCart := r.PostgreSqlConn.Db.Create(data)
	err := createdUserCart.Error
	if err != nil {
		log.Error().Err(err).Msg("[UserCartCreate] Repository error creating user cart")
		return err
	}

	return nil
}

func (r *Repository) UserCartUpdateById(primaryId models.UserCartPrimaryId, data *models.UserCart) error {
	var userCart models.UserCart
	updatedUserCart := r.PostgreSqlConn.Db.Model(&userCart).Where("id = ?", primaryId.Id).Updates(data)
	err := updatedUserCart.Error
	if err != nil {
		log.Error().Err(err).Msg("[UserCartUpdateById] Repository error updating user cart by id")
		return err
	}

	if updatedUserCart.RowsAffected == 0 {
		err = failure.NotFound("User cart not found")
		return err
	}

	return nil
}

func (r *Repository) UserCartFindManyAndCountByFilter(filter models.Filter) ([]models.UserCart, int64, error) {
	var userCarts []models.UserCart
	var totalUserCarts int64
	var err error

	// Start a transaction
	err = r.PostgreSqlConn.Db.Transaction(func(tx *gorm.DB) error {
		// Start building the query with the transaction context (tx)
		query := tx.Model(&models.UserCart{})
		countQuery := tx.Model(&models.UserCart{})

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
					log.Warn().Msgf("[UserCartFindManyAndCountByFilter] Unsupported filter operator: %s for field: %s", filterField.Operator, filterField.Field)
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
					log.Warn().Msgf("[UserCartFindManyAndCountByFilter] Unknown sort order: %s for field: %s", sort.Order, sort.Field)
				}
			}
		}

		// Finds all records matching given conditions
		err = query.Find(&userCarts).Error
		if err != nil {
			log.Error().Err(err).Msg("[UserCartFindManyAndCountByFilter] Repository error retrieving user carts by filter")
			return err
		}

		// Count total user carts based on the filtered conditions
		err = countQuery.Count(&totalUserCarts).Error
		if err != nil {
			log.Error().Err(err).Msg("[UserFindManyAndCountByFilter] Repository error counting total user carts")
			return err
		}

		return nil
	})

	if err != nil {
		return nil, 0, err
	}

	return userCarts, totalUserCarts, nil
}
