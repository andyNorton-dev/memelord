package config

import (
	"os"
	"fmt"
	"github.com/joho/godotenv"
)

type Config struct {
	DB_HOST string
	DB_PORT string
	DB_USER string
	DB_PASSWORD string
	DB_NAME string
	TELEGRAM_BOT_TOKEN string
	TELEGRAM_WEBHOOK_URL string
	PORT string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	config := GetConfig()
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

func GetConfig() *Config {
	return &Config{
		DB_HOST: os.Getenv("DB_HOST"),
		DB_PORT: os.Getenv("DB_PORT"),
		DB_USER: os.Getenv("DB_USER"),
		DB_PASSWORD: os.Getenv("DB_PASSWORD"),
		DB_NAME: os.Getenv("DB_NAME"),
		TELEGRAM_BOT_TOKEN: os.Getenv("TELEGRAM_BOT_TOKEN"),
		PORT: os.Getenv("PORT"),
	}
}

func validateConfig(config *Config) error {
	requiredFields := map[string]string{
		"DB_HOST": config.DB_HOST,
		"DB_PORT": config.DB_PORT,
		"DB_USER": config.DB_USER,
		"DB_PASSWORD": config.DB_PASSWORD,
		"DB_NAME": config.DB_NAME,
		"TELEGRAM_BOT_TOKEN": config.TELEGRAM_BOT_TOKEN,
		"PORT": config.PORT,
	}

	for field, value := range requiredFields {
		if value == "" {
			return fmt.Errorf("required environment variable %s is not set", field)
		}
	}

	return nil
}

