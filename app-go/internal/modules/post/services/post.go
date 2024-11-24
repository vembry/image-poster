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
)

type post struct {
	postRepo          repositories.IPost
	postStructureRepo repositories.IPostStructure
	uploader          filestorage.IUpload
}

// New initialize post domain's service
func New(postRepo repositories.IPost, postStructureRepo repositories.IPostStructure, uploader filestorage.IUpload) *post {
	return &post{
		postRepo:          postRepo,
		postStructureRepo: postStructureRepo,
		uploader:          uploader,
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

	// construct new post entry
	post := &models.Post{
		Text:      args.Text,
		Image:     raw,
		CreatedBy: args.Creator,
	}

	// save post entry
	err = p.postRepo.Create(ctx, post)
	if err != nil {
		log.Printf("error on saving post to database. err=%v", err)
		return errors.New("error on saving post to database")
	}

	// append rollback handler
	rollbackFns = append(rollbackFns, func() {
		_ = p.postRepo.RollbackCreate(ctx, post)
	})

	// create post-structure entry
	err = p.postStructureRepo.Create(ctx, models.PostStructure{
		PostId: post.Id,
	})
	if err != nil {
		log.Printf("error on saving post structures to database. err=%v", err)
		return errors.New("error on saving post structures to database")
	}

	// enqueue post's for image transform

	// TODO: implements enquque for image transform
	// ...

	// return
	return nil
}

// GetPosts return a list of posts based on provided arguments
func (p *post) GetPosts(ctx context.Context, args models.GetPostsArg) (*models.GetPostsResponse, error) {
	// attempt to get available posts based on parameter
	postStructures, err := p.postStructureRepo.GetMultipleWithCursor(ctx, args.Limit, args.Limit*args.Page)
	if err != nil {
		log.Printf("error on getting post structures from database. err=%v", err)
		return nil, errors.New("error on getting post structures from database")
	}

	if len(postStructures) == 0 {
		return &models.GetPostsResponse{
			List: make([]models.PostResponse, 0),
		}, nil
	}

	// prep arg to retrieve posts by ids
	postIds := []string{}
	for _, postStpostStructure := range postStructures {
		postIds = append(postIds, postStpostStructure.PostId.String())
	}

	// retrieve full posts
	posts, err := p.postRepo.GetByMultipleIds(ctx, postIds)
	if err != nil {
		log.Printf("error on getting posts from database. err=%v", err)
		return nil, errors.New("error on getting post structures from database")
	}

	// convert posts into post-responses
	postResponses := []models.PostResponse{}
	for _, post := range posts {
		// read image links
		var image models.PostImage
		err = json.Unmarshal(post.Image, &image)
		if err != nil {
			log.Printf("error on parsing post's image json into a structured format")
			continue
		}

		postResponses = append(postResponses, models.PostResponse{
			Id:      post.Id,
			Text:    post.Text,
			Creator: post.CreatedBy,
			Image:   image,
		})
	}

	return &models.GetPostsResponse{
		List: postResponses,
	}, nil
}

// PostComment creates a comment entry of a post
func (p *post) PostComment(ctx context.Context, args models.PostCommentArg) error {
	return nil
}

// DeleteComment deletes a comment entry from a post
func (p *post) DeleteComment(ctx context.Context, args models.DeleteCommentArg) error {
	return nil
}
