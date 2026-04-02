package domain

import (
	"errors"
	"time"
)

type Mark struct {
	ID        string // Generated on client via UUID v4. Slim (1/2^122) chance of collision!
	UserID    string
	Date      time.Time
	Content   string // Markdown string
	UpdatedAt time.Time
}

var ErrMarkNotFound = errors.New("mark not found")
var ErrMarkAlreadyExists = errors.New("mark already exists")
