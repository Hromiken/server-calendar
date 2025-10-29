package service

import (
	"errors"
	"sync"
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
	mu   sync.RWMutex
	data map[int]entity.Event
}

// NewStorage создаёт и возвращает новое пустое хранилище событий.
func NewStorage() *Storage {
	return &Storage{
		data: make(map[int]entity.Event),
	}
}

// Create добавляет новое событие в хранилище.
func (s *Storage) Create(e entity.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.data[e.ID]
	if ok {
		return ErrAlreadyExist
	}
	s.data[e.ID] = e
	return nil
}

// Update обновляет существующее событие.
func (s *Storage) Update(e entity.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.data[e.ID]
	if !ok {
		return ErrNotFound
	}
	s.data[e.ID] = e
	return nil
}

// Delete удаляет событие по ID.
func (s *Storage) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.data[id]
	if !ok {
		return ErrNotFound
	}
	delete(s.data, id)
	return nil
}

// EventsForDay возвращает события пользователя за день.
func (s *Storage) EventsForDay(userID int, date time.Time) []entity.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Filter(userID, date, date)
}

// EventsForWeek возвращает события пользователя за неделю.
func (s *Storage) EventsForWeek(userID int, date time.Time) []entity.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()
	end := date.AddDate(0, 0, 7)
	return s.Filter(userID, date, end)
}

// EventsForMonth возвращает события пользователя за месяц.
func (s *Storage) EventsForMonth(userID int, date time.Time) []entity.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()
	end := date.AddDate(0, 1, 0)
	return s.Filter(userID, date, end)
}

// Filter утилита для выборки событий по диапазону дат
func (s *Storage) Filter(userID int, from, to time.Time) []entity.Event {
	var res []entity.Event
	for _, e := range s.data {
		if e.UserID == userID && !e.Date.Before(from) && !e.Date.After(to) {
			res = append(res, e)
		}
	}
	return res
}
