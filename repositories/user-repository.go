package repositories

import (
	"errors"

	"github.com/nexmedis-be-technical-test/models"
	"github.com/nexmedis-be-technical-test/utils/failure"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func (r *Repository) UserCreate(data *models.User) error {
	createdUser := r.PostgreSqlConn.Db.Create(data)
	err := createdUser.Error
	if err != nil {
		log.Error().Err(err).Msg("[UserCreate] Repository error creating user")
		return err
	}

	return nil
}

func (r *Repository) UserUpdateById(primaryId models.UserPrimaryId, data *models.User) error {
	var user models.User
	updatedUser := r.PostgreSqlConn.Db.Model(&user).Where("id = ?", primaryId.Id).Updates(data)
	err := updatedUser.Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return failure.BadRequest("User is not found")
		}
		log.Error().Err(err).Msg("[UserUpdateById] Repository error updating user by id")
		return err
	}

	if updatedUser.RowsAffected == 0 {
		err = failure.NotFound("User not found")
		return err
	}

	return nil
}

func (r *Repository) UserFindById(primaryId models.UserPrimaryId) (models.User, error) {
	var user models.User
	userData := r.PostgreSqlConn.Db.First(&user, primaryId.Id)
	err := userData.Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, failure.BadRequest("User not found")
		}
		log.Error().Err(err).Msg("[UserFindById] Repository error retrieving user by id")
		return user, err
	}

	return user, nil
}

func (r *Repository) UserFindManyAndCountByFilter(filter models.Filter) ([]models.User, int64, error) {
	var users []models.User
	var totalUsers int64
	var err error

	// Start a transaction
	err = r.PostgreSqlConn.Db.Transaction(func(tx *gorm.DB) error {
		// Start building the query with the transaction context (tx)
		query := tx.Model(&models.User{})
		countQuery := tx.Model(&models.User{})

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
					log.Warn().Msgf("[UserFindManyAndCountByFilter] Unsupported filter operator: %s for field: %s", filterField.Operator, filterField.Field)
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
					log.Warn().Msgf("[UserFindManyAndCountByFilter] Unknown sort order: %s for field: %s", sort.Order, sort.Field)
				}
			}
		}

		// Finds all records matching given conditions
		err = query.Find(&users).Error
		if err != nil {
			log.Error().Err(err).Msg("[UserFindManyAndCountByFilter] Repository error retrieving users by filter")
			return err
		}

		// Count total users based on the filtered conditions
		err = countQuery.Count(&totalUsers).Error
		if err != nil {
			log.Error().Err(err).Msg("[UserFindManyAndCountByFilter] Repository error counting total users")
			return err
		}

		return nil
	})

	if err != nil {
		return nil, 0, err
	}

	return users, totalUsers, nil
}
