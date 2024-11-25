package services

import (
	internalmodels "app-go/internal/models"
	"app-go/internal/modules/post/models"
	"app-go/internal/modules/post/repositories"
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/segmentio/ksuid"
)

type comment struct {
	postRepo          repositories.IPost
	postStructureRepo repositories.IPostStructure
}

func NewComment(
	postRepo repositories.IPost,
	postStructureRepo repositories.IPostStructure,
) *comment {
	return &comment{
		postRepo:          postRepo,
		postStructureRepo: postStructureRepo,
	}
}

func (c *comment) Post(ctx context.Context, args models.PostCommentArg) error {
	// validate post id format
	postId, err := ksuid.Parse(args.PostId)
	if err != nil {
		return internalmodels.ErrorInvalidId
	}

	// validate post existence
	existingPost, err := c.postRepo.GetById(ctx, postId)
	if err != nil {
		log.Printf("error on getting existing post from database. err=%v", err)
		return errors.New("error on getting existing post from database")
	}

	// construct comment post
	comment := &models.Post{
		Text:      args.Text,
		CreatedBy: args.Creator,
	}

	// create comment entry
	err = c.postRepo.Create(ctx, comment)
	if err != nil {
		log.Printf("error on creating comment entry to database. err=%v", err)
		return errors.New("error on creating comment entry to database")
	}

	// construct post-structure
	commentStructure := models.PostStructure{
		PostId: comment.Id,
		ParentPostId: sql.NullString{
			String: existingPost.Id.String(),
		},
	}

	// create comment-structure entry
	err = c.postStructureRepo.Create(ctx, commentStructure)
	if err != nil {
		log.Printf("error on creating comment structure to database. err=%v", err)
		return errors.New("error on creating comment structure to database")
	}

	return nil
}

func (c *comment) Delete(ctx context.Context, args models.DeleteCommentArg) error {
	// validate post id format
	commentId, err := ksuid.Parse(args.CommentId)
	if err != nil {
		return internalmodels.ErrorInvalidId
	}

	// validate post existence
	existingComment, err := c.postRepo.GetById(ctx, commentId)
	if err != nil {
		log.Printf("error on getting existing comment from database. commentId=%s. err=%v", commentId, err)
		return errors.New("error on getting existing comment from database")
	}

	// validate if comment belong to requester
	if existingComment.CreatedBy != args.Requester {
		log.Printf("'%s' is attempting to delete comment made by '%s'. commentId=%s", args.Requester, existingComment.CreatedBy, existingComment.Id.String())
		return nil
	}

	// do soft delete
	existingComment.DeletedAt = sql.NullTime{
		Time: time.Now().UTC(),
	}

	// update comment
	err = c.postRepo.Update(ctx, existingComment)
	if err != nil {
		log.Printf("error on updating existing comment into the database. commentId=%s. err=%v", existingComment.Id, err)
		return errors.New("error on getting existing post from database")
	}

	return nil
}
