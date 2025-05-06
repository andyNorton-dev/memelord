package indexer

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

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

func (h *IndexHandler) GetIndex(c echo.Context) error {
	return h.service.GetCurrentIndex(c)
}

func (h *IndexHandler) AddIndex(c echo.Context) error {
	return h.service.IncrementIndex(c)
}

func (h *IndexHandler) DoubleIndex(c echo.Context) error {
	return h.service.DoubleIndex(c)
}

func (h *IndexHandler) AddToIndex(c echo.Context) error {
	return h.service.AddToIndex(c)
}

// SetupIndexer инициализирует сервис и регистрирует роуты
func SetupIndexer(e *echo.Echo, db *sql.DB) {
	repo := NewIndexRepository(db)
	service := NewIndexService(repo)
	handler := NewIndexHandler(service)
	handler.RegisterRoutes(e)
}