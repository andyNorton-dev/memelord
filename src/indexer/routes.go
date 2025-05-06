package indexer

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

// @Summary Получить текущий индекс
// @Description Получает текущее значение индекса
// @Tags indexer
// @Accept json
// @Produce json
// @Success 200 {object} map[string]int
// @Failure 500 {object} ErrorResponse
// @Router /index [get]
// @Security TelegramAuth
func (h *IndexHandler) GetIndex(c echo.Context) error {
	return h.service.GetCurrentIndex(c)
}

// @Summary Увеличить индекс
// @Description Увеличивает значение индекса на 1
// @Tags indexer
// @Accept json
// @Produce json
// @Success 204 "No Content"
// @Failure 500 {object} ErrorResponse
// @Router /index [post]
// @Security TelegramAuth
func (h *IndexHandler) AddIndex(c echo.Context) error {
	return h.service.IncrementIndex(c)
}

// @Summary Удвоить индекс
// @Description Удваивает текущее значение индекса
// @Tags indexer
// @Accept json
// @Produce json
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /index/double [post]
// @Security TelegramAuth
func (h *IndexHandler) DoubleIndex(c echo.Context) error {
	return h.service.DoubleIndex(c)
}

// @Summary Добавить к индексу
// @Description Добавляет указанное число к текущему значению индекса
// @Tags indexer
// @Accept json
// @Produce json
// @Param number query int true "Число для добавления"
// @Success 200 {object} AddToIndexResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /index/add [post]
// @Security TelegramAuth
func (h *IndexHandler) AddToIndex(c echo.Context) error {
	return h.service.AddToIndex(c)
}

type IndexHandler struct {
	service *IndexService
}

func NewIndexHandler(service *IndexService) *IndexHandler {
	return &IndexHandler{
		service: service,
	}
}

func validateAddToIndex(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		numberStr := c.QueryParam("number")
		if numberStr == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "параметр number обязателен",
			})
		}
		number, err := strconv.Atoi(numberStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "параметр number должен быть числом",
			})
		}
		c.Set("request", &AddToIndexRequest{Number: number})
		return next(c)
	}
}

func (h *IndexHandler) RegisterRoutes(e *echo.Echo) {
	e.GET("/index", h.GetIndex)
	e.POST("/index", h.AddIndex)
	e.POST("/index/double", h.DoubleIndex)
	e.POST("/index/add", h.AddToIndex, validateAddToIndex)
}

// SetupIndexer инициализирует сервис и регистрирует роуты
func SetupIndexer(e *echo.Echo, db *sql.DB) {
	repo := NewIndexRepository(db)
	service := NewIndexService(repo)
	handler := NewIndexHandler(service)
	handler.RegisterRoutes(e)
}