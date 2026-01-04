package media

import (
	"context"
	"mime/multipart"
)

// Service defines the contract for media storage operations such as uploading files and performing storage health checks.
type Service interface {
	// Upload stores a file under the given category and returns the publicly accessible object path or identifier.
	Upload(ctx context.Context, file *multipart.FileHeader, subDir fileCategory) (string, error)

	// HealthCheck verifies that the underlying storage service is reachable and operational.
	HealthCheck(ctx context.Context) error // Check if storage service is healthy
}
