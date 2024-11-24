package postgres

import (
	"app-go/internal/clients"
	"app-go/internal/modules/post/models"
	"context"
	"fmt"

	"github.com/segmentio/ksuid"
)

type post struct {
	dbProvider clients.IDb
}

func NewPost(dbProvider clients.IDb) *post {
	return &post{dbProvider: dbProvider}
}

func (p *post) Create(ctx context.Context, entry *models.Post) error {
	entry.Id = ksuid.New()

	// insert entry
	err := p.dbProvider.GetDb().
		WithContext(ctx).
		Table("posts").
		Create(&entry).
		Error
	if err != nil {
		return fmt.Errorf("error on creating post entry to database. err=%w", err)
	}

	return nil
}

func (p *post) RollbackCreate(ctx context.Context, entry *models.Post) error {
	// delete entry
	err := p.dbProvider.GetDb().
		WithContext(ctx).
		Table("posts").
		Delete(&entry).
		Error
	if err != nil {
		return fmt.Errorf("error on rolling back post creation to database. err=%w", err)
	}

	return nil
}

func (p *post) GetByMultipleIds(ctx context.Context, ids []string) ([]*models.Post, error) {
	var out []*models.Post

	err := p.dbProvider.GetDb().
		WithContext(ctx).
		Table("posts").
		Find(&out, ids).
		Error
	if err != nil {
		return out, fmt.Errorf("error on getting list of posts by multiple ids from database. err=%w", err)
	}

	return out, nil

}
