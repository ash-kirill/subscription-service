package service

import (
    "fmt"
    "time"

    "github.com/google/uuid"
    "github.com/ash-kirill/subscription-service/internal/model"
    "github.com/ash-kirill/subscription-service/internal/repository"
)

type SubscriptionService struct {
    repo *repository.SubscriptionRepository
}

func NewSubscriptionService() *SubscriptionService {
    return &SubscriptionService{
        repo: repository.NewSubscriptionRepository(),
    }
}

// Create создает новую подписку
func (s *SubscriptionService) Create(req *model.CreateSubscriptionRequest) (*model.Subscription, error) {
    // Парсим дату начала
    startDate, err := time.Parse("01-2006", req.StartDate)
    if err != nil {
        return nil, fmt.Errorf("invalid start_date format, expected MM-YYYY: %w", err)
    }

    var endDate *time.Time
    if req.EndDate != "" {
        parsed, err := time.Parse("01-2006", req.EndDate)
        if err != nil {
            return nil, fmt.Errorf("invalid end_date format, expected MM-YYYY: %w", err)
        }
        endDate = &parsed

        // Проверяем, что дата окончания не раньше даты начала
        if endDate.Before(startDate) {
            return nil, fmt.Errorf("end_date cannot be before start_date")
        }
    }

    sub := &model.Subscription{
        ID:          uuid.New(),
        ServiceName: req.ServiceName,
        Price:       req.Price,
        UserID:      req.UserID,
        StartDate:   startDate,
        EndDate:     endDate,
    }

    err = s.repo.Create(sub)
    if err != nil {
        return nil, err
    }

    return sub, nil
}

// GetByID получает подписку по ID
func (s *SubscriptionService) GetByID(id string) (*model.Subscription, error) {
    uuid, err := uuid.Parse(id)
    if err != nil {
        return nil, fmt.Errorf("invalid UUID format: %w", err)
    }

    return s.repo.GetByID(uuid)
}

// Update обновляет подписку
func (s *SubscriptionService) Update(id string, req *model.UpdateSubscriptionRequest) (*model.Subscription, error) {
    uuid, err := uuid.Parse(id)
    if err != nil {
        return nil, fmt.Errorf("invalid UUID format: %w", err)
    }

    // Получаем существующую подписку
    existing, err := s.repo.GetByID(uuid)
    if err != nil {
        return nil, err
    }

    // Обновляем только те поля, которые переданы
    if req.ServiceName != "" {
        existing.ServiceName = req.ServiceName
    }
    if req.Price != nil {
        existing.Price = *req.Price
    }
    if req.EndDate != "" {
        parsed, err := time.Parse("01-2006", req.EndDate)
        if err != nil {
            return nil, fmt.Errorf("invalid end_date format, expected MM-YYYY: %w", err)
        }
        existing.EndDate = &parsed

        if existing.EndDate.Before(existing.StartDate) {
            return nil, fmt.Errorf("end_date cannot be before start_date")
        }
    }

    err = s.repo.Update(existing)
    if err != nil {
        return nil, err
    }

    return existing, nil
}

// Delete удаляет подписку
func (s *SubscriptionService) Delete(id string) error {
    uuid, err := uuid.Parse(id)
    if err != nil {
        return fmt.Errorf("invalid UUID format: %w", err)
    }

    return s.repo.Delete(uuid)
}

// List возвращает все подписки
func (s *SubscriptionService) List() ([]model.Subscription, error) {
    return s.repo.List()
}

// GetTotalPriceByPeriod считает сумму за период
func (s *SubscriptionService) GetTotalPriceByPeriod(
    userIDStr, serviceName, startDateStr, endDateStr string,
) (int, error) {
    // Парсим даты
    startDate, err := time.Parse("01-2006", startDateStr)
    if err != nil {
        return 0, fmt.Errorf("invalid start_date format, expected MM-YYYY: %w", err)
    }

    endDate, err := time.Parse("01-2006", endDateStr)
    if err != nil {
        return 0, fmt.Errorf("invalid end_date format, expected MM-YYYY: %w", err)
    }

    if endDate.Before(startDate) {
        return 0, fmt.Errorf("end_date cannot be before start_date")
    }

    var userID *uuid.UUID
    if userIDStr != "" {
        parsed, err := uuid.Parse(userIDStr)
        if err != nil {
            return 0, fmt.Errorf("invalid user_id format: %w", err)
        }
        userID = &parsed
    }

    var serviceNamePtr *string
    if serviceName != "" {
        serviceNamePtr = &serviceName
    }

    return s.repo.GetTotalPriceByPeriod(userID, serviceNamePtr, startDate, endDate)
}