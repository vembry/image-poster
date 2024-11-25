package workers

import "context"

type IImageTransformWorker interface {
	Enqueue(ctx context.Context, postId string) error
}
