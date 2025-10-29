package service_test

import (
	"errors"
	"testing"
	"time"

	"server-calendar/internal/entity"
	"server-calendar/internal/service"
)

func TestStorage_CreateUpdateDelete(t *testing.T) {
	s := service.NewStorage()
	event := entity.Event{
		ID:     1,
		UserID: 42,
		Title:  "Test Event",
		Date:   time.Now(),
	}

	//Create
	err := s.Create(event)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Попробуем создать ещё раз — должно вернуть ошибку ErrAlreadyExist
	errCreate := s.Create(event)
	if !errors.Is(errCreate, service.ErrAlreadyExist) {
		t.Errorf("expected ErrAlreadyExist, got %v", err)
	}

	//Update
	event.Title = "Updated Title"
	errUpdate := s.Update(event)
	if errUpdate != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// Попробуем обновить несуществующий ID
	errUpdateExist := s.Update(entity.Event{ID: 999})
	if !errors.Is(errUpdateExist, service.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}

	//Delete
	errDelete := s.Delete(1)
	if errDelete != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	errSecondDelete := s.Delete(1)
	if !errors.Is(errSecondDelete, service.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestStorage_EventsForPeriods(t *testing.T) {
	s := service.NewStorage()
	now := time.Now().Truncate(24 * time.Hour)
	user := 99

	s.Create(entity.Event{ID: 1, UserID: user, Title: "Today", Date: now})
	s.Create(entity.Event{ID: 2, UserID: user, Title: "NextWeek", Date: now.AddDate(0, 0, 7)})
	s.Create(entity.Event{ID: 3, UserID: user, Title: "NextMonth", Date: now.AddDate(0, 1, 0)})

	if got := s.EventsForDay(user, now); len(got) != 1 {
		t.Errorf("expected 1 event for day, got %d", len(got))
	}
	if got := s.EventsForWeek(user, now); len(got) != 2 {
		t.Errorf("expected 2 events for week, got %d", len(got))
	}
	if got := s.EventsForMonth(user, now); len(got) != 3 {
		t.Errorf("expected 3 events for month, got %d", len(got))
	}
}

func TestStorage_ConcurrentAccess(t *testing.T) {
	s := service.NewStorage()
	user := 42
	now := time.Now().Truncate(24 * time.Hour)

	s.Create(entity.Event{ID: 1, UserID: user, Title: "Initial", Date: now})

	done := make(chan struct{})
	const goroutines = 20

	for i := 0; i < goroutines; i++ {
		go func(i int) {
			defer func() { done <- struct{}{} }()

			e := entity.Event{ID: i + 2, UserID: user, Title: "Event", Date: now.AddDate(0, 0, i)}

			s.Create(e)
			s.Update(e)
			s.EventsForDay(user, now)
			s.EventsForWeek(user, now)
			s.Delete(e.ID)
		}(i)
	}

	for i := 0; i < goroutines; i++ {
		<-done
	}

}
