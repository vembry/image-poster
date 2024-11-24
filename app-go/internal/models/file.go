package models

import "mime/multipart"

// File contain basic information to pass file informations
type File struct {
	Name        string
	ContentType string
	Content     multipart.File
}
