package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
)

// TelegramUser представляет данные пользователя из Telegram
type TelegramUser struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

// TelegramAuthConfig конфигурация для middleware
type TelegramAuthConfig struct {
	BotToken string
}

// TelegramAuth middleware проверяет подпись Telegram WebApp данных
func TelegramAuth(config TelegramAuthConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Получаем initData из заголовка
			initData := c.Request().Header.Get("X-Telegram-Init-Data")
			if initData == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "X-Telegram-Init-Data header is required",
				})
			}
			config.BotToken = "6885676739:AAFP8P6v51rXXdQzpH04EhQNdPVpHVJ-26Y"
			// Проверяем подпись
			ok, data, err := verifyTelegramWebAppData(config.BotToken, initData)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{
					"error": "Failed to verify Telegram data",
				})
			}
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid Telegram data signature",
				})
			}

			// Извлекаем user данные
			userJSON, exists := data["user"]
			if !exists {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "User data not found in initData",
				})
			}

			// Парсим JSON в структуру
			var telegramUser TelegramUser
			if err := json.Unmarshal([]byte(userJSON), &telegramUser); err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{
					"error": "Invalid user data format",
				})
			}

			// Устанавливаем структуру пользователя в контекст
			c.Set("telegram_user", &telegramUser)

			return next(c)
		}
	}
} 