package services

import (
	"time"

	"github.com/nexmedis-be-technical-test/models"
	"github.com/nexmedis-be-technical-test/models/dto"
	"github.com/nexmedis-be-technical-test/utils/constant"
	"github.com/nexmedis-be-technical-test/utils/failure"
	"github.com/nexmedis-be-technical-test/utils/jwt"
	"github.com/nexmedis-be-technical-test/utils/password"
	"github.com/rs/zerolog/log"
)

func (s *Service) AuthRegister(req dto.AuthRegisterRequest) (string, error) {
	message := "Failed"

	// Check if email already exists
	_, totalUsers, err := s.Repository.UserFindManyAndCountByFilter(models.Filter{
		FilterFields: []models.FilterField{
			{
				Field:    models.UserDbField.Email,
				Operator: models.OperatorEqual,
				Value:    req.Email,
			},
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("[AuthRegister] Service error getting users")
		return message, err
	}

	if totalUsers > 0 {
		err = failure.Conflict("User with this email already exists")
		return message, err
	}

	hashedPassword, err := password.HashPassword(req.Password)
	if err != nil {
		log.Error().Err(err).Msg("[AuthRegister] Service error hashing password")
		return message, err
	}

	req.Password = hashedPassword
	user := req.ToModel()
	err = s.Repository.UserCreate(&user)
	if err != nil {
		log.Error().Err(err).Msg("[AuthRegister] Service error creating user")
		return message, err
	}

	message = "Success"
	return message, nil
}

func (s *Service) AuthLogin(req dto.AuthLoginRequest) (dto.AuthLoginResponse, error) {
	users, totalUsers, err := s.Repository.UserFindManyAndCountByFilter(models.Filter{
		FilterFields: []models.FilterField{
			{
				Field:    models.UserDbField.Email,
				Operator: models.OperatorEqual,
				Value:    req.Email,
			},
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("[AuthLogin] Service error getting users")
		return dto.AuthLoginResponse{}, err
	}

	if totalUsers == 0 {
		err = failure.Unauthorized("Invalid credentials")
		return dto.AuthLoginResponse{}, err
	}

	userRole, err := s.Repository.MasterUserRoleFindById(models.MasterUserRolePrimaryId{
		Id: users[0].RoleId},
	)
	if err != nil {
		log.Error().Err(err).Msg("[AuthLogin] Service error getting user role")
		return dto.AuthLoginResponse{}, err
	}

	err = password.ComparePassword(req.Password, users[0].Password)
	if err != nil {
		log.Error().Err(err).Msg("[AuthLogin] Service error comparing password")
		err = failure.Unauthorized("Invalid credentials")
		return dto.AuthLoginResponse{}, err
	}

	/* Update last login value*/
	lastLogin := req.UpdateLastLogin()
	err = s.Repository.UserUpdateById(models.UserPrimaryId{
		Id: users[0].Id,
	},
		&lastLogin,
	)
	if err != nil {
		log.Error().Err(err).Msg("[AuthLogin] Service error updating last login")
		return dto.AuthLoginResponse{}, err
	}

	response, err := jwt.SignJwtToken(dto.AuthTokenPayload{
		UserId: users[0].Id.String(),
		Email:  users[0].Email,
		Role:   userRole.Name,
	},
		constant.BearerTokenType,
		time.Hour,
	)
	if err != nil {
		log.Error().Err(err).Msg("[AuthLogin] Service error signing jwt token")
		err = failure.InternalError("Failed to login user")
		return dto.AuthLoginResponse{}, err
	}

	return response, nil
}
