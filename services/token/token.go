package token

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	ExpireCount        = 2
	ExpireRefreshCount = 168
)

type JwtCustomClaims struct {
	UserID    uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	CiSession int       `json:"ci_session"`
	jwt.RegisteredClaims
}
