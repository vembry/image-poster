package models

// FileContentType is enum types containing file type
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
	ContentType string
	Content     []byte
}
