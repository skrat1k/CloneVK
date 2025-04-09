package services

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("your_secret_key")

type JWTService interface {
	GenerateToken(userID int) (string, error)
	ValidateToken(tokenString string) (int, error)
}

type jwtService struct{}

func NewJWTService() JWTService {
	return &jwtService{}
}

func (s *jwtService) GenerateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(), //токен живёт сутки
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func (s *jwtService) ValidateToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return 0, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		userID := int(claims["user_id"].(float64))
		return userID, nil
	}
	return 0, jwt.ErrTokenInvalidClaims
}
