package filestorage

import (
	"app-go/internal/modules/file_storage/models"
	"context"
)

// IUpload specify expectation of an upload module
type IUpload interface {
	// Upload uploads file to the storage provider
	Upload(ctx context.Context, args models.UploadArgs) error

	// RollbackUpload rolls back file uploaded to storage provider
	RollbackUpload(ctx context.Context, args models.UploadArgs) error
}

type IDownload interface {
	// Download downloads file from storage provider
	Download(ctx context.Context, fileKey string) (*models.DownloadResponse, error)
}
