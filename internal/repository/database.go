package repository

import (
    "database/sql"
    "fmt"
    "log"

    _"github.com/lib/pq" // PostgreSQL драйвер(подчеркивание - импорт для побочного эффекта)
    "github.com/ash-kirill/subscription-service/internal/config"
)

var DB *sql.DB

// InitDB инициализирует подключение к базе данных
func InitDB(cfg *config.Config) error {
    // Формируем строку подключения
    // Пример: "host=localhost port=5432 user=postgres password=postgres dbname=subscription_db sslmode=disable"
    connStr := fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
    )

    var err error
    DB, err = sql.Open("postgres", connStr)
    if err != nil {
        return fmt.Errorf("failed to open database: %w", err)
    }

    // Проверяем подключение
    err = DB.Ping()
    if err != nil {
        return fmt.Errorf("failed to ping database: %w", err)
    }

    log.Println("Database connected successfully!")
    return nil
}

// CloseDB закрывает подключение к базе данных
func CloseDB() error {
    if DB != nil {
        return DB.Close()
    }
    return nil
}