package services

import (
	"time"

	logger "CloneVK/pkg/Logger"
	"log/slog"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("your_secret_key")

const (
	accessTokenTTL  = 15 * time.Minute
	refreshTokenTTL = 7 * 24 * time.Hour
)

type TokenPair struct {
	AccessToken           string
	RefreshToken          string
	RefreshTokenExpiresAt time.Time
}

type JWTService interface {
	GenerateTokens(userID int) (*TokenPair, error)
	GenerateAccessToken(userID int) (string, error)
	ValidateToken(tokenString string) (int, error)
}

type jwtService struct {
	Log *slog.Logger
}

func NewJWTService(log *slog.Logger) JWTService {
	lg := logger.WithService(log, "JWTService")
	return &jwtService{Log: lg}
}

func (s *jwtService) GenerateTokens(userID int) (*TokenPair, error) {
	s.Log.Debug("Generating access and refresh tokens", slog.Int("userID", userID))

	accessToken, err := s.generateToken(userID, accessTokenTTL)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshExpiresAt := s.generateTokenWithExpiry(userID, refreshTokenTTL)
	if refreshToken == "" {
		return nil, err
	}

	return &TokenPair{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshExpiresAt,
	}, nil
}

func (s *jwtService) GenerateAccessToken(userID int) (string, error) {
	s.Log.Debug("Generating access token", slog.Int("userID", userID))
	return s.generateToken(userID, accessTokenTTL)
}

func (s *jwtService) generateToken(userID int, ttl time.Duration) (string, error) {
	s.Log.Debug("Generating token", slog.Int("userID", userID))

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(ttl).Unix(),
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

func (s *jwtService) generateTokenWithExpiry(userID int, ttl time.Duration) (string, time.Time) {
	expiresAt := time.Now().Add(ttl)

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expiresAt.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		s.Log.Error("Failed to sign token", slog.Int("userID", userID), slog.String("error", err.Error()))
		return "", time.Time{}
	}

	s.Log.Info("Token successfully generated", slog.Int("userID", userID))
	return tokenString, expiresAt
}
