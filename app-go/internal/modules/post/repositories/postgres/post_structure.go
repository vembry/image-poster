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

func (p *postStructure) RollbackCreate(ctx context.Context, entry models.PostStructure) error {
	// delete entry
	err := p.dbProvider.GetDb().
		WithContext(ctx).
		Table("post_structures").
		Where("post_id = ?", entry.PostId).
		Delete(&models.PostStructure{}).
		Error
	if err != nil {
		return fmt.Errorf("error on rolling back post structure creation from database. err=%w", err)
	}

	return nil
}

func (p *postStructure) GetMultipleParentlessWithCursor(ctx context.Context, limit int, offset int) ([]*models.PostStructure, error) {
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
	if (offset) > int(totalRows) {
		return []*models.PostStructure{}, nil
	}

	// retrieve rows
	var out []*models.PostStructure
	err = p.dbProvider.GetDb().
		WithContext(ctx).
		Raw(`
			SELECT 
				p.post_id,
				p.parent_post_id,
				COUNT(c.post_id) as comment_count
			FROM post_structures p
				LEFT JOIN post_structures c ON p.post_id = c.parent_post_id
			WHERE
				p.parent_post_id IS NULL 
			GROUP BY p.post_id,
				p.parent_post_id 
			ORDER BY comment_count DESC
			LIMIT ?
			OFFSET ?
		`, limit, offset).
		Scan(&out).
		Error
	if err != nil {
		return out, fmt.Errorf("error on getting list post entry from database. err=%w", err)
	}

	return out, nil
}

func (p *postStructure) GetChildrenWithParentPostIdsAndChildrenLimitCount(ctx context.Context, ids []string, childrenLimit int) ([]*models.PostStructure, error) {
	// retrieve rows
	var out []*models.PostStructure
	err := p.dbProvider.GetDb().
		WithContext(ctx).
		Raw(`
			select 
				sub.post_id,
				sub.parent_post_id
			from post_structures ps 
				left join lateral (
					select 
						ps1.post_id,
						ps1.parent_post_id
					from post_structures ps1
					where ps1.parent_post_id = ps.post_id
					order by ps1.post_id desc
					limit ?
				) sub on true and sub.post_id is not NULL
			where 
				ps.parent_post_id is null 
				and sub.post_id is not null
				and ps.post_id IN(?) 
		`, childrenLimit, ids).
		Scan(&out).
		Error
	if err != nil {
		return out, fmt.Errorf("error on getting list post entry from database. err=%w", err)
	}

	return out, nil
}
