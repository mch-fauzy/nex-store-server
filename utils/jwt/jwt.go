package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/nexmedis-be-technical-test/configs"
	"github.com/nexmedis-be-technical-test/models/dto"
	"github.com/nexmedis-be-technical-test/utils/constant"
	"github.com/rs/zerolog/log"
)

func SignJwtToken(req dto.AuthTokenPayload, tokenType string, expiryDuration time.Duration) (dto.AuthLoginResponse, error) {
	config := configs.Get()

	expireTime := time.Now().Add(expiryDuration)
	expireTimeUnix := expireTime.Unix()
	expireTimeString := expireTime.Format(constant.DateTimeUTCFormat)

	// Create a new token with standard and custom claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": req.UserId,
		"email":  req.Email,
		"role":   req.Role,
		"exp":    expireTimeUnix,
	})

	// Sign the token with the provided secret key
	tokenString, err := token.SignedString([]byte(config.App.JwtAccessKey))
	if err != nil {
		log.Error().Err(err).Msg("[SignJWTToken] Failed to sign JWT Token")
		return dto.AuthLoginResponse{}, err
	}

	return dto.AuthLoginResponse{
		Token:     tokenString,
		TokenType: tokenType,
		ExpiresIn: expireTimeString,
	}, nil
}
