package s3

import (
	"app-go/internal/modules/file_storage/models"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
)

type s3 struct {
	client   *awss3.Client
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
		// this probably better to be defined, but i defined it
		// directly here since theres no point doing so for now
		s3bucket: "image-poster",
	}
}

// Upload handle upload process to aws' s3
func (s *s3) Upload(ctx context.Context, args models.UploadArgs) error {
	return nil
}
