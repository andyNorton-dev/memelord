package db

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/lib/pq"
	"api/src/core/config"
)



func Connect(config *config.Config) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DB_HOST, config.DB_PORT, config.DB_USER, config.DB_PASSWORD, config.DB_NAME)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("ошибка при подключении к базе данных: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("ошибка при проверке подключения: %v", err)
	}

	log.Println("Успешное подключение к базе данных")
	return db, nil
} 