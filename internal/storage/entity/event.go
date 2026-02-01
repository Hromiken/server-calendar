package entity

import "time"

// UserID - псевдоним для user ключа.
type UserID int

// EventID - псевдоним для event ключа
type EventID int

// Event описывает событие календаря.
type Event struct {
	EventID EventID    `json:"event_id" validate:"required"`
	UserID  UserID     `json:"user_id" validate:"required"`
	Date    *time.Time `json:"date" validate:"omitempty"`
	Title   *string    `json:"title" validate:"omitempty"`
}
