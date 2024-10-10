package repository

import (
	"context"

	"gorm.io/gorm"
)

type WithTransaction interface {
	WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error
}
