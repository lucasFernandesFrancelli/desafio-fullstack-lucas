// Package auth cuida da emissão e validação dos JWTs usados pelo "login"
// mock (seletor de perfil): o backend confia no userId enviado por
// /auth/login e emite um token real, sem senha.
package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"ekaizen-backend/internal/models"
)

var ErrInvalidToken = errors.New("token inválido ou expirado")

const tokenTTL = 24 * time.Hour

type Claims struct {
	UserID uuid.UUID   `json:"sub"`
	Role   models.Role `json:"role"`
	Name   string      `json:"name"`
	jwt.RegisteredClaims
}

type JWTIssuer struct {
	secret []byte
}

func NewJWTIssuer(secret string) *JWTIssuer {
	return &JWTIssuer{secret: []byte(secret)}
}

func (j *JWTIssuer) Issue(user models.User) (string, error) {
	claims := Claims{
		UserID: user.ID,
		Role:   user.Role,
		Name:   user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID.String(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

func (j *JWTIssuer) Parse(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return j.secret, nil
	})
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
