package service

import (
	"context"
	"errors"
	"time"

	"server-calendar/internal/entity"
)

type CalendarService struct {
	repo repository
}

type repository interface {
	CreateEvent(event entity.Event) error
	UpdateEvent(event entity.Event) error
	DeleteEvent(event entity.Event) error
	GetEventsByDateRange(userID entity.UserID, from, to time.Time) ([]entity.Event, error)
}

func NewCalendarService(repo repository) *CalendarService {
	return &CalendarService{repo: repo}
}

func (s *CalendarService) CreateEvent(_ context.Context, e entity.Event) error {
	if e.Date.Before(time.Now()) {
		return errors.New("cannot create event in the past")
	}
	return s.repo.CreateEvent(e)
}

func (s *CalendarService) UpdateEvent(_ context.Context, e entity.Event) error {
	if e.Date.Before(time.Now()) {
		return errors.New("cannot update event in the past")
	}
	return s.repo.UpdateEvent(e)
}

func (s *CalendarService) DeleteEvent(_ context.Context, id entity.Event) error {
	return s.repo.DeleteEvent(id)
}

func (s *CalendarService) EventsForDay(_ context.Context, userID entity.UserID) ([]entity.Event, error) {
	from := time.Now()
	to := from.
		AddDate(0, 0, 1).
		Truncate(24 * time.Hour)

	return s.repo.GetEventsByDateRange(userID, from, to)
}

func (s *CalendarService) EventsForWeek(_ context.Context, userID entity.UserID) ([]entity.Event, error) {
	from := time.Now()
	to := from.
		AddDate(0, 0, 7).
		Truncate(24 * time.Hour)
	return s.repo.GetEventsByDateRange(userID, from, to)
}

func (s *CalendarService) EventsForMonth(_ context.Context, userID entity.UserID) ([]entity.Event, error) {
	from := time.Now()
	to := from.
		AddDate(0, 1, 0).
		Truncate(24 * time.Hour)
	return s.repo.GetEventsByDateRange(userID, from, to)
}

// 1) шардирование ,2) shutdown, 3) map в map
