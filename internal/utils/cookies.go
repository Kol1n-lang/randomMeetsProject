package utils

import (
	"net/http"

	"randomMeetsProject/internal/models/validators"
)

func CreateCookies(tokens validators.JWTTokens) map[string]*http.Cookie {
	cookies := make(map[string]*http.Cookie)

	accessCookie := new(http.Cookie)
	accessCookie.Name = "access_token"
	accessCookie.Value = tokens.AccessToken
	accessCookie.HttpOnly = true
	accessCookie.Path = "/"

	refreshCookie := new(http.Cookie)
	refreshCookie.Name = "refresh_token"
	refreshCookie.Value = tokens.RefreshToken
	refreshCookie.HttpOnly = true
	refreshCookie.Path = "/"

	cookies[accessCookie.Name] = accessCookie
	cookies[refreshCookie.Name] = refreshCookie

	return cookies
}
