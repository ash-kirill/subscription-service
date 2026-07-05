package repository

import (
    "database/sql"
    "fmt"
    "time"

    "github.com/google/uuid"
    "github.com/ash-kirill/subscription-service/internal/model"
)

// SubscriptionRepository - структура для работы с подписками в БД
type SubscriptionRepository struct {
    db *sql.DB
}

// NewSubscriptionRepository создает новый репозиторий
func NewSubscriptionRepository() *SubscriptionRepository {
    return &SubscriptionRepository{
        db: DB,
    }
}

// Create создает новую подписку
func (r *SubscriptionRepository) Create(sub *model.Subscription) error {
    query := `
        INSERT INTO subscriptions (id, service_name, price, user_id, start_date, end_date)
        VALUES ($1, $2, $3, $4, $5, $6)
    `

    _, err := r.db.Exec(query,
        sub.ID,
        sub.ServiceName,
        sub.Price,
        sub.UserID,
        sub.StartDate,
        sub.EndDate,
    )

    if err != nil {
        return fmt.Errorf("failed to create subscription: %w", err)
    }

    return nil
}

// GetByID получает подписку по ID
func (r *SubscriptionRepository) GetByID(id uuid.UUID) (*model.Subscription, error) {
    query := `
        SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
        FROM subscriptions
        WHERE id = $1
    `

    var sub model.Subscription
    err := r.db.QueryRow(query, id).Scan(
        &sub.ID,
        &sub.ServiceName,
        &sub.Price,
        &sub.UserID,
        &sub.StartDate,
        &sub.EndDate,
        &sub.CreatedAt,
        &sub.UpdatedAt,
    )

    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("subscription not found")
    }
    if err != nil {
        return nil, fmt.Errorf("failed to get subscription: %w", err)
    }

    return &sub, nil
}

// Update обновляет подписку
func (r *SubscriptionRepository) Update(sub *model.Subscription) error {
    query := `
        UPDATE subscriptions
        SET service_name = $1, price = $2, end_date = $3, updated_at = CURRENT_TIMESTAMP
        WHERE id = $4
    `

    result, err := r.db.Exec(query,
        sub.ServiceName,
        sub.Price,
        sub.EndDate,
        sub.ID,
    )

    if err != nil {
        return fmt.Errorf("failed to update subscription: %w", err)
    }

    rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 0 {
        return fmt.Errorf("subscription not found")
    }

    return nil
}

// Delete удаляет подписку
func (r *SubscriptionRepository) Delete(id uuid.UUID) error {
    query := `DELETE FROM subscriptions WHERE id = $1`

    result, err := r.db.Exec(query, id)
    if err != nil {
        return fmt.Errorf("failed to delete subscription: %w", err)
    }

    rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 0 {
        return fmt.Errorf("subscription not found")
    }

    return nil
}

// List возвращает список всех подписок
func (r *SubscriptionRepository) List() ([]model.Subscription, error) {
    query := `
        SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
        FROM subscriptions
        ORDER BY created_at DESC
    `

    rows, err := r.db.Query(query)
    if err != nil {
        return nil, fmt.Errorf("failed to list subscriptions: %w", err)
    }
    defer rows.Close()

    var subscriptions []model.Subscription
    for rows.Next() {
        var sub model.Subscription
        err := rows.Scan(
            &sub.ID,
            &sub.ServiceName,
            &sub.Price,
            &sub.UserID,
            &sub.StartDate,
            &sub.EndDate,
            &sub.CreatedAt,
            &sub.UpdatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan subscription: %w", err)
        }
        subscriptions = append(subscriptions, sub)
    }

    return subscriptions, nil
}

// GetTotalPriceByPeriod считает сумму за период с фильтрацией
func (r *SubscriptionRepository) GetTotalPriceByPeriod(
    userID *uuid.UUID,
    serviceName *string,
    startDate, endDate time.Time,
) (int, error) {
    query := `
        SELECT COALESCE(SUM(price), 0)
        FROM subscriptions
        WHERE start_date >= $1 
          AND (end_date IS NULL OR end_date <= $2)
    `
    args := []interface{}{startDate, endDate}
    paramIndex := 3

    if userID != nil {
        query += fmt.Sprintf(" AND user_id = $%d", paramIndex)
        args = append(args, *userID)
        paramIndex++
    }

    if serviceName != nil {
        query += fmt.Sprintf(" AND service_name = $%d", paramIndex)
        args = append(args, *serviceName)
    }

    var total int
    err := r.db.QueryRow(query, args...).Scan(&total)
    if err != nil {
        return 0, fmt.Errorf("failed to calculate total: %w", err)
    }

    return total, nil
}