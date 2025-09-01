package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secret []byte
}

func NewJWTService(secret string) *JWTService {
	return &JWTService{
		secret: []byte(secret),
	}
}

// CreateToken generates a new JWT token
func (j *JWTService) CreateToken(userID uint) (string, error) {
	// Add nil checks
	if j == nil {
		return "", errors.New("JWT service is nil")
	}

	if len(j.secret) == 0 {
		return "", errors.New("JWT secret is not configured")
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"iat":     time.Now().Unix(),
	})

	// Sign token
	return token.SignedString(j.secret)
}

func (j *JWTService) VerifyToken(tokenString string) (*jwt.MapClaims, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return j.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*jwt.MapClaims)

	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
