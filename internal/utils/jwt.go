package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken  = errors.New("invalid token")
	ErrExpiredToken  = errors.New("token has expired")
	ErrInvalidClaims = errors.New("invalid token claims")
)

// Claims represents the JWT claims
type Claims struct {
	UserID   int64  `json:"user_id"`
	Email    string `json:"email"`
	UserType string `json:"user_type"`
	jwt.RegisteredClaims
}

type JWTConfig struct {
	SecretKey            string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

func GenerateAccessToken(userID int64, email, userType, secretKey string, duration time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:   userID,
		Email:    email,
		UserType: userType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "keerja-api",
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

func GenerateRefreshToken(userID int64, email, userType, secretKey string, duration time.Duration) (string, error) {
	return GenerateAccessToken(userID, email, userType, secretKey, duration)
}

func ValidateToken(tokenString, secretKey string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
		},
	)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidClaims
	}

	return claims, nil
}

func ExtractUserID(claims *Claims) int64 {
	if claims == nil {
		return 0
	}
	return claims.UserID
}

func ExtractEmail(claims *Claims) string {
	if claims == nil {
		return ""
	}
	return claims.Email
}

func ExtractUserType(claims *Claims) string {
	if claims == nil {
		return ""
	}
	return claims.UserType
}

func IsTokenExpired(claims *Claims) bool {
	if claims == nil {
		return true
	}
	return claims.ExpiresAt.Before(time.Now())
}

func GetTokenExpirationTime(claims *Claims) time.Time {
	if claims == nil {
		return time.Time{}
	}
	return claims.ExpiresAt.Time
}

func GetTokenRemainingTime(claims *Claims) time.Duration {
	if claims == nil {
		return 0
	}
	return time.Until(claims.ExpiresAt.Time)
}
