package services

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/guregu/null"
	"github.com/nexmedis-be-technical-test/models"
	"github.com/nexmedis-be-technical-test/models/dto"
	"github.com/nexmedis-be-technical-test/utils/failure"
	"github.com/nexmedis-be-technical-test/utils/invoice"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func (s *Service) TransactionTopUpBalanceByUserId(req dto.TransactionTopUpBalanceByUserIdRequest) (dto.TransactionTopUpBalanceByUserIdResponse, error) {
	var response dto.TransactionTopUpBalanceByUserIdResponse
	users, totalUsers, err := s.Repository.UserFindManyAndCountByFilter(models.Filter{
		FilterFields: []models.FilterField{
			{
				Field:    models.UserDbField.Id,
				Operator: models.OperatorEqual,
				Value:    req.UserId,
			},
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("[TransactionTopUpBalanceByUserId] Service error getting user")
		return response, err
	}

	if totalUsers == 0 {
		return response, failure.BadRequest("User not found")
	}

	err = s.Repository.UserUpdateById(models.UserPrimaryId{Id: users[0].Id}, &models.User{
		Balance:   users[0].Balance + float32(req.TopupAmount),
		UpdatedBy: req.Email,
	})
	if err != nil {
		log.Error().Err(err).Msg("[TransactionTopUpBalanceByUserId] Service error updating user balance")
		return response, err
	}

	user, err := s.Repository.UserFindById(models.UserPrimaryId{Id: users[0].Id})
	if err != nil {
		log.Error().Err(err).Msg("[TransactionTopUpBalanceByUserId] Service error retrieving user by id")
		return response, err
	}

	response = dto.TransactionTopUpBalanceByUserIdResponse{
		Balance: user.Balance,
	}

	return response, nil
}

func (s *Service) TransactionWithdrawBalanceByUserId(req dto.TransactionWithdrawBalanceByUserIdRequest) (dto.TransactionWithdrawBalanceByUserIdResponse, error) {
	var response dto.TransactionWithdrawBalanceByUserIdResponse

	users, totalUsers, err := s.Repository.UserFindManyAndCountByFilter(models.Filter{
		FilterFields: []models.FilterField{
			{
				Field:    models.UserDbField.Id,
				Operator: models.OperatorEqual,
				Value:    req.UserId,
			},
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("[TransactionWithdrawBalanceByUserId] Service error getting user")
		return response, err
	}

	if totalUsers == 0 {
		return response, failure.BadRequest("User not found")
	}

	currentBalance := users[0].Balance
	withdrawAmount := float32(req.WithdrawAmount)
	if currentBalance < withdrawAmount {
		err = failure.BadRequest("Insufficient balance")
		return response, err
	}

	err = s.Repository.UserUpdateById(models.UserPrimaryId{Id: users[0].Id}, &models.User{
		Balance:   currentBalance - withdrawAmount,
		UpdatedBy: req.Email,
	})
	if err != nil {
		log.Error().Err(err).Msg("[TransactionWithdrawBalanceByUserId] Service error updating user balance")
		return response, err
	}

	user, err := s.Repository.UserFindById(models.UserPrimaryId{Id: users[0].Id})
	if err != nil {
		log.Error().Err(err).Msg("[TransactionWithdrawBalanceByUserId] Service error retrieving user by id")
		return response, err
	}

	response = dto.TransactionWithdrawBalanceByUserIdResponse{
		Balance: user.Balance,
	}

	return response, nil
}

func (s *Service) TransactionPurchaseCart(req dto.TransactionPurchaseCartRequest) (dto.TransactionPurchaseCartResponse, error) {
	var response dto.TransactionPurchaseCartResponse

	userCarts, totalUserCarts, err := s.Repository.UserCartFindManyAndCountByFilter(models.Filter{
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
	})
	if err != nil {
		log.Error().Err(err).Msg("[TransactionPurchaseCart] Service error retrieving user carts")
		return response, err
	}
	if totalUserCarts == 0 {
		return response, failure.BadRequest("No items in cart to purchase")
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
		log.Error().Err(err).Msg("[TransactionPurchaseCart] Service error retrieving products")
		return response, err
	}

	productMap := make(map[int]models.Product)
	for _, product := range products {
		productMap[product.Id] = product
	}

	var totalCost float32 = 0
	for _, cart := range userCarts {
		product, ok := productMap[cart.ProductId]
		if ok {
			totalCost += float32(cart.Quantity) * product.Price
		}
	}

	user, err := s.Repository.UserFindById(models.UserPrimaryId{Id: uuid.FromStringOrNil(req.UserId)})
	if err != nil {
		log.Error().Err(err).Msg("[TransactionPurchaseCart] Service error retrieving user")
		return response, err
	}

	if user.Balance < totalCost {
		return response, failure.BadRequest("Insufficient balance")
	}

	err = s.Repository.UserUpdateById(models.UserPrimaryId{Id: user.Id}, &models.User{
		Balance:   user.Balance - totalCost,
		UpdatedBy: req.Email,
	})
	if err != nil {
		log.Error().Err(err).Msg("[TransactionPurchaseCart] Service error updating user balance")
		return response, err
	}

	invoiceNumber := invoice.GenerateNumber()
	err = s.Repository.UserTransactionCreate(&models.UserTransaction{
		UserId:              uuid.FromStringOrNil(req.UserId),
		TransactionStatusId: models.MasterTransactionStatusId.Complete,
		TotalAmount:         float32(totalCost),
		InvoiceNumber:       invoiceNumber,
		CreatedBy:           req.Email,
		UpdatedBy:           req.Email,
	})
	if err != nil {
		log.Error().Err(err).Msg("[TransactionPurchaseCart] Service error creating transaction record")
		return response, err
	}

	/* Mark each cart item as purchased */
	for _, cart := range userCarts {
		updatedCart := cart
		updatedCart.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
		updatedCart.DeletedBy = null.StringFrom(req.Email)

		err = s.Repository.UserCartUpdateById(models.UserCartPrimaryId{Id: cart.Id}, &updatedCart)
		if err != nil {
			log.Error().Err(err).Msg("[TransactionPurchaseCart] Service error updating user cart")
			return response, err
		}
	}

	response = dto.TransactionPurchaseCartResponse{
		InvoiceNumber: invoiceNumber,
		TotalAmount:   totalCost,
	}

	return response, nil
}
