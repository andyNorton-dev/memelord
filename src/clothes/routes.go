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
	clothesGroup.GET("", handler.GetClothes)
	clothesGroup.GET("/:id", handler.GetClothe)
	clothesGroup.POST("/buy/:id", handler.BuyClothe)
	clothesGroup.POST("/equip/:id", handler.EquipClothe)
}

// @Summary Получить список одежды
// @Description Получает список доступной одежды
// @Tags clothes
// @Accept json
// @Produce json
// @Success 200 {array} ClotheUserResponse
// @Failure 500 {object} map[string]string
// @Router /clothes [get]
// @Security TelegramAuth
func (h *ClothesHandler) GetClothes(c echo.Context) error {
	return h.service.GetClothes(c)
}

// @Summary Получить предмет одежды
// @Description Получает информацию о конкретном предмете одежды
// @Tags clothes
// @Accept json
// @Produce json
// @Param id path string true "ID предмета одежды"
// @Success 200 {object} ClotheUserResponse
// @Failure 500 {object} map[string]string
// @Router /clothes/{id} [get]
// @Security TelegramAuth
func (h *ClothesHandler) GetClothe(c echo.Context) error {
	return h.service.GetClothe(c)
}

// @Summary Купить предмет одежды
// @Description Покупает предмет одежды
// @Tags clothes
// @Accept json
// @Produce json
// @Param id path string true "ID предмета одежды"
// @Success 200 {object} ClotheUserResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /clothes/buy/{id} [post]
// @Security TelegramAuth
func (h *ClothesHandler) BuyClothe(c echo.Context) error {
	return h.service.BuyClothe(c)
}

// @Summary Экипировать предмет одежды
// @Description Экипирует предмет одежды
// @Tags clothes
// @Accept json
// @Produce json
// @Param id path string true "ID предмета одежды"
// @Success 200 {object} ClotheUserResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /clothes/equip/{id} [post]
// @Security TelegramAuth
func (h *ClothesHandler) EquipClothe(c echo.Context) error {
	return h.service.EquipClothe(c)
}

