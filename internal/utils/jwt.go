package utils

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"randomMeetsProject/config"
	"randomMeetsProject/internal/models/sql_models"
	"randomMeetsProject/internal/models/validators"
)

func NewTokens(user sql_models.User) (validators.JWTTokens, error) {
	cfg, _ := config.LoadConfig("config.toml")
	accessToken := jwt.New(jwt.SigningMethodHS256)
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	accessExpTime := time.Now().Add(time.Minute * time.Duration(cfg.JWT.AccessExpireMinutes)).Unix()
	refreshExpTime := time.Now().Add(time.Minute * 60 * 24 * time.Duration(cfg.JWT.RefreshExpireDays)).Unix()

	accessClaims := accessToken.Claims.(jwt.MapClaims)
	accessClaims["uid"] = user.ID
	accessClaims["email"] = user.Email
	accessClaims["username"] = user.Username
	accessClaims["exp"] = accessExpTime

	refreshClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshClaims["uid"] = user.ID
	refreshClaims["email"] = user.Email
	refreshClaims["username"] = user.Username
	refreshClaims["exp"] = refreshExpTime

	accessTokenString, accessErr := accessToken.SignedString([]byte(cfg.JWT.SecretKey))
	refreshTokenString, refreshErr := refreshToken.SignedString([]byte(cfg.JWT.SecretKey))

	if refreshErr != nil || accessErr != nil {
		return validators.JWTTokens{}, errors.New("Error signing token")
	}

	tokens := validators.JWTTokens{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}

	return tokens, nil
}

func VerifyToken(tokenString string, secretKey string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return nil, errors.New("token expired")
			}
		} else {
			return nil, errors.New("invalid expiration time")
		}

		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func GetUserIDbyToken(tokenString string) (string, error) {
	cfg, err := config.LoadConfig("config.toml")
	if err != nil {
		return "", err
	}
	secretKey := cfg.JWT.SecretKey
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userID, ok := claims["uid"].(string); ok {
			return userID, nil
		}
	}

	return "", errors.New("invalid token")
}

func GetTokenFromAuthHeader(str string) string {
	token := strings.Split(str, " ")[1]
	return token
}
