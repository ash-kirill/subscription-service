package handler

import (
    "encoding/json"
    "net/http"
    //"strconv"

    "github.com/ash-kirill/subscription-service/internal/model"
    "github.com/ash-kirill/subscription-service/internal/service"
)

type SubscriptionHandler struct {
    service *service.SubscriptionService
}

func NewSubscriptionHandler() *SubscriptionHandler {
    return &SubscriptionHandler{
        service: service.NewSubscriptionService(),
    }
}

// CreateSubscription создает новую подписку
// @Summary Создать подписку
// @Description Создает новую запись о подписке
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param request body model.CreateSubscriptionRequest true "Данные подписки"
// @Success 201 {object} model.Subscription
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions [post]
func (h *SubscriptionHandler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
    var req model.CreateSubscriptionRequest
    
    // Декодируем JSON из запроса
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
        return
    }

    // Вызываем сервис для создания подписки
    subscription, err := h.service.Create(&req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Отправляем ответ с кодом 201 Created
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(subscription)
}

// GetSubscription получает подписку по ID
// @Summary Получить подписку
// @Description Получает подписку по ее ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "ID подписки (UUID)"
// @Success 200 {object} model.Subscription
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetSubscription(w http.ResponseWriter, r *http.Request) {
    // Получаем ID из URL: /subscriptions/123e4567-e89b-12d3-a456-426614174000
    id := r.PathValue("id")
    if id == "" {
        http.Error(w, "ID is required", http.StatusBadRequest)
        return
    }

    subscription, err := h.service.GetByID(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(subscription)
}

// UpdateSubscription обновляет подписку
// @Summary Обновить подписку
// @Description Обновляет существующую подписку
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "ID подписки (UUID)"
// @Param request body model.UpdateSubscriptionRequest true "Данные для обновления"
// @Success 200 {object} model.Subscription
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [put]
func (h *SubscriptionHandler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    if id == "" {
        http.Error(w, "ID is required", http.StatusBadRequest)
        return
    }

    var req model.UpdateSubscriptionRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
        return
    }

    subscription, err := h.service.Update(id, &req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    json.NewEncoder(w).Encode(subscription)
}

// DeleteSubscription удаляет подписку
// @Summary Удалить подписку
// @Description Удаляет подписку по ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "ID подписки (UUID)"
// @Success 204 "No Content"
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    if id == "" {
        http.Error(w, "ID is required", http.StatusBadRequest)
        return
    }

    err := h.service.Delete(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

// ListSubscriptions возвращает список всех подписок
// @Summary Список подписок
// @Description Возвращает список всех подписок
// @Tags subscriptions
// @Produce json
// @Success 200 {array} model.Subscription
// @Failure 500 {object} map[string]string
// @Router /subscriptions [get]
func (h *SubscriptionHandler) ListSubscriptions(w http.ResponseWriter, r *http.Request) {
    subscriptions, err := h.service.List()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(subscriptions)
}

// GetTotalPrice считает сумму подписок за период
// @Summary Сумма подписок за период
// @Description Считает суммарную стоимость подписок за выбранный период
// @Tags subscriptions
// @Produce json
// @Param start_date query string true "Дата начала (MM-YYYY)"
// @Param end_date query string true "Дата окончания (MM-YYYY)"
// @Param user_id query string false "ID пользователя (UUID)"
// @Param service_name query string false "Название сервиса"
// @Success 200 {object} map[string]int
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/total [get]
func (h *SubscriptionHandler) GetTotalPrice(w http.ResponseWriter, r *http.Request) {
    // Получаем параметры из query-строки: /subscriptions/total?start_date=07-2025&end_date=12-2025
    startDate := r.URL.Query().Get("start_date")
    endDate := r.URL.Query().Get("end_date")
    userID := r.URL.Query().Get("user_id")
    serviceName := r.URL.Query().Get("service_name")

    // Проверяем обязательные параметры
    if startDate == "" || endDate == "" {
        http.Error(w, "start_date and end_date are required", http.StatusBadRequest)
        return
    }

    total, err := h.service.GetTotalPriceByPeriod(userID, serviceName, startDate, endDate)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    json.NewEncoder(w).Encode(map[string]int{"total": total})
}