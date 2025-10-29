package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"server-calendar/internal/entity"
	"server-calendar/internal/service"
)

type EventHandler struct {
	storage *service.Storage
}

// NewEventHandler создаёт обработчик с бизнес-логикой
func NewEventHandler(storage *service.Storage) *EventHandler {
	return &EventHandler{storage: storage}
}

// CreateEvent — создание нового события
func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var e entity.Event
	errDecoder := json.NewDecoder(r.Body).Decode(&e)
	if errDecoder != nil {
		writeError(w, "Json invalid", http.StatusBadRequest)
		return
	}

	err := h.storage.Create(e)
	if err != nil {
		writeError(w, "Create error", http.StatusBadRequest)
		return
	}

	writeJSON(w, map[string]string{"result": "created"})
}

// UpdateEvent — обновление события
func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	var e entity.Event
	errDecoder := json.NewDecoder(r.Body).Decode(&e)
	if errDecoder != nil {
		writeError(w, "Json invalid", http.StatusBadRequest)
		return
	}

	err := h.storage.Update(e)
	if err != nil {
		writeError(w, "Update error", http.StatusBadRequest)
		return
	}

	writeJSON(w, map[string]string{"result": "updated"})
}

// DeleteEvent — удаление события по id
func (h *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)

	err := h.storage.Delete(id)
	if err != nil {
		writeError(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	writeJSON(w, map[string]string{"result": "deleted"})
}

// EventsForDay — события за день
func (h *EventHandler) EventsForDay(w http.ResponseWriter, r *http.Request) {
	h.eventsForRange(w, r, "day")
}

// EventsForWeek — события за неделю
func (h *EventHandler) EventsForWeek(w http.ResponseWriter, r *http.Request) {
	h.eventsForRange(w, r, "week")
}

// EventsForMonth — события за месяц
func (h *EventHandler) EventsForMonth(w http.ResponseWriter, r *http.Request) {
	h.eventsForRange(w, r, "month")
}

// универсальный метод для выборки по диапазону
func (h *EventHandler) eventsForRange(w http.ResponseWriter, r *http.Request, period string) {
	userID, errorAtoi := strconv.Atoi(r.URL.Query().Get("user_id"))
	if errorAtoi != nil {
		writeError(w, errorAtoi.Error(), http.StatusBadRequest)
		return
	}
	dateStr := r.URL.Query().Get("date")

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		writeError(w, "invalid date", http.StatusBadRequest)
		return
	}

	var events []entity.Event

	switch period {
	case "day":
		events = h.storage.EventsForDay(userID, date)
	case "week":
		events = h.storage.EventsForWeek(userID, date)
	case "month":
		events = h.storage.EventsForMonth(userID, date)
	}

	writeJSON(w, events)
}

// writeJSON — утилита для возврата JSON-ответа
func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
