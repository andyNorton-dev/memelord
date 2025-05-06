package workers

import (
	"github.com/labstack/echo/v4"
	"api/src/middleware"
)

type WorkerHandler struct {
	service *WorkerService
}

func NewWorkerHandler(service *WorkerService) *WorkerHandler {
	return &WorkerHandler{service: service}
}

func RegisterRoutes(e *echo.Echo, service *WorkerService) {
	handler := NewWorkerHandler(service)
	
	workersGroup := e.Group("/workers")
	workersGroup.Use(middleware.TelegramAuth(middleware.TelegramAuthConfig{
		BotToken: "6885676739:AAFP8P6v51rXXdQzpH04EhQNdPVpHVJ-26Y",
	}))
	workersGroup.GET("/worker", handler.GetWorkers)
	workersGroup.GET("/army", handler.GetArmy)
	workersGroup.POST("/buy/:id", handler.BuyWorker)
}

func (h *WorkerHandler) GetWorkers(c echo.Context) error {
	return h.service.GetWorkers(c)
}

func (h *WorkerHandler) GetArmy(c echo.Context) error {
	return h.service.GetArmy(c)
}

func (h *WorkerHandler) BuyWorker(c echo.Context) error {
	return h.service.BuyWorker(c)
}