package entity

import "time"

// Event описывает событие календаря.
type Event struct {
	ID     int
	UserID int
	Date   time.Time
	Title  string
}
