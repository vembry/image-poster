package models

import "mime/multipart"

type FileContentType string

const (
	FileContentTypeJPEG FileContentType = "image/jpeg"
	FileContentTypeJPG  FileContentType = "image/jpg"
	FileContentTypePNG  FileContentType = "image/png"
	FileContentTypeBMP  FileContentType = "image/bmp"

	// add more if there are needed
	// ...
)

// File contain basic information to pass file informations
type File struct {
	Name        string
	ContentType FileContentType
	Content     multipart.File
}
