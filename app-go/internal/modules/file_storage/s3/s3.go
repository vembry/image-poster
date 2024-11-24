package s3

import (
	"app-go/internal/modules/file_storage/models"
	"context"
)

type s3 struct {
	// dependencies
}

// New initiate s3 apis to do file upload/download
func New() *s3 {
	return &s3{}
}

// Upload handle upload process to aws' s3
func (s *s3) Upload(ctx context.Context, args models.UploadArgs) error {
	return nil
}
