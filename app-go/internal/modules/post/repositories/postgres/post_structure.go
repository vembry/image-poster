package postgres

import (
	"app-go/internal/clients"
	"app-go/internal/modules/post/models"
	"context"
	"fmt"
)

type postStructure struct {
	dbProvider clients.IDb
}

func NewPostStructure(dbProvider clients.IDb) *postStructure {
	return &postStructure{dbProvider: dbProvider}
}

func (p *postStructure) Create(ctx context.Context, entry models.PostStructure) error {
	err := p.dbProvider.GetDb().
		WithContext(ctx).
		Table("post_structures").
		Create(&entry).
		Error
	if err != nil {
		return fmt.Errorf("error on creating post structure entry to database. err=%w", err)
	}

	return nil
}

func (p *postStructure) GetMultipleWithCursor(ctx context.Context, limit int, offset int) ([]*models.PostStructure, error) {
	// retrieve row count
	var totalRows int64
	err := p.dbProvider.GetDb().
		WithContext(ctx).
		Table("post_structures").
		Where("parent_post_id IS NULL").
		Count(&totalRows).Error
	if err != nil {
		return nil, fmt.Errorf("error on counting list of post entries. err=%w", err)
	}

	// NOTE:
	// doing this because if offset went beyond actual data count
	// it will return all data, which is bad
	if (offset) >= int(totalRows) {
		return []*models.PostStructure{}, nil
	}

	// retrieve rows
	var out []*models.PostStructure
	err = p.dbProvider.GetDb().
		WithContext(ctx).
		Table("post_structures").
		Offset(offset).
		Limit(limit).
		Find(&out).
		Where("parent_post_id = NULL").
		Error
	if err != nil {
		return out, fmt.Errorf("error on getting list post entry from database. err=%w", err)
	}

	return out, nil
}
