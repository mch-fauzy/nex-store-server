package dto

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
	"github.com/nexmedis-be-technical-test/models"
	"github.com/nexmedis-be-technical-test/utils/failure"
)

type AuthRegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	RoleId   int    `json:"-"`
}

func (r AuthRegisterRequest) Validate() error {
	if r.Email == "" {
		return failure.BadRequest("Email is required")
	}

	if r.Password == "" {
		return failure.BadRequest("Password is required")
	}

	if len(r.Password) < 8 {
		return failure.BadRequest("Password must be at least 8 characters")
	}

	return nil
}

func (r AuthRegisterRequest) ToModel() models.User {
	id, _ := uuid.NewV4()
	return models.User{
		Id:        id,
		RoleId:    r.RoleId,
		Email:     r.Email,
		Password:  r.Password,
		LastLogin: nil,
		CreatedBy: r.Email,
		UpdatedBy: r.Email,
	}
}

type AuthLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r AuthLoginRequest) Validate() error {
	if r.Email == "" {
		return failure.BadRequest("Email is required")
	}

	if r.Password == "" {
		return failure.BadRequest("Password is required")
	}

	return nil
}

func (r AuthLoginRequest) UpdateLastLogin() models.User {
	now := time.Now()
	return models.User{
		LastLogin: &now,
	}
}

type AuthLoginResponse struct {
	Token     string `json:"token"`
	TokenType string `json:"tokenType"`
	ExpiresIn string `json:"expiresIn"`
}

type AuthTokenPayload struct {
	UserId string `json:"userId"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.StandardClaims
}
