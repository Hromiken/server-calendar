package service_test

import (
	"testing"
	"time"

	"server-calendar/internal/entity"
	"server-calendar/internal/service"
)

func TestStorageCRUD(t *testing.T) {
	s := service.NewStorage()
	e := entity.Event{ID: 1, UserID: 42, Title: "Test", Date: time.Now()}

	// Create
	if err := s.Create(e); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Update
	e.Title = "Updated"
	if err := s.Update(e); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// Delete
	if err := s.Delete(1); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}

func TestEventsForDayWeekMonth(t *testing.T) {
	s := service.NewStorage()
	now := time.Now()
	user := 99

	s.Create(entity.Event{ID: 1, UserID: user, Title: "Today", Date: now})
	s.Create(entity.Event{ID: 2, UserID: user, Title: "NextWeek", Date: now.AddDate(0, 0, 7)})
	s.Create(entity.Event{ID: 3, UserID: user, Title: "NextMonth", Date: now.AddDate(0, 1, 0)})

	if len(s.EventsForDay(user, now)) == 0 {
		t.Error("expected at least 1 event for day")
	}
	if len(s.EventsForWeek(user, now)) == 0 {
		t.Error("expected at least 1 event for week")
	}
	if len(s.EventsForMonth(user, now)) == 0 {
		t.Error("expected at least 1 event for month")
	}
}
