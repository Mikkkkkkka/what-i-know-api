package repository

import (
	"time"

	"github.com/mikkkkkkka/what-i-know-api/internal/domain"
)

type userModel struct {
	ID        string    `gorm:"primaryKey;type:text"`
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
	ID        string    `gorm:"primaryKey;type:text"`
	UserID    string    `gorm:"type:text;not null;index"`
	Title     string    `gorm:"type:text;not null"`
	Content   string    `gorm:"type:text;not null"`
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
		Title:     note.Title,
		Content:   note.Content,
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
		Title:     model.Title,
		Content:   model.Content,
		UpdatedAt: model.UpdatedAt,
	}
}

type markModel struct {
	ID        string    `gorm:"primaryKey"`
	UserID    string    `gorm:"type:text;not null;index"`
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
