package services

import (
	"app-go/modules/post/models"
	"app-go/modules/post/repositories"
	"context"
)

type post struct {
	postRepo repositories.IPost
}

// New initialize post domain's service
func New(postRepo repositories.IPost) *post {
	return &post{
		postRepo: postRepo,
	}
}

// GetPosts return a list of posts based on provided arguments
func (p *post) GetPosts(ctx context.Context, args models.GetPostsArg) ([]models.Post, error) {
	return nil, nil
}

// CreatePost creates a new post entry
func (p *post) CreatePost(ctx context.Context, args models.CreatePostArg) error {
	return nil
}

// PostComment creates a comment entry of a post
func (p *post) PostComment(ctx context.Context, args models.PostCommentArg) error {
	return nil
}

// DeleteComment deletes a comment entry from a post
func (p *post) DeleteComment(ctx context.Context, args models.DeleteCommentArg) error {
	return nil
}
