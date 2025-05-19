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
	"api/src/core/config"
	"api/src/core/loger"
	customMiddleware "api/src/core/middleware"
	"go.uber.org/zap"
	"net/http"
)

// @title API Documentation
// @version 1.0
// @description API сервер для управления пользователями, работниками и одеждой
// @host localhost:8081
// @BasePath /

// @securityDefinitions.apikey TelegramAuth
// @in header
// @name X-Telegram-Init-Data
// @description Данные инициализации Telegram WebApp
func main() {
	// Инициализация логгера
	if err := loger.InitLogger("logs/app.log"); err != nil {
		panic(err)
	}
	defer loger.Logger.Sync()

	config, err := config.LoadConfig()
	if err != nil {
		loger.Logger.Fatal("Ошибка загрузки конфигурации", zap.Error(err))
	}
	loger.Logger.Info("Конфигурация загружена", zap.Any("config", config))
	
	e := echo.New()
	
	// Кастомный обработчик ошибок
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		var statusCode int
		var message string

		if he, ok := err.(*echo.HTTPError); ok {
			statusCode = he.Code
			message = he.Message.(string)
		} else {
			statusCode = http.StatusInternalServerError
			message = "Внутренняя ошибка сервера"
		}

		// Логируем ошибку
		loger.Logger.Error("HTTP Error",
			zap.Int("status_code", statusCode),
			zap.String("message", message),
			zap.String("path", c.Request().URL.Path),
			zap.String("method", c.Request().Method),
			zap.Error(err),
		)

		// Отправляем ответ клиенту
		if !c.Response().Committed {
			if c.Request().Header.Get("Content-Type") == "application/json" {
				c.JSON(statusCode, map[string]string{
					"error": message,
				})
			} else {
				c.String(statusCode, message)
			}
		}
	}
	
	// Middleware
	e.Use(customMiddleware.LoggerMiddleware())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	
	// Swagger endpoint
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	
	// Подключение к базе данных
	database, err := db.Connect(config)
	if err != nil {
		loger.Logger.Fatal("Ошибка подключения к базе данных", zap.Error(err))
	}
	defer database.Close()
	
	loger.Logger.Info("Подключение к базе данных установлено")
	
	// Setup routes
	indexer.SetupIndexer(e, database)
	
	// Создаем сервис пользователей
	userRepo := user.NewUserRepository(database)
	userService := user.NewService(userRepo)
	user.SetupUser(e, database, config)
	
	// Создаем сервис для работников
	workerRepo := workers.NewWorkerRepository(database)
	workerService := workers.NewWorkerService(workerRepo, userService)
	workers.RegisterRoutes(e, workerService, config)

	// Создаем сервис для одежды
	clothesRepo := clothes.NewClothesRepository(database)
	clothesService := clothes.NewClothesService(clothesRepo, userService)
	clothes.RegisterRoutes(e, clothesService, config)
	
	loger.Logger.Info("Сервер запускается", zap.String("port", config.PORT))
	
	// Start server
	if err := e.Start(":" + config.PORT); err != nil {
		loger.Logger.Fatal("Ошибка запуска сервера", zap.Error(err))
	}
}