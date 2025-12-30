package media

import (
	"context"
	"mime/multipart"
)

type Service interface {
	Upload(file *multipart.FileHeader, subDir FileCategory) (string, error)
	HealthCheck(ctx context.Context) error // Check if storage service is healthy
}
