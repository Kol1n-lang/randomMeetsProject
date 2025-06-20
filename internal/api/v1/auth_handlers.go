package v1

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"

	"randomMeetsProject/internal/models/validators"
	"randomMeetsProject/internal/services"
	"randomMeetsProject/internal/utils"
	"randomMeetsProject/pkg/docs"
)

func AuthGroup(e *echo.Echo) {
	authGroup := e.Group("/api/v1/auth")
	authGroup.POST("/register", registerUser)
	authGroup.POST("/login", loginUser)
	authGroup.GET("/me", me)
	authGroup.POST("/add-photo-url", addPhotoUrl)
	authGroup.GET("/confirm-email/", confirmUser)
}

func registerUser(c echo.Context) error {
	loginUser := new(validators.LoginUser)
	if err := c.Bind(loginUser); err != nil {
		return utils.CustomError(c, 422, err)
	}
	if err := loginUser.Validate(); err != nil {
		return utils.CustomError(c, 422, err)
	}
	authService, err := services.NewAuthService()
	if err != nil {
		return utils.CustomError(c, 500, err)
	}
	result, err := authService.RegisterUser(loginUser)
	if err != nil {
		return utils.CustomError(c, 500, err)
	}

	return c.JSON(http.StatusOK, result)
}

func loginUser(c echo.Context) error {
	fmt.Println("loginUser")
	loginUser := new(validators.LoginUser)
	if err := c.Bind(loginUser); err != nil {
		return utils.CustomError(c, 422, err)
	}
	if err := loginUser.Validate(); err != nil {
		return utils.CustomError(c, 422, err)
	}

	authService, err := services.NewAuthService()
	if err != nil {
		return utils.CustomError(c, 500, err)
	}
	jwt, err := authService.LoginUser(loginUser)
	if err != nil {
		return utils.CustomError(c, 500, err)
	}

	cookies := utils.CreateCookies(jwt)
	c.SetCookie(cookies["access_token"])
	c.SetCookie(cookies["refresh_token"])

	return c.JSON(http.StatusOK, docs.Response{
		Message: jwt,
	})
}

func confirmUser(c echo.Context) error {
	userID := uuid.MustParse(c.QueryParam("user_id"))
	authService, err := services.NewAuthService()
	if err != nil {
		return utils.CustomError(c, 500, err)
	}
	err = authService.ConfirmEmail(userID)
	if err != nil {
		return utils.CustomError(c, 500, err)
	}
	return c.JSON(http.StatusOK, docs.Response{
		Message: "Вы успешно подтвепдили email",
	})
}

func me(c echo.Context) error {
	token := utils.GetTokenFromAuthHeader(c.Request().Header.Get("Authorization"))
	userID, _ := utils.GetUserIDbyToken(token)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return utils.CustomError(c, 422, err)
	}
	authService, err := services.NewAuthService()
	if err != nil {
		return utils.CustomError(c, 500, err)
	}

	result, err := authService.GetMe(userUUID)
	return c.JSON(http.StatusOK, docs.Response{
		Message: result,
	})
}

func addPhotoUrl(c echo.Context) error {
	url := new(validators.PhotoURL)
	if err := c.Bind(url); err != nil {
		return utils.CustomError(c, 422, err)
	}
	if err := url.Validate(); err != nil {
		return utils.CustomError(c, 422, err)
	}
	token := utils.GetTokenFromAuthHeader(c.Request().Header.Get("Authorization"))
	userID, _ := utils.GetUserIDbyToken(token)
	userUUID, _ := uuid.Parse(userID)
	authService, _ := services.NewAuthService()
	err := authService.AddPhotoUrl(userUUID, url.URL)
	if err != nil {
		return utils.CustomError(c, 500, err)
	}
	return c.JSON(http.StatusOK, docs.Response{
		Message: "success",
	})
}
