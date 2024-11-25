package sqs

import (
	"app-go/internal/models"
	"app-go/internal/modules/post"
	"context"
	"errors"
	"fmt"

	awssqs "github.com/aws/aws-sdk-go-v2/service/sqs"
)

type imageTransformer struct {
	client   *awssqs.Client
	queueUrl string

	postProvider post.IPost
}

func NewImageTransformer(queueUrl string) *imageTransformer {
	return &imageTransformer{
		queueUrl: queueUrl,
	}
}

func (it *imageTransformer) InjectDeps(postProvider post.IPost) {
	it.postProvider = postProvider
}

func (it *imageTransformer) GetQueueUrl() *string {
	return &it.queueUrl
}

func (it *imageTransformer) injectClient(client *awssqs.Client) {
	it.client = client
}

func (it *imageTransformer) Handle(ctx context.Context, body string) error {

	post, err := it.postProvider.GetPost(ctx, body)
	if err != nil {
		if errors.Is(err, models.ErrorInvalidId) {
			// break away when error return invalid id
			return nil
		}
	}

	if post == nil {
		// break away when post not found
		return nil
	}

	// start transforming image
	// ...

	// update post
	err = it.postProvider.Update(ctx, post)
	if err != nil {
		return fmt.Errorf("error on updating post entry. postId=%s. err=%v", post.Id.String(), err)
	}

	return nil
}

func (it *imageTransformer) Enqueue(ctx context.Context, postId string) error {
	_, err := it.client.SendMessage(ctx, &awssqs.SendMessageInput{
		QueueUrl:    it.GetQueueUrl(),
		MessageBody: &postId,
	})
	return err
}
