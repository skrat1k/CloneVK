package services

import (
	"time"

	logger "CloneVK/pkg/Logger"
	"log/slog"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("your_secret_key")

type JWTService interface {
	GenerateToken(userID int) (string, error)
	ValidateToken(tokenString string) (int, error)
}

type jwtService struct {
	Log *slog.Logger
}

func NewJWTService(log *slog.Logger) JWTService {
	lg := logger.WithService(log, "JWTService")
	return &jwtService{Log: lg}
}

func (s *jwtService) GenerateToken(userID int) (string, error) {
	s.Log.Debug("Generating token", slog.Int("userID", userID))

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // токен живёт сутки
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		s.Log.Error("Failed to sign token", slog.Int("userID", userID), slog.String("error", err.Error()))
		return "", err
	}

	s.Log.Info("Token successfully generated", slog.Int("userID", userID))
	return tokenString, nil
}

func (s *jwtService) ValidateToken(tokenString string) (int, error) {
	s.Log.Debug("Validating token")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		s.Log.Warn("Invalid token", slog.String("error", err.Error()))
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		userID := int(claims["user_id"].(float64))
		s.Log.Info("Token successfully validated", slog.Int("userID", userID))
		return userID, nil
	}

	s.Log.Error("Invalid token claims structure")
	return 0, jwt.ErrTokenInvalidClaims
}
