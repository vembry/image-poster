package s3

import (
	internalmodels "app-go/internal/models"
	"app-go/internal/modules/file_storage/models"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type s3 struct {
	client *awss3.Client

	defaultS3Bucket string
}

// New initiate s3 apis to do file upload/download
func New(awsCfg aws.Config) *s3 {
	opt := func(o *awss3.Options) {
		o.UsePathStyle = true
	}

	return &s3{
		// NOTE:
		// deliberately placing s3 initialization directly here,
		// if it's required by other modules, we should create an
		// isolated package to initialize this to avoid discrepancies
		client: awss3.NewFromConfig(awsCfg, opt),

		// NOTE:
		// this probably better to be defined on env var or on
		// parameter, but since theres no point doing so for now,
		// i'll just define it explicitly here
		defaultS3Bucket: "public-images",
	}
}

// Upload handle upload process to aws' s3
func (s *s3) Upload(ctx context.Context, args models.UploadArgs) error {
	_, err := s.client.PutObject(ctx, &awss3.PutObjectInput{
		Bucket: &s.defaultS3Bucket,
		Key:    &args.File.Name,
		ACL:    types.ObjectCannedACLPublicRead,

		Body:        bytes.NewReader(args.File.Content),
		ContentType: aws.String(string(args.File.ContentType)),
	})
	if err != nil {
		err = fmt.Errorf("error on uploading file to s3. err=%w", err)
		log.Print(err)
		return err
	}

	return nil
}

// RollbackUpload handle rollback process upload
func (s *s3) RollbackUpload(ctx context.Context, args models.UploadArgs) error {
	_, err := s.client.DeleteObject(ctx, &awss3.DeleteObjectInput{
		Bucket: &s.defaultS3Bucket,
		Key:    &args.File.Name,
	})
	if err != nil {
		err = fmt.Errorf("error on rolling back upload file to s3. err=%w", err)
		log.Print(err)
		return err
	}

	return nil
}

func (s *s3) Download(ctx context.Context, fileKey string) (*models.DownloadResponse, error) {
	// download from s3
	res, err := s.client.GetObject(ctx, &awss3.GetObjectInput{
		Bucket: &s.defaultS3Bucket,
		Key:    &fileKey,
	})
	if err != nil {
		return nil, fmt.Errorf("error on getting file from s3. err=%w", err)
	}
	if res != nil {
		defer res.Body.Close()
	}

	// read file
	imageRaw, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error on reading file from s3. err=%w", err)
	}

	fileBuffer := bytes.NewBuffer(imageRaw)

	return &models.DownloadResponse{
		File: internalmodels.File{
			ContentType: *res.ContentType,
			Content:     fileBuffer.Bytes(),
		},
	}, nil
}
