package service_test

import (
	"context"
	"server-calendar/internal/storage"
	"testing"
	"time"

	"server-calendar/internal/entity"
	"server-calendar/internal/service"
)

func TestCalendarService_CreateEvent(t *testing.T) {
	repo := storage.NewStorage()
	svc := service.NewCalendarService(repo)

	date := time.Now().
		AddDate(0, 0, 1).
		Truncate(24 * time.Hour)

	ev := entity.Event{
		EventID: 1,
		UserID:  42,
		Date:    &date,
	}

	err := svc.CreateEvent(context.Background(), ev)
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	x, err := repo.GetEventsByDateRange(entity.UserID(42), time.Now(), date)
	if len(x) != 1 {
		t.Fatalf("event not saved into repo")
	}
}

func TestCalendarService_EventsForDay(t *testing.T) {
	repo := storage.NewStorage()
	svc := service.NewCalendarService(repo)

	now := time.Now().
		AddDate(0, 0, 1).
		Truncate(24 * time.Hour)

	ev := entity.Event{
		EventID: 3,
		UserID:  1,
		Date:    &now,
	}
	_ = repo.CreateEvent(ev)

	res, err := svc.EventsForDay(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	if len(res) != 1 {
		t.Fatalf("expected 1 event, got %d", len(res))
	}
}
