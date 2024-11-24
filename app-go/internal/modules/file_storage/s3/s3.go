package s3

import (
	"app-go/internal/modules/file_storage/models"
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type s3 struct {
	client *awss3.Client

	s3bucket string
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
		// this probably better to be defined on env var, but
		// since theres no point doing so for now, i'll define
		// it explicitly here
		s3bucket: "public-images",
	}
}

// Upload handle upload process to aws' s3
func (s *s3) Upload(ctx context.Context, args models.UploadArgs) error {
	res, err := s.client.PutObject(ctx, &awss3.PutObjectInput{
		Bucket: &s.s3bucket,
		Key:    aws.String(fmt.Sprintf("file-%s", args.File.Name)),
		ACL:    types.ObjectCannedACLPublicRead,

		Body:        args.File.Content,
		ContentType: &args.File.ContentType,
	})
	if err != nil {
		err = fmt.Errorf("error on uploading file to s3. err=%w", err)
		log.Print(err)
		return err
	}

	log.Print(res)

	return nil
}
