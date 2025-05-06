package clothes

import (
	"github.com/labstack/echo/v4"
	"api/src/middleware"
)

type ClothesHandler struct {
	service *ClothesService
}

func NewClothesHandler(service *ClothesService) *ClothesHandler {
	return &ClothesHandler{service: service}
}

func RegisterRoutes(e *echo.Echo, service *ClothesService) {
	handler := NewClothesHandler(service)
	
	clothesGroup := e.Group("/clothes")
	clothesGroup.Use(middleware.TelegramAuth(middleware.TelegramAuthConfig{
		BotToken: "6885676739:AAFP8P6v51rXXdQzpH04EhQNdPVpHVJ-26Y",
	}))
	clothesGroup.GET("/", handler.GetClothes)
	clothesGroup.GET("/:id", handler.GetClothe)
	clothesGroup.POST("/buy/:id", handler.BuyClothe)
	clothesGroup.POST("/equip/:id", handler.EquipClothe)
}

func (h *ClothesHandler) GetClothes(c echo.Context) error {
	return h.service.GetClothes(c)
}

func (h *ClothesHandler) GetClothe(c echo.Context) error {
	return h.service.GetClothe(c)
}

func (h *ClothesHandler) BuyClothe(c echo.Context) error {
	return h.service.BuyClothe(c)
}

func (h *ClothesHandler) EquipClothe(c echo.Context) error {
	return h.service.EquipClothe(c)
}

