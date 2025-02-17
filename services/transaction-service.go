package services

import (
	"github.com/nexmedis-be-technical-test/models"
	"github.com/nexmedis-be-technical-test/models/dto"
	"github.com/nexmedis-be-technical-test/utils/failure"
	"github.com/rs/zerolog/log"
)

func (s *Service) TransactionTopUpBalanceByUserId(req dto.TransactionTopUpBalanceByUserIdRequest) (dto.TransactionTopUpBalanceByUserIdResponse, error) {
	var response dto.TransactionTopUpBalanceByUserIdResponse
	users, _, err := s.Repository.UserFindManyAndCountByFilter(models.Filter{
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

	users, _, err := s.Repository.UserFindManyAndCountByFilter(models.Filter{
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
