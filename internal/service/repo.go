package service

import (
	"server-calendar/internal/entity"
	"time"
)

type Repository interface {
	Create(e entity.Event) error
	Update(e entity.Event) error
	Delete(id int) error
	EventsForDay(userID int, date time.Time) []entity.Event
	EventsForWeek(userID int, date time.Time) []entity.Event
	EventsForMonth(userID int, date time.Time) []entity.Event
	Filter(userID int, from, to time.Time) []entity.Event
}
