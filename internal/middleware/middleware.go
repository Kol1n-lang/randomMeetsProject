package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"randomMeetsProject/config"
	"randomMeetsProject/internal/utils"
)

//	func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
//		return func(c echo.Context) error {
//			cfg, err := config.LoadConfig("config.toml")
//			if err != nil {
//				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "server configuration error"})
//			}
//
//			accessToken, err := c.Cookie("access_token")
//			accessTokenString := accessToken.Value
//			if err != nil {
//				return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
//			}
//			if accessTokenString == "" {
//				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "authorization header required"})
//			}
//
//			claims, err := utils.VerifyToken(accessTokenString, cfg.JWT.SecretKey)
//			if err != nil {
//				return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
//			}
//
//			c.Set("userClaims", claims)
//
//			return next(c)
//		}
//	}
func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cfg, err := config.LoadConfig("config.toml")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "server configuration error"})
		}

		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "authorization header required"})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid authorization format"})
		}

		tokenString := parts[1]
		claims, err := utils.VerifyToken(tokenString, cfg.JWT.SecretKey)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		}

		c.Set("userClaims", claims)

		return next(c)
	}
}
