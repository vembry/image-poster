package filestorage

import (
	"app-go/internal/modules/file_storage/models"
	"context"
)

// IUpload specify expectation of an upload module
type IUpload interface {
	Upload(ctx context.Context, args models.UploadArgs) error
}
