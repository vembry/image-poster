package models

import (
	"app-go/internal/models"
	"encoding/json"
	"time"

	"github.com/segmentio/ksuid"
)

// Post is an entity that contain information of a post
type Post struct {
	Id        ksuid.KSUID     `gorm:"column:id"`
	Text      string          `gorm:"column:text"`
	Image     json.RawMessage `gorm:"column:image"`
	CreatedBy string          `gorm:"column:created_by"`
	CreatedAt time.Time       `gorm:"column:created_at"`
	UpdatedAt time.Time       `gorm:"column:updated_at"`
	DeletedAt time.Time       `gorm:"column:deleted_at"`
}

// PostImage contain image stored
type PostImage struct {
	Original    string `json:"original"`    // contain image link which uploaded by user
	Transformed string `json:"transformed"` // contain transformed image link from 'Original'

	// NOTE:
	// its probably better to define 'Original' and 'Transform' as array
	// for future usage, if we want to accept multiple image entry like
	// instagram. For now, keeping things simple
}

type CreatePostArg struct {
	Creator string      `json:"creator"`
	Text    string      `json:"text"`
	File    models.File `json:"file"`
}

type GetPostsArg struct{}
type PostCommentArg struct{}
type DeleteCommentArg struct{}
