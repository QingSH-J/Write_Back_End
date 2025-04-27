package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jinxinyu/go_backend/internal/config"
)

// define a struct to represent the token
type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

// define a struct to represent the token
type ToKenGenerator interface {
	GenerateToken(userID uuid.UUID, email string) (string, error)
	ValidateToken(token string) (*Claims, error)
}

type jwtTokenGenerator struct {
	SecretKey     []byte
	TokenDuration time.Duration
}

func NewTokenGenerator(cfg *config.Config) (ToKenGenerator, error) {
	if cfg.JWTSecret == "" {
		return nil, errors.New("JWT_SECRET is not set")
	}
	if cfg.JWTExpirationMinutes <= 0 {
		return nil, errors.New("JWT_EXPIRATION_MINUTES must bigger than 0")
	}
	return &jwtTokenGenerator{
		SecretKey:     []byte(cfg.JWTSecret),
		TokenDuration: time.Duration(cfg.JWTExpirationMinutes) * time.Minute,
	}, nil
}

func (t *jwtTokenGenerator) GenerateToken(userID uuid.UUID, email string) (string, error) {
	expirationTime := time.Now().Add(t.TokenDuration)
	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "go_backend",
			Subject:   userID.String(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(t.SecretKey)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return tokenString, nil
}

func (t *jwtTokenGenerator) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return t.SecretKey, nil
	})
	//check some errors like invalid token, expired token, etc.
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, fmt.Errorf("token is expired or not valid yet")
		}
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("invalid token")
		}
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}
