package repository

import (
	"context"

	"gorm.io/gorm"

	"what-i-know-api/internal/domain"
)

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) GetByID(ctx context.Context, id int64) (*domain.Session, error) {
	var model sessionModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		return nil, translateError(err)
	}

	return toDomainSession(&model), nil
}

func (r *SessionRepository) GetByToken(ctx context.Context, token string) (*domain.Session, error) {
	var model sessionModel
	if err := r.db.WithContext(ctx).Where("token = ?", token).First(&model).Error; err != nil {
		return nil, translateError(err)
	}

	return toDomainSession(&model), nil
}

func (r *SessionRepository) Create(ctx context.Context, session *domain.Session) error {
	model := toSessionModel(session)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return translateError(err)
	}

	session.Id = model.ID
	return nil
}

func (r *SessionRepository) Delete(ctx context.Context, session *domain.Session) error {
	result := r.db.WithContext(ctx).Delete(&sessionModel{}, session.Id)
	if result.Error != nil {
		return translateError(result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

var _ domain.SessionRepository = (*SessionRepository)(nil)
