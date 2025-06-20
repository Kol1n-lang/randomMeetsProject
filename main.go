package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"randomMeetsProject/config"
	v1 "randomMeetsProject/internal/api/v1"
	"randomMeetsProject/pkg/database"
)

func initServer() *echo.Echo {
	server := echo.New()
	return server
}

func initGroups(e *echo.Echo) {
	v1.AuthGroup(e)
	v1.NewMeetsGroup(e)
}

func main() {
	cfg, err := config.LoadConfig("config.toml")
	if err != nil {
		log.Fatal("Error loading config " + err.Error())
	}
	err = database.InitDB()
	if err != nil {
		log.Fatal("Error connecting to database " + err.Error())
	}

	server := initServer()
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost",
			"http://localhost:8000",
			"http://127.0.0.1:5500",
			"https://oddly-lucid-tilefish.cloudpub.ru",
			"https://oddly-lucid-tilefish.cloudpub.ru:443",
		},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			"X-Requested-With",
		},
		AllowCredentials: true,
		MaxAge:           86400,
	}))
	initGroups(server)
	err = server.Start(cfg.ServerURL())
	if err != nil {
		log.Fatal("Error starting server " + err.Error())
	}
}
