package utils

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/mikromolekula2002/Testovoe/internal/config"
)

// Функция для создания базы данных
func CreateDatabase(config *config.Config) {
	// Подключение к PostgreSQL без указания базы данных
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=%s",
		config.Database.Host,
		config.Database.Port,
		config.Database.User,
		config.Database.Password,
		config.Database.Sslmode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Ошибка подключения к PostgreSQL: \n%v", err)
	}
	defer db.Close()

	// Проверяем существует ли база данных и создаем её если нет
	_, err = db.Exec("CREATE DATABASE " + config.Database.DBName)
	if err != nil {
		// Игнорируем ошибку, если база данных уже существует
		if err.Error() != "pq: база данных \""+config.Database.DBName+"\" уже существует" {
			log.Fatalf("Ошибка создания базы данных: \n%v", err)
		} else {
			fmt.Println("База данных уже существует.")
		}
	} else {
		fmt.Println("База данных создана.")
	}

	CreateTable(db)
}

// скорее всего бесполезный метод, оставил на всякий случай
// Создает таблицу для работы кода если ее нет
func CreateTable(db *sql.DB) {
	query := `
    CREATE TABLE IF NOT EXISTS refresh_tokens (
        id SERIAL PRIMARY KEY,
        user_id VARCHAR(255) NOT NULL,
        token_hash VARCHAR(255) NOT NULL UNIQUE,
		blocked boolean NOT NULL DEFAULT false,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		expires_at TIMESTAMP NOT NULL,
        ip_adress VARCHAR(255) NOT NULL
    );`
	_, err := db.Exec(query)
	if err != nil {
		fmt.Println("Ошибка создания таблицы:", err)
	} else {
		fmt.Println("Таблица refresh_tokens создана.")
	}
}
