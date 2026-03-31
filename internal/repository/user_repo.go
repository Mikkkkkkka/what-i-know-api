package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/mikkkkkkka/what-i-know-api/internal/domain"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var model userModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		return nil, translateError(err)
	}

	return toDomainUser(&model), nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	var model userModel
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&model).Error; err != nil {
		return nil, translateError(err)
	}

	return toDomainUser(&model), nil
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	model := toUserModel(user)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return translateError(err)
	}

	user.CreatedAt = model.CreatedAt
	return nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	model := toUserModel(user)
	result := r.db.WithContext(ctx).
		Model(&userModel{}).
		Where("id = ?", user.ID).
		Updates(map[string]any{
			"username": model.Username,
			"password": model.Password,
		})
	if result.Error != nil {
		return translateError(result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, userID string) error {
	result := r.db.WithContext(ctx).Delete(&userModel{}, userID)
	if result.Error != nil {
		return translateError(result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

var _ domain.UserRepository = (*UserRepository)(nil)
