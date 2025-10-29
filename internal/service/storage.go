package service

import (
	"errors"
	"time"

	"server-calendar/internal/entity"
)

var (
	// ErrNotFound возвращается, если событие не найдено.
	ErrNotFound = errors.New("event not found")
	// ErrAlreadyExist возвращается, если событие с таким ID уже существует.
	ErrAlreadyExist = errors.New("event already exists")
)

// Storage — простое хранилище событий в памяти.
type Storage struct {
	data map[int]entity.Event
}

// NewStorage создаёт и возвращает новое пустое хранилище событий.
func NewStorage() *Storage {
	return &Storage{data: make(map[int]entity.Event)}
}

// Create добавляет новое событие в хранилище.
func (s *Storage) Create(e entity.Event) error {
	if _, ok := s.data[e.ID]; ok {
		return ErrAlreadyExist
	}
	s.data[e.ID] = e
	return nil
}

// Update обновляет существующее событие.
func (s *Storage) Update(e entity.Event) error {
	if _, ok := s.data[e.ID]; !ok {
		return ErrNotFound
	}
	s.data[e.ID] = e
	return nil
}

// Delete удаляет событие по ID.
func (s *Storage) Delete(id int) error {
	if _, ok := s.data[id]; !ok {
		return ErrNotFound
	}
	delete(s.data, id)
	return nil
}

// EventsForDay возвращает события пользователя за день.
func (s *Storage) EventsForDay(userID int, date time.Time) []entity.Event {
	return s.filter(userID, date, date)
}

// EventsForWeek возвращает события пользователя за неделю.
func (s *Storage) EventsForWeek(userID int, date time.Time) []entity.Event {
	end := date.AddDate(0, 0, 7)
	return s.filter(userID, date, end)
}

// EventsForMonth возвращает события пользователя за месяц.
func (s *Storage) EventsForMonth(userID int, date time.Time) []entity.Event {
	end := date.AddDate(0, 1, 0)
	return s.filter(userID, date, end)
}

// filter — утилита для выборки событий по диапазону дат.
func (s *Storage) filter(userID int, from, to time.Time) []entity.Event {
	var res []entity.Event
	for _, e := range s.data {
		if e.UserID == userID && (e.Date.Equal(from) || (e.Date.After(from) && e.Date.Before(to))) {
			res = append(res, e)
		}
	}
	return res
}
