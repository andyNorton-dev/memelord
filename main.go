package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/swaggo/echo-swagger"
	_ "api/docs"
	"api/src/indexer"
	"api/src/user"
	"api/src/workers"
	"api/src/clothes"
	"api/db"
)

// @title API Documentation
// @version 1.0
// @description API сервер для управления пользователями, работниками и одеждой
// @host localhost:8081
// @BasePath /
func main() {
	e := echo.New()
	
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	
	// Swagger endpoint
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	
	// Подключение к базе данных
	database, err := db.Connect()
	if err != nil {
		e.Logger.Fatal(err)
	}
	defer database.Close()
	
	// Setup routes
	indexer.SetupIndexer(e, database)
	
	// Создаем сервис пользователей
	userRepo := user.NewUserRepository(database)
	userService := user.NewService(userRepo)
	user.SetupUser(e, database)
	
	// Создаем сервис для работников
	workerRepo := workers.NewWorkerRepository(database)
	workerService := workers.NewWorkerService(workerRepo, userService)
	workers.RegisterRoutes(e, workerService)

	// Создаем сервис для одежды
	clothesRepo := clothes.NewClothesRepository(database)
	clothesService := clothes.NewClothesService(clothesRepo, userService)
	clothes.RegisterRoutes(e, clothesService)
	
	// Start server
	e.Logger.Fatal(e.Start(":8081"))
}