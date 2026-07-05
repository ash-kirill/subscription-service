package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/ash-kirill/subscription-service/internal/config"
    "github.com/ash-kirill/subscription-service/internal/handler"
    "github.com/ash-kirill/subscription-service/internal/repository"
)

func main() {
    // 1. Загружаем конфигурацию
    cfg := config.LoadConfig()
    log.Printf("🚀 Starting server on port %s", cfg.ServerPort)

    // 2. Подключаемся к базе данных
    if err := repository.InitDB(cfg); err != nil {
        log.Fatalf("❌ Failed to connect to database: %v", err)
    }
    defer repository.CloseDB()
    log.Println("✅ Database connected successfully!")

    // 3. Создаем обработчики
    subscriptionHandler := handler.NewSubscriptionHandler()

    // 4. Настраиваем маршруты
    // Создаем мультиплексор (роутер)
    mux := http.NewServeMux()

    // Регистрируем маршруты
    mux.HandleFunc("POST /subscriptions", subscriptionHandler.CreateSubscription)
    mux.HandleFunc("GET /subscriptions", subscriptionHandler.ListSubscriptions)
    mux.HandleFunc("GET /subscriptions/{id}", subscriptionHandler.GetSubscription)
    mux.HandleFunc("PUT /subscriptions/{id}", subscriptionHandler.UpdateSubscription)
    mux.HandleFunc("DELETE /subscriptions/{id}", subscriptionHandler.DeleteSubscription)
    mux.HandleFunc("GET /subscriptions/total", subscriptionHandler.GetTotalPrice)

    // 5. Запускаем сервер
    server := &http.Server{
        Addr:    fmt.Sprintf(":%s", cfg.ServerPort),
        Handler: mux,
    }

    log.Printf("✅ Server is running on http://localhost:%s", cfg.ServerPort)
    log.Println("📋 Available endpoints:")
    log.Println("  POST   /subscriptions")
    log.Println("  GET    /subscriptions")
    log.Println("  GET    /subscriptions/{id}")
    log.Println("  PUT    /subscriptions/{id}")
    log.Println("  DELETE /subscriptions/{id}")
    log.Println("  GET    /subscriptions/total?start_date=MM-YYYY&end_date=MM-YYYY")

    if err := server.ListenAndServe(); err != nil {
        log.Fatalf("❌ Failed to start server: %v", err)
    }
}