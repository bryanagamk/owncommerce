package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

type Claims struct {
	UserID    uuid.UUID  `json:"user_id"`
	TenantID  *uuid.UUID `json:"tenant_id,omitempty"`
	TokenType TokenType  `json:"token_type"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type JWTManager struct {
	secret        []byte
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

func NewJWTManager(secret string, accessExpiry, refreshExpiry time.Duration) *JWTManager {
	return &JWTManager{
		secret:        []byte(secret),
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
	}
}

func (m *JWTManager) GenerateTokenPair(userID uuid.UUID, tenantID *uuid.UUID) (*TokenPair, string, error) {
	accessToken, err := m.generateToken(userID, tenantID, TokenTypeAccess, m.accessExpiry)
	if err != nil {
		return nil, "", err
	}

	refreshToken, err := m.generateToken(userID, tenantID, TokenTypeRefresh, m.refreshExpiry)
	if err != nil {
		return nil, "", err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(m.accessExpiry.Seconds()),
		TokenType:    "Bearer",
	}, HashToken(refreshToken), nil
}

func (m *JWTManager) generateToken(userID uuid.UUID, tenantID *uuid.UUID, tokenType TokenType, expiry time.Duration) (string, error) {
	now := time.Now().UTC()
	claims := Claims{
		UserID:    userID,
		TenantID:  tenantID,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "owncommerce",
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

func (m *JWTManager) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func (m *JWTManager) RefreshExpiry() time.Duration {
	return m.refreshExpiry
}

func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
