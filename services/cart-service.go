package services

import (
	"github.com/gofrs/uuid"
	"github.com/nexmedis-be-technical-test/models"
	"github.com/nexmedis-be-technical-test/models/dto"
	"github.com/rs/zerolog/log"
)

func (s *Service) UserCartAddItem(req dto.UserCartAddItemRequest) (string, error) {
	message := "Failed"

	_, err := s.Repository.UserFindById(models.UserPrimaryId{Id: uuid.FromStringOrNil(req.UserId)})
	if err != nil {
		log.Error().Err(err).Msg("[UserCartAddItem] Service error retrieving user by id")
		return message, err
	}

	_, err = s.Repository.ProductFindById(models.ProductPrimaryId{Id: req.ProductId})
	if err != nil {
		log.Error().Err(err).Msg("[UserCartAddItem] Service error retrieving product by id")
		return message, err
	}

	userCarts, totalUserCarts, err := s.Repository.UserCartFindManyAndCountByFilter(models.Filter{
		FilterFields: []models.FilterField{
			{
				Field:    models.UserCartDbField.UserId,
				Operator: models.OperatorEqual,
				Value:    req.UserId,
			},
			{
				Field:    models.UserCartDbField.ProductId,
				Operator: models.OperatorEqual,
				Value:    req.ProductId,
			},
			{
				Field:    models.UserCartDbField.DeletedAt,
				Operator: models.OperatorIsNull,
				Value:    true,
			},
		},
		Sorts: []models.Sort{{
			Field: models.UserCartDbField.UpdatedAt,
			Order: models.SortDesc,
		}},
	})
	if err != nil {
		log.Error().Err(err).Msg("[UserCartAddItem] Service error retrieving user carts")
		return message, err
	}

	/* If a cart entry with given product id already exists, update its quantity */
	if totalUserCarts > 0 {
		updatedUserCart := req.UpdateCartQuantity(userCarts[0])

		err := s.Repository.UserCartUpdateById(models.UserCartPrimaryId{Id: userCarts[0].Id}, &updatedUserCart)
		if err != nil {
			log.Error().Err(err).Msg("[UserCartAddItem] Service error updating user cart")
			return message, err
		}

		message = "Success"
		return message, nil
	}

	userCart := req.ToModel()
	err = s.Repository.UserCartCreate(&userCart)
	if err != nil {
		log.Error().Err(err).Msg("[UserCartAddItem] Service error creating user cart")
		return message, err
	}

	message = "Success"
	return message, nil
}

func (s *Service) UserCartGetList(req dto.UserCartGetListRequest) ([]dto.UserCartGetListResponse, error) {
	var responses []dto.UserCartGetListResponse

	_, err := s.Repository.UserFindById(models.UserPrimaryId{Id: uuid.FromStringOrNil(req.UserId)})
	if err != nil {
		log.Error().Err(err).Msg("[UserCartGetList] Service error retrieving user by id")
		return responses, err
	}

	userCarts, _, err := s.Repository.UserCartFindManyAndCountByFilter(models.Filter{
		FilterFields: []models.FilterField{
			{
				Field:    models.UserCartDbField.UserId,
				Operator: models.OperatorEqual,
				Value:    req.UserId,
			},
			{
				Field:    models.UserCartDbField.DeletedAt,
				Operator: models.OperatorIsNull,
				Value:    true,
			},
		},
		Sorts: []models.Sort{
			{
				Field: models.UserCartDbField.UpdatedAt,
				Order: models.SortDesc,
			},
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("[UserCartGetList] Service error getting user carts")
		return responses, err
	}

	var productIDs []int
	for _, cart := range userCarts {
		productIDs = append(productIDs, cart.ProductId)
	}

	products, _, err := s.Repository.ProductFindManyAndCountByFilter(models.Filter{
		FilterFields: []models.FilterField{
			{
				Field:    models.ProductDbField.Id,
				Operator: models.OperatorIn,
				Value:    productIDs,
			},
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("[UserCartGetList] Service error getting products")
		return responses, err
	}

	/* Build a mapping of productId */
	productMap := make(map[int]models.Product)
	for _, product := range products {
		productMap[product.Id] = product
	}

	for _, cart := range userCarts {
		product, ok := productMap[cart.ProductId]
		if !ok {
			continue
		}

		response := dto.UserCartGetListResponse{
			ProductName: product.Name,
			Quantity:    cart.Quantity,
			TotalPrice:  float32(cart.Quantity) * product.Price,
		}
		responses = append(responses, response)
	}

	return responses, nil
}
