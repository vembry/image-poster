package models

import "mime/multipart"

// Post is an entity that contain information of a post
type Post struct{}

type CreatePostArg struct {
	Text string         `json:"text"`
	File multipart.File `json:"file"`
}

type GetPostsArg struct{}
type PostCommentArg struct{}
type DeleteCommentArg struct{}
