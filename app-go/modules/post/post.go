package post

import (
	"app-go/modules/post/models"
	"context"
)

type IPost interface {
	GetPosts(ctx context.Context, args models.GetPostsArg) ([]models.Post, error)
	CreatePost(ctx context.Context, args models.CreatePostArg) error
	PostComment(ctx context.Context, args models.PostCommentArg) error
	DeleteComment(ctx context.Context, args models.DeleteCommentArg) error
}
