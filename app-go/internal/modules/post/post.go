package post

import (
	"app-go/internal/modules/post/models"
	"context"
)

type IPost interface {
	GetPost(ctx context.Context, postId string) (*models.Post, error)
	GetPosts(ctx context.Context, args models.GetPostsArg) (*models.GetPostsResponse, error)
	CreatePost(ctx context.Context, args models.CreatePostArg) error
	Update(ctx context.Context, post *models.Post) error
}

type IComment interface {
	// Post creates a comment on a post
	Post(ctx context.Context, args models.PostCommentArg) error

	// Delete deletes a comment on a post. Deletion only allow for self-made comment
	Delete(ctx context.Context, args models.DeleteCommentArg) error
}
