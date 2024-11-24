package repositories

import (
	"app-go/internal/modules/post/models"
	"context"
)

type IPost interface {
	Create(ctx context.Context, entry models.Post) error
}
