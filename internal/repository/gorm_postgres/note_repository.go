package gorm_postgres

import (
	"context"

	"github.com/mikkkkkkka/what-i-know-api/internal/service"
	"gorm.io/gorm"

	"github.com/mikkkkkkka/what-i-know-api/internal/domain"
)

type NoteRepository struct {
	db *gorm.DB
}

func NewNoteRepository(db *gorm.DB) *NoteRepository {
	return &NoteRepository{db: db}
}

func (r *NoteRepository) GetByID(ctx context.Context, id string) (*domain.Note, error) {
	var model noteModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return nil, err
	}

	return toDomainNote(&model), nil
}

func (r *NoteRepository) GetByUserID(ctx context.Context, userID string) ([]*domain.Note, error) {
	var models []noteModel
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("title").Find(&models).Error; err != nil {
		return nil, err
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
		return err
	}

	note.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *NoteRepository) Update(ctx context.Context, note *domain.Note) error {
	model := toNoteModel(note)
	result := r.db.WithContext(ctx).
		Model(&noteModel{}).
		Where("id = ?", note.ID).
		Updates(map[string]any{
			"user_id":    model.UserID,
			"title":      model.Title,
			"content":    model.Content,
			"updated_at": model.UpdatedAt,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrNoteNotFound
	}

	var updated noteModel
	if err := r.db.WithContext(ctx).First(&updated, note.ID).Error; err != nil {
		return err
	}
	note.UpdatedAt = updated.UpdatedAt
	return nil
}

func (r *NoteRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&noteModel{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrNoteNotFound
	}

	return nil
}

var _ service.NoteRepository = (*NoteRepository)(nil)
