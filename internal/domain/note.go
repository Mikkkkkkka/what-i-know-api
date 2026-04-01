package domain

import (
	"errors"
	"time"
)

type Note struct {
	ID        string // Generated on client via UUID v4. Slim (1/2^122) chance of collision!
	UserID    string
	Title     string
	Content   string // Markdown string
	UpdatedAt time.Time
}

var ErrNoteNotFound = errors.New("note not found")
