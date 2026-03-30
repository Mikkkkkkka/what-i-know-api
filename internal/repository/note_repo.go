package repository

import (
	"context"

	"gorm.io/gorm"

	"what-i-know-api/internal/domain"
)

type NoteRepository struct {
	db *gorm.DB
}

func NewNoteRepository(db *gorm.DB) *NoteRepository {
	return &NoteRepository{db: db}
}

func (r *NoteRepository) GetById(ctx context.Context, id string) (*domain.Note, error) {
	var model noteModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		return nil, translateError(err)
	}

	return toDomainNote(&model), nil
}

func (r *NoteRepository) GetByUserId(ctx context.Context, userId string) ([]*domain.Note, error) {
	var models []noteModel
	if err := r.db.WithContext(ctx).Where("user_id = ?", userId).Order("date DESC, id DESC").Find(&models).Error; err != nil {
		return nil, translateError(err)
	}

	notes := make([]*domain.Note, 0, len(models))
	for i := range models {
		notes = append(notes, toDomainNote(&models[i]))
	}

	return notes, nil
}

func (r *NoteRepository) Create(ctx context.Context, note *domain.Note) error {
	model := toNoteModel(note)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return translateError(err)
	}

	note.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *NoteRepository) Update(ctx context.Context, note *domain.Note) error {
	model := toNoteModel(note)
	result := r.db.WithContext(ctx).
		Model(&noteModel{}).
		Where("id = ?", note.Id).
		Updates(map[string]any{
			"user_id":    model.UserID,
			"name":       model.Title,
			"content":    model.Content,
			"updated_at": model.UpdatedAt,
		})
	if result.Error != nil {
		return translateError(result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}

	var updated noteModel
	if err := r.db.WithContext(ctx).First(&updated, note.Id).Error; err != nil {
		return translateError(err)
	}
	note.UpdatedAt = updated.UpdatedAt
	return nil
}

func (r *NoteRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&noteModel{}, id)
	if result.Error != nil {
		return translateError(result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

var _ domain.NoteRepository = (*NoteRepository)(nil)
