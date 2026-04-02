package domain

import (
	"errors"
	"time"
)

type User struct {
	ID        string
	Username  string
	Password  string
	CreatedAt time.Time
}

var ErrUserNotFound = errors.New("user not found")
var ErrUsernameAlreadyExists = errors.New("username already exists")
var ErrIncorrectCredentials = errors.New("incorrect credentials")
