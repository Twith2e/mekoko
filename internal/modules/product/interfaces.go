package product

import (
	"context"
	"mime/multipart"
)

type FileUploader interface {
	UploadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error)
}
