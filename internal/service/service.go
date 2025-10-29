package service

import (
	"errors"
	"time"

	"server-calendar/internal/entity"
)

type CalendarService struct {
	repo Repository
}

func NewCalendarService(repo Repository) *CalendarService {
	return &CalendarService{repo: repo}
}

func (s *CalendarService) CreateEvent(e entity.Event) error {
	if e.Date.Before(time.Now()) {
		return errors.New("cannot create event in the past")
	}
	return s.repo.Create(e)
}

func (s *CalendarService) UpdateEvent(e entity.Event) error {
	return s.repo.Update(e)
}

func (s *CalendarService) DeleteEvent(id int) error {
	return s.repo.Delete(id)
}

func (s *CalendarService) EventsForDay(userID int, date time.Time) []entity.Event {
	return s.repo.Filter(userID, date, date)
}

func (s *CalendarService) EventsForWeek(userID int, date time.Time) []entity.Event {
	end := date.AddDate(0, 0, 7)
	return s.repo.Filter(userID, date, end)
}

func (s *CalendarService) EventsForMonth(userID int, date time.Time) []entity.Event {
	end := date.AddDate(0, 1, 0)
	return s.repo.Filter(userID, date, end)
}
