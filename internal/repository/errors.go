package repository

import (
	"errors"

	"gorm.io/gorm"

	"what-i-know-api/internal/domain"
)

func translateError(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		return domain.ErrNotFound
	case errors.Is(err, gorm.ErrDuplicatedKey):
		return domain.ErrAlreadyExists
	default:
		return err
	}
}
