package utils

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"randomMeetsProject/pkg/docs"
)

func CustomError(c echo.Context, code int, err error) error {
	if code == 422 {
		return c.JSON(http.StatusUnprocessableEntity, docs.Response{
			Message: "Validation Error: " + err.Error(),
		})
	} else if code == 500 {
		return c.JSON(http.StatusInternalServerError, docs.Response{
			Message: "Server Error: " + err.Error(),
		})
	} else if code == 401 {
		return c.JSON(http.StatusUnauthorized, docs.Response{
			Message: "Unauthorized Error",
		})
	} else if code == 409 {
		return c.JSON(http.StatusConflict, docs.Response{
			Message: "Data already exists : " + err.Error(),
		})
	} else if code == 404 {
		return c.JSON(http.StatusNotFound, docs.Response{
			Message: "Data Not Found",
		})
	}
	return c.JSON(http.StatusOK, docs.Response{
		Message: code,
	})
}
