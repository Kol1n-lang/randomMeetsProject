package v1

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"

	"randomMeetsProject/internal/middleware"
	"randomMeetsProject/internal/models/validators"
	"randomMeetsProject/internal/services"
	"randomMeetsProject/internal/utils"
	"randomMeetsProject/pkg/docs"
)

func NewMeetsGroup(c *echo.Echo) {
	meetsGroup := c.Group("/api/v1/meets")
	meetsGroup.Use(middleware.JWTMiddleware)
	meetsGroup.GET("/docs", meetDocs)
	meetsGroup.POST("/create-meetup", createMeetUp)
	meetsGroup.GET("/get-meetup", getMeetUp)
	meetsGroup.GET("/get-all-meetups", getAllMeetUps)
	meetsGroup.POST("/delete-meetup", deleteMeetUp)
	meetsGroup.PATCH("/update-meetup", updateMeetUp)
}

func meetDocs(c echo.Context) error {
	return c.JSON(http.StatusOK, docs.Response{
		Message: "Create Meets, Check Meets",
	})
}

func createMeetUp(c echo.Context) error {
	token := utils.GetTokenFromAuthHeader(c.Request().Header.Get("Authorization"))
	meetUp := new(validators.MeetUp)
	if err := c.Bind(meetUp); err != nil {
		return utils.CustomError(c, 422, err)
	}
	if err := meetUp.Validate(); err != nil {
		return utils.CustomError(c, 422, err)
	}
	dateErr := utils.CheckDate(meetUp.Date)
	if dateErr != nil {
		return utils.CustomError(c, 422, dateErr)
	}
	meetService, err := services.NewMeetUpService()
	if err != nil {
		return utils.CustomError(c, 500, err)
	}
	userID, err := utils.GetUserIDbyToken(token)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return utils.CustomError(c, 500, err)
	}
	result, err := meetService.CreateMeetUp(meetUp, userUUID)
	if err != nil {
		return utils.CustomError(c, 500, err)
	}
	return c.JSON(http.StatusOK, docs.Response{
		Message: result,
	})
}

func getMeetUp(c echo.Context) error {
	userID := c.QueryParam("user_id")
	meetService, err := services.NewMeetUpService()
	if err != nil {
		return utils.CustomError(c, 500, err)
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return utils.CustomError(c, 500, err)
	}
	result, err := meetService.GetMeetUp(userUUID)
	if err != nil {
		return utils.CustomError(c, 500, err)
	}
	return c.JSON(http.StatusOK, docs.Response{
		Message: result,
	})
}

func getAllMeetUps(c echo.Context) error {
	meetService, err := services.NewMeetUpService()
	if err != nil {
		return utils.CustomError(c, 500, err)
	}
	result, err := meetService.GetAllMeetUps(c.Request().Context())
	if err != nil {
		return utils.CustomError(c, 500, err)
	}
	return c.JSON(http.StatusOK, docs.Response{
		Message: result,
	})
}

func deleteMeetUp(c echo.Context) error {
	token := utils.GetTokenFromAuthHeader(c.Request().Header.Get("Authorization"))
	userID, err := utils.GetUserIDbyToken(token)
	if err != nil {
		return utils.CustomError(c, 500, err)
	}
	UserUUID, _ := uuid.Parse(userID)
	meetService, err := services.NewMeetUpService()
	if err != nil {
		return utils.CustomError(c, 500, err)
	}
	err = meetService.DeleteMeetUp(UserUUID)
	if err != nil {
		return utils.CustomError(c, 500, err)
	}
	return c.JSON(http.StatusOK, docs.Response{
		Message: "Already delete MeetUps",
	})
}

func updateMeetUp(c echo.Context) error {
	token := utils.GetTokenFromAuthHeader(c.Request().Header.Get("Authorization"))
	userID, err := utils.GetUserIDbyToken(token)
	if err != nil {
		return utils.CustomError(c, 500, err)
	}
	UserUUID, _ := uuid.Parse(userID)
	meetUp := new(validators.PutMeetUp)
	if err := c.Bind(meetUp); err != nil {
		return utils.CustomError(c, 422, err)
	}
	if err := meetUp.Validate(); err != nil {
		return utils.CustomError(c, 422, err)
	}

	if meetUp.Date != "" {
		errDate := utils.CheckDate(meetUp.Date)
		if errDate != nil {
			return utils.CustomError(c, 422, errDate)
		}
	}

	meetService, err := services.NewMeetUpService()
	if err != nil {
		return utils.CustomError(c, 500, err)
	}
	err = meetService.UpdateMeetUp(meetUp, UserUUID)
	if err != nil {
		return utils.CustomError(c, 500, err)
	}
	return c.JSON(http.StatusOK, docs.Response{
		Message: meetUp,
	})
}
