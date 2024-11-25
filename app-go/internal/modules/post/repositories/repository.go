package repositories

import (
	"app-go/internal/modules/post/models"
	"context"

	"github.com/segmentio/ksuid"
)

type IPost interface {
	// Create creates post entry
	Create(ctx context.Context, entry *models.Post) error

	// RollbackCreate rollbacks post entry creation
	RollbackCreate(ctx context.Context, entry *models.Post) error

	GetByMultipleIds(ctx context.Context, ids []string) ([]*models.Post, error)
	GetById(ctx context.Context, id ksuid.KSUID) (*models.Post, error)

	// Update updates existing post entry
	Update(ctx context.Context, entry *models.Post) error
}

type IPostStructure interface {
	// Create create post structure entry
	Create(ctx context.Context, entry models.PostStructure) error

	// RollbackCreate rollbacks post structure entry creation
	RollbackCreate(ctx context.Context, entry models.PostStructure) error

	// GetMultipleWithCursor get list of posts with cursor
	GetMultipleWithCursor(ctx context.Context, limit int, offset int) ([]*models.PostStructure, error)
}
