package postgres

import (
	"app-go/internal/modules/post/models"
	"context"

	"gorm.io/gorm"
)

type post struct {
	db *gorm.DB
}

func New(db *gorm.DB) *post {
	return &post{db: db}
}

func (p *post) CreatePost(ctx context.Context, entry models.Post) error {
	return nil
}
