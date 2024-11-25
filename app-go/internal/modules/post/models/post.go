package models

import (
	"app-go/internal/models"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/segmentio/ksuid"
)

// Post is an entity that contain fields of a 'posts' table
type Post struct {
	Id        ksuid.KSUID     `gorm:"column:id"`
	Text      string          `gorm:"column:text"`
	Image     json.RawMessage `gorm:"column:image"`
	CreatedBy string          `gorm:"column:created_by"`
	CreatedAt time.Time       `gorm:"column:created_at"`
	UpdatedAt time.Time       `gorm:"column:updated_at"`
	DeletedAt sql.NullTime    `gorm:"column:deleted_at"`
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

// PostStructure is an entity that contain fields of a 'post_structures' table
type PostStructure struct {
	PostId       ksuid.KSUID    `gorm:"column:post_id"`
	ParentPostId sql.NullString `gorm:"column:parent_post_id"`
}

type CreatePostArg struct {
	Creator string
	Text    string
	File    models.File
}

type GetPostsArg struct {
	Limit int
	Page  int
}

type GetPostsResponse struct {
	List []PostResponse `json:"list"`
}

type PostResponse struct {
	Id      ksuid.KSUID `json:"id"`
	Text    string      `json:"text"`
	Creator string      `json:"creator"`
	Image   PostImage   `json:"image"`
}

type PostCommentArg struct {
	PostId  string
	Text    string
	Creator string
}

type DeleteCommentArg struct {
	CommentId string
	Requester string
}
