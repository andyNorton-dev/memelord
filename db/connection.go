package db

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "root"
	password = "root"
	dbname   = "gobd"
)

func Connect() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

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