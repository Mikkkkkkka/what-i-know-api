package repository

import (
	"time"

	"what-i-know-api/internal/domain"
)

type userModel struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	Username  string    `gorm:"type:text;not null;uniqueIndex"`
	Password  string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime"`
}

func (userModel) TableName() string {
	return "users"
}

func toUserModel(user *domain.User) *userModel {
	if user == nil {
		return nil
	}

	return &userModel{
		ID:        user.Id,
		Username:  user.Username,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
	}
}

func toDomainUser(model *userModel) *domain.User {
	if model == nil {
		return nil
	}

	return &domain.User{
		Id:        model.ID,
		Username:  model.Username,
		Password:  model.Password,
		CreatedAt: model.CreatedAt,
	}
}

type noteModel struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	UserID    int64     `gorm:"not null;index"`
	Name      string    `gorm:"type:text;not null"`
	Content   string    `gorm:"type:text;not null"`
	Date      time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null;autoUpdateTime"`
}

func (noteModel) TableName() string {
	return "notes"
}

func toNoteModel(note *domain.Note) *noteModel {
	if note == nil {
		return nil
	}

	return &noteModel{
		ID:        note.Id,
		UserID:    note.UserId,
		Name:      note.Name,
		Content:   note.Content,
		Date:      note.Date,
		UpdatedAt: note.UpdatedAt,
	}
}

func toDomainNote(model *noteModel) *domain.Note {
	if model == nil {
		return nil
	}

	return &domain.Note{
		Id:        model.ID,
		UserId:    model.UserID,
		Name:      model.Name,
		Content:   model.Content,
		Date:      model.Date,
		UpdatedAt: model.UpdatedAt,
	}
}

type markModel struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	UserID    int64     `gorm:"not null;index"`
	Date      time.Time `gorm:"not null"`
	Content   string    `gorm:"type:text;not null"`
	UpdatedAt time.Time `gorm:"not null;autoUpdateTime"`
}

func (markModel) TableName() string {
	return "marks"
}

func toMarkModel(mark *domain.Mark) *markModel {
	if mark == nil {
		return nil
	}

	return &markModel{
		ID:        mark.Id,
		UserID:    mark.UserId,
		Date:      mark.Date,
		Content:   mark.Content,
		UpdatedAt: mark.UpdatedAt,
	}
}

func toDomainMark(model *markModel) *domain.Mark {
	if model == nil {
		return nil
	}

	return &domain.Mark{
		Id:        model.ID,
		UserId:    model.UserID,
		Date:      model.Date,
		Content:   model.Content,
		UpdatedAt: model.UpdatedAt,
	}
}

type sessionModel struct {
	ID       int64  `gorm:"primaryKey;autoIncrement"`
	UserID   int64  `gorm:"not null;index"`
	Username string `gorm:"type:text;not null"`
	Token    string `gorm:"type:text;not null;uniqueIndex"`
}

func (sessionModel) TableName() string {
	return "sessions"
}

func toSessionModel(session *domain.Session) *sessionModel {
	if session == nil {
		return nil
	}

	return &sessionModel{
		ID:       session.Id,
		UserID:   session.UserId,
		Username: session.Username,
		Token:    session.Token,
	}
}

func toDomainSession(model *sessionModel) *domain.Session {
	if model == nil {
		return nil
	}

	return &domain.Session{
		Id:       model.ID,
		UserId:   model.UserID,
		Username: model.Username,
		Token:    model.Token,
	}
}
