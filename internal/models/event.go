package models

import "time"

// Event represents event object
type Event struct {
	ID      int       `db:"id"`
	When    time.Time `db:"when"`
	Message string    `db:"message"`
	Done    bool      `db:"done"`
}
