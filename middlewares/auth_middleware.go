package middlewares

import (
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"

	s "github.com/Project-IPCA/ipca-realtime-go/server"
	"github.com/Project-IPCA/ipca-realtime-go/services/token"
)

type AuthMiddleware struct {
	server *s.Server
}

func NewAuthMiddleware(server *s.Server) *AuthMiddleware {
	return &AuthMiddleware{
		server: server,
	}
}

func (authMiddleware *AuthMiddleware) GetJwtConfig() echojwt.Config {
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(token.JwtCustomClaims)
		},
		SigningKey: []byte(authMiddleware.server.Config.Auth.AccessSecret),
		TokenLookup: "query:token",
	}
	return config
}
