package services

import (
	internalmodels "app-go/internal/models"
	filestorage "app-go/internal/modules/file_storage"
	filestoragemodels "app-go/internal/modules/file_storage/models"
	"app-go/internal/modules/post/models"
	"app-go/internal/modules/post/repositories"
	"context"
	"errors"
	"fmt"
	"log"
	"time"
)

type post struct {
	postRepo repositories.IPost
	uploader filestorage.IUpload
}

// New initialize post domain's service
func New(postRepo repositories.IPost, uploader filestorage.IUpload) *post {
	return &post{
		postRepo: postRepo,
		uploader: uploader,
	}
}

// CreatePost creates a new post entry
func (p *post) CreatePost(ctx context.Context, args models.CreatePostArg) error {
	// upload file
	err := p.uploader.Upload(ctx, filestoragemodels.UploadArgs{
		File: internalmodels.File{
			Name:        fmt.Sprintf("%d-%s", time.Now().UnixMilli(), args.File.Name), // append
			ContentType: args.File.ContentType,
			Content:     args.File.Content,
		},
	})
	if err != nil {
		log.Printf("error on uploading file. err=%v", err)
		return errors.New("failed to upload file")
	}

	// save as post entry
	// ...

	log.Print(args.File)

	// return
	return nil
}

// GetPosts return a list of posts based on provided arguments
func (p *post) GetPosts(ctx context.Context, args models.GetPostsArg) ([]models.Post, error) {
	return nil, nil
}

// PostComment creates a comment entry of a post
func (p *post) PostComment(ctx context.Context, args models.PostCommentArg) error {
	return nil
}

// DeleteComment deletes a comment entry from a post
func (p *post) DeleteComment(ctx context.Context, args models.DeleteCommentArg) error {
	return nil
}
