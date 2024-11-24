package post

import (
	"app-go/internal/modules/post/models"
	"context"
)

type IPost interface {
	GetPosts(ctx context.Context, args models.GetPostsArg) (*models.GetPostsResponse, error)
	CreatePost(ctx context.Context, args models.CreatePostArg) error
	PostComment(ctx context.Context, args models.PostCommentArg) error
	DeleteComment(ctx context.Context, args models.DeleteCommentArg) error
}
