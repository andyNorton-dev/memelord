package custom_middleware

import (
	"api/src/core/loger"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// LoggerMiddleware возвращает middleware для логирования HTTP-запросов
func LoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			req := c.Request()
			res := c.Response()

			// Выполняем следующий обработчик
			err := next(c)

			// Логируем информацию о запросе
			loger.Logger.Info("HTTP Request",
				zap.String("method", req.Method),
				zap.String("path", req.URL.Path),
				zap.String("remote_ip", c.RealIP()),
				zap.Int("status", res.Status),
				zap.Int64("latency", time.Since(start).Milliseconds()),
				zap.String("user_agent", req.UserAgent()),
			)

			return err
		}
	}
} 