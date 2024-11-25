package sqs

import (
	"app-go/internal/models"
	filestorage "app-go/internal/modules/file_storage"
	filestoragemodels "app-go/internal/modules/file_storage/models"
	"app-go/internal/modules/post"
	postmodels "app-go/internal/modules/post/models"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"strings"
	"time"

	awssqs "github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/nfnt/resize"
)

type imageTransformer struct {
	client   *awssqs.Client
	queueUrl string

	postProvider     post.IPost
	downloadProvider filestorage.IDownload
	uploadProvider   filestorage.IUpload
}

func NewImageTransformer(queueUrl string) *imageTransformer {
	return &imageTransformer{
		queueUrl: queueUrl,
	}
}

func (it *imageTransformer) InjectDeps(postProvider post.IPost, downloadProvider filestorage.IDownload, uploadProvider filestorage.IUpload) {
	it.postProvider = postProvider
	it.downloadProvider = downloadProvider
	it.uploadProvider = uploadProvider
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

	// download file
	var postImage postmodels.PostImage
	err = json.Unmarshal(post.Image, &postImage)
	if err != nil {
		log.Printf("error on unmarshaling post's image for image transform. postId=%s. err=%v", post.Id.String(), err)
		return nil
	}

	file, err := it.downloadProvider.Download(ctx, postImage.Original)
	if err != nil {
		return fmt.Errorf("error on downloading file from download-provider. err=%w", err)
	}
	buffer := bytes.NewReader(file.File.Content)

	img, _, err := image.Decode(buffer)
	if err != nil {
		log.Printf("error on decoding image to be transformed. postId=%s. err=%v", post.Id.String(), err)
		return nil
	}

	// start transforming image
	resizedImage := resize.Resize(600, 600, img, resize.Lanczos2)
	resizedImageBuffer := new(bytes.Buffer)
	err = jpeg.Encode(resizedImageBuffer, resizedImage, nil)
	if err != nil {
		log.Printf("error on encoding transformed image into jpeg. postId=%s. err=%v", post.Id.String(), err)
		return nil
	}

	// construct new name
	filenameArr := strings.Split(postImage.Original, ".")
	filenameArr = filenameArr[:len(filenameArr)-1] // remove file type
	filename := strings.Join(filenameArr, "")
	filename = fmt.Sprintf("%s-transformed-%d.jpg", filename, time.Now().UnixMilli())
	postImage.Transformed = filename

	// upload
	err = it.uploadProvider.Upload(ctx, filestoragemodels.UploadArgs{
		File: models.File{
			Name:        filename, // append indicator that this is transformed file
			ContentType: string(models.FileContentTypeJPEG),
			Content:     resizedImageBuffer.Bytes(),
		},
	})
	if err != nil {
		return fmt.Errorf("error on uploading file to file storage provider. err=%w", err)
	}

	// assign image to the post
	post.Image, _ = json.Marshal(postImage)

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
