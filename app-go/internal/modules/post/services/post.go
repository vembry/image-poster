package services

import (
	internalmodels "app-go/internal/models"
	filestorage "app-go/internal/modules/file_storage"
	filestoragemodels "app-go/internal/modules/file_storage/models"
	"app-go/internal/modules/post/models"
	"app-go/internal/modules/post/repositories"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/ksuid"
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
	var (
		rollbackFns []func() // to contain rollback functionalities, if needed
		err         error
	)
	defer func() {
		if err != nil {
			// rollback when function ends with error
			for _, rollbackFn := range rollbackFns {
				rollbackFn()
			}
		}
	}()

	// construct upload arg
	filename := fmt.Sprintf("%d-%s", time.Now().UnixMilli(), args.File.Name)
	uploadArg := filestoragemodels.UploadArgs{
		File: internalmodels.File{
			Name:        filename,
			ContentType: args.File.ContentType,
			Content:     args.File.Content,
		},
	}

	// upload file
	err = p.uploader.Upload(ctx, uploadArg)
	if err != nil {
		log.Printf("error on uploading file. err=%v", err)
		return errors.New("failed to upload file")
	}

	// append rollback handler
	rollbackFns = append(rollbackFns, func() {
		_ = p.uploader.RollbackUpload(ctx, uploadArg)
	})

	// construct image
	image := models.PostImage{
		Original: filename,
	}
	raw, _ := json.Marshal(image)

	// save post entry
	err = p.postRepo.Create(ctx, models.Post{
		Id:        ksuid.New(),
		Text:      args.Text,
		Image:     raw,
		CreatedBy: args.Creator,
	})
	if err != nil {
		log.Printf("error on saving post to database. err=%v", err)
		return errors.New("error on saving post to database")
	}

	// enqueue post's for image transform
	// ...

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
