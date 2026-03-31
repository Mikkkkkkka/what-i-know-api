package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/mikkkkkkka/what-i-know-api/internal/domain"
)

type MarkRepository struct {
	db *gorm.DB
}

func NewMarkRepository(db *gorm.DB) *MarkRepository {
	return &MarkRepository{db: db}
}

func (r *MarkRepository) GetByID(ctx context.Context, id string) (*domain.Mark, error) {
	var model markModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		return nil, translateError(err)
	}

	return toDomainMark(&model), nil
}

func (r *MarkRepository) GetByUserID(ctx context.Context, userID string) ([]*domain.Mark, error) {
	var models []markModel
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("date DESC").Find(&models).Error; err != nil {
		return nil, translateError(err)
	}

	marks := make([]*domain.Mark, 0, len(models))
	for i := range models {
		marks = append(marks, toDomainMark(&models[i]))
	}

	return marks, nil
}

func (r *MarkRepository) Create(ctx context.Context, mark *domain.Mark) error {
	model := toMarkModel(mark)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return translateError(err)
	}

	mark.ID = model.ID
	mark.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *MarkRepository) Update(ctx context.Context, mark *domain.Mark) error {
	model := toMarkModel(mark)
	result := r.db.WithContext(ctx).
		Model(&markModel{}).
		Where("id = ?", mark.ID).
		Updates(map[string]any{
			"user_id":    model.UserID,
			"date":       model.Date,
			"content":    model.Content,
			"updated_at": model.UpdatedAt,
		})
	if result.Error != nil {
		return translateError(result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}

	var updated markModel
	if err := r.db.WithContext(ctx).First(&updated, mark.ID).Error; err != nil {
		return translateError(err)
	}
	mark.UpdatedAt = updated.UpdatedAt
	return nil
}

func (r *MarkRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&markModel{}, id)
	if result.Error != nil {
		return translateError(result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

var _ domain.MarkRepository = (*MarkRepository)(nil)
