package loger

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitLogger(logPath string) error {
	// Создаем директорию для логов, если она не существует
	if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
		return err
	}

	// Открываем файл для записи логов
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	// Создаем конфигурацию для логгера
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Создаем core для записи в файл
	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(config.EncoderConfig),
		zapcore.AddSync(file),
		zapcore.InfoLevel,
	)

	// Создаем core для вывода в консоль
	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(config.EncoderConfig),
		zapcore.AddSync(os.Stdout),
		zapcore.InfoLevel,
	)

	// Объединяем cores
	core := zapcore.NewTee(fileCore, consoleCore)

	// Создаем логгер
	Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return nil
}

