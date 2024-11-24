package models

import "app-go/internal/models"

// Post is an entity that contain information of a post
type Post struct{}

type CreatePostArg struct {
	Creator string      `json:"creator"`
	Text    string      `json:"text"`
	File    models.File `json:"file"`
}

type GetPostsArg struct{}
type PostCommentArg struct{}
type DeleteCommentArg struct{}
