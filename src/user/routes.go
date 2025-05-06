package user

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	"api/src/middleware"
)

type UserHandler struct {
	service UserService
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) RegisterRoutes(e *echo.Echo) {
	// Создаем группу роутов с middleware
	userGroup := e.Group("/user")
	
	// Применяем TelegramAuth middleware ко всем роутам в группе
	userGroup.Use(middleware.TelegramAuth(middleware.TelegramAuthConfig{
		BotToken: "6885676739:AAFP8P6v51rXXdQzpH04EhQNdPVpHVJ-26Y",
	}))

	userGroup.GET("", h.GetUser)
	userGroup.POST("", h.CreateUser)
	userGroup.POST("/tap", h.TapUser)
}

func (h *UserHandler) GetUser(c echo.Context) error {
	return h.service.GetUserHandler(c)
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	return h.service.CreateUserHandler(c)
}

func (h *UserHandler) TapUser(c echo.Context) error {
	return h.service.TapUserHandler(c)
}

func SetupUser(e *echo.Echo, db *sql.DB) {
	repo := NewUserRepository(db)
	service := NewService(repo)
	handler := NewUserHandler(service)
	handler.RegisterRoutes(e)
}



