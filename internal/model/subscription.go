package model

import (
    "time"
    "github.com/google/uuid"
)

type Subscription struct {
    ID          uuid.UUID  `json:"id"`
    ServiceName string     `json:"service_name"`
    Price       int        `json:"price"`
    UserID      uuid.UUID  `json:"user_id"`
    StartDate   time.Time  `json:"start_date"`
    EndDate     *time.Time `json:"end_date,omitempty"`
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
}

type CreateSubscriptionRequest struct {
    ServiceName string `json:"service_name" binding:"required"`
    Price       int    `json:"price" binding:"required,min=1"`
    UserID      uuid.UUID `json:"user_id" binding:"required"`
    StartDate   string `json:"start_date" binding:"required"` // формат "MM-YYYY"
    EndDate     string `json:"end_date,omitempty"`            // формат "MM-YYYY"
}

type UpdateSubscriptionRequest struct {
    ServiceName string `json:"service_name,omitempty"`
    Price       *int   `json:"price,omitempty"`
    EndDate     string `json:"end_date,omitempty"`
}