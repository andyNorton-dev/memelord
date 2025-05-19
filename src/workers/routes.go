package workers

import (
	"github.com/labstack/echo/v4"
	"api/src/middleware"
	"api/src/core/config"
)

type WorkerHandler struct {
	service *WorkerService
	config *config.Config
}

func NewWorkerHandler(service *WorkerService, config *config.Config) *WorkerHandler {
	return &WorkerHandler{service: service, config: config}
}

func RegisterRoutes(e *echo.Echo, service *WorkerService, config *config.Config) {
	handler := NewWorkerHandler(service, config)
	
	workersGroup := e.Group("/workers")
	workersGroup.Use(middleware.TelegramAuth(middleware.TelegramAuthConfig{
		BotToken: handler.config.TELEGRAM_BOT_TOKEN,
	}))
	workersGroup.GET("/worker", handler.GetWorkers)
	workersGroup.GET("/army", handler.GetArmy)
	workersGroup.POST("/buy/:id", handler.BuyWorker)
}

// @Summary Получить список работников
// @Description Получает список доступных работников
// @Tags workers
// @Accept json
// @Produce json
// @Success 200 {array} UserWorkerResponse
// @Failure 500 {object} map[string]string
// @Router /workers/worker [get]
// @Security TelegramAuth
func (h *WorkerHandler) GetWorkers(c echo.Context) error {
	return h.service.GetWorkers(c)
}

// @Summary Получить армию
// @Description Получает список доступных воинов
// @Tags workers
// @Accept json
// @Produce json
// @Success 200 {array} UserWorkerResponse
// @Failure 500 {object} map[string]string
// @Router /workers/army [get]
// @Security TelegramAuth
func (h *WorkerHandler) GetArmy(c echo.Context) error {
	return h.service.GetArmy(c)
}

// @Summary Купить работника
// @Description Покупает или улучшает работника
// @Tags workers
// @Accept json
// @Produce json
// @Param id path int true "ID работника"
// @Success 200 {array} UserWorkerResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /workers/buy/{id} [post]
// @Security TelegramAuth
func (h *WorkerHandler) BuyWorker(c echo.Context) error {
	return h.service.BuyWorker(c)
}