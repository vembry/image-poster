package postgres

import (
	"app-go/internal/clients"
	"app-go/internal/modules/post/models"
	"context"
	"fmt"
)

type post struct {
	dbProvider clients.IDb
}

func New(dbProvider clients.IDb) *post {
	return &post{dbProvider: dbProvider}
}

// Create create post entry
func (p *post) Create(ctx context.Context, entry models.Post) error {
	err := p.dbProvider.GetDb().
		WithContext(ctx).
		Table("posts").
		Create(&entry).
		Error
	if err != nil {
		return fmt.Errorf("error on creating post entry to database. err=%v", err)
	}

	return nil
}
