package services

import (
	internalmodels "app-go/internal/models"
	filestorage "app-go/internal/modules/file_storage"
	filestoragemodels "app-go/internal/modules/file_storage/models"
	"app-go/internal/modules/post/models"
	"app-go/internal/modules/post/repositories"
	"app-go/internal/workers"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/ksuid"
)

type post struct {
	postRepo             repositories.IPost
	postStructureRepo    repositories.IPostStructure
	uploader             filestorage.IUpload
	imageTransformWorker workers.IImageTransformWorker
}

// New initialize post domain's service
func New(
	postRepo repositories.IPost,
	postStructureRepo repositories.IPostStructure,
	uploader filestorage.IUpload,
	imageTransformWorker workers.IImageTransformWorker,
) *post {
	return &post{
		postRepo:             postRepo,
		postStructureRepo:    postStructureRepo,
		uploader:             uploader,
		imageTransformWorker: imageTransformWorker,
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

	// append rollback handler for upload file
	rollbackFns = append(
		[]func(){
			func() {
				_ = p.uploader.RollbackUpload(ctx, uploadArg)
			},
		},
		rollbackFns...,
	)

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

	// append rollback handler for post creation
	rollbackFns = append(
		[]func(){
			func() {
				_ = p.postRepo.RollbackCreate(ctx, post)
			},
		},
		rollbackFns...,
	)

	// construct payload
	postStructure := models.PostStructure{
		PostId: post.Id,
	}
	// create post-structure entry
	err = p.postStructureRepo.Create(ctx, postStructure)
	if err != nil {
		log.Printf("error on saving post structures to database. err=%v", err)
		return errors.New("error on saving post structures to database")
	}

	// append rollback handler for post-structure creation
	rollbackFns = append(
		[]func(){
			func() {
				_ = p.postStructureRepo.RollbackCreate(ctx, postStructure)
			},
		},
		rollbackFns...,
	)

	// enqueue post's for image transform
	err = p.imageTransformWorker.Enqueue(ctx, post.Id.String())
	if err != nil {
		log.Printf("error on enqueuing post for image transform. err=%v", err)
		return errors.New("error on enqueuing post for image transform")
	}

	// return
	return nil
}

func (p *post) GetPost(ctx context.Context, postId string) (*models.Post, error) {
	id, err := ksuid.Parse(postId)
	if err != nil {
		return nil, internalmodels.ErrorInvalidId
	}

	post, err := p.postRepo.GetById(ctx, id)
	if err != nil {
		log.Printf("error on retrieving post from database. postId=%s. err=%v", postId, err)
		return nil, errors.New("error on retrieving post from database")
	}
	return post, nil
}

// GetPosts return a list of posts based on provided arguments
func (p *post) GetPosts(ctx context.Context, args models.GetPostsArg) (*models.GetPostsResponse, error) {
	// attempt to get available posts based on parameter
	postStructures, err := p.postStructureRepo.GetMultipleWithCursor(ctx, args.Limit, args.Limit*(args.Page-1))
	if err != nil {
		log.Printf("error on getting post structures from database. err=%v", err)
		return nil, errors.New("error on getting post structures from database")
	}

	if len(postStructures) == 0 {
		return &models.GetPostsResponse{
			List: make([]*models.PostResponse, 0),
		}, nil
	}

	// prep arg to retrieve posts by ids
	postIds := []string{}

	// constructing this so we could establish list's order
	// and rely on pointer to populate the responses
	postResponseMap := map[string]*models.PostResponse{}
	postResponses := []*models.PostResponse{}

	// iterate through structure retrieved
	for _, postStpostStructure := range postStructures {
		postId := postStpostStructure.PostId.String()
		postIds = append(postIds, postId)

		postResponse := &models.PostResponse{
			Id: postStpostStructure.PostId,
		}

		postResponseMap[postId] = postResponse
		postResponses = append(postResponses, postResponse)
	}

	// retrieve full posts
	posts, err := p.postRepo.GetByMultipleIds(ctx, postIds)
	if err != nil {
		log.Printf("error on getting posts from database. err=%v", err)
		return nil, errors.New("error on getting post structures from database")
	}

	// convert posts into post-responses
	for _, post := range posts {
		// read image links
		var image models.PostImage
		err = json.Unmarshal(post.Image, &image)
		if err != nil {
			log.Printf("error on parsing post's image json into a structured format")
			continue
		}

		val, ok := postResponseMap[post.Id.String()]
		if !ok {
			continue
		}

		val.Text = post.Text
		val.Creator = post.CreatedBy
		val.Image = image

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

// PostComment creates a comment entry of a post
func (p *post) Update(ctx context.Context, post *models.Post) error {
	err := p.postRepo.Update(ctx, post)
	if err != nil {
		log.Printf("error on updating post to database. postId=%s. err=%v", post.Id.String(), err)
		return errors.New("error on updating post to database")
	}
	return nil
}
