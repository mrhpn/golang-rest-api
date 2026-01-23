package media

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"

	"github.com/mrhpn/go-rest-api/internal/apperror"
	"github.com/mrhpn/go-rest-api/internal/constants"
)

const healthCheckTimeout = 5 * time.Second

// minioService implements the Service interface using MinIO as the underlying object storage backend.
type minioService struct {
	client     *minio.Client
	bucketName string
}

// NewMinioService initializes a MinIO-backed media service with the provided connection credentials and bucket configuration.
func NewMinioService(client *minio.Client, bucketName string) Service {
	return &minioService{
		client:     client,
		bucketName: bucketName,
	}
}

// Upload streams the file to MinIO and returns the path
func (s *minioService) Upload(ctx context.Context, file *multipart.FileHeader, subDir fileCategory) (string, error) {
	src, err := file.Open() // 80
	if err != nil {
		return "", apperror.Wrap(
			apperror.Internal,
			errFileOpen.Code,
			errFileOpen.Message,
			err,
		)
	}
	defer func() { _ = src.Close() }()

	// 1. initialize vars with raw upload data
	var reader io.Reader = src
	size := file.Size
	contentType := file.Header.Get("Content-Type")
	ext := filepath.Ext(file.Filename)

	// 2. define processing rules based on directory
	var opts *imageOptions
	switch subDir {
	case fileCategoryProfile:
		opts = &imageOptions{
			MaxWidth:  constants.MaxProfileImageWidth,
			MaxHeight: constants.MaxProfileImageHeight,
			Quality:   constants.MaxProfileImageQuality,
		}
	case fileCategoryThumbnail:
		opts = &imageOptions{
			MaxWidth:  constants.MaxThumbnailImageWidth,
			MaxHeight: constants.MaxThumbnailImageHeight,
			Quality:   constants.MaxThumbnailImageQuality,
		}
	default:
		// no processing
	}

	// 3. apply processing if options were found
	if opts != nil {
		processed, newSize, pErr := processImage(src, *opts) // 116
		if pErr != nil {
			return "", apperror.Wrap(
				apperror.BadRequest,
				errInvalidFile.Code,
				errInvalidFile.Message,
				pErr,
			)
		}
		reader = processed
		size = newSize
		contentType = "image/jpeg"
		ext = ".jpg"
	}

	// 4. construct the obj name
	objectName := fmt.Sprintf("%s/%s%s", subDir, uuid.New().String(), ext)

	_, err = s.client.PutObject(ctx, s.bucketName, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", apperror.Wrap(
			apperror.Internal,
			errUploadToStorage.Code,
			errUploadToStorage.Message,
			err,
		)
	}

	return fmt.Sprintf("/%s", objectName), nil
}

// HealthCheck verifies that MinIO is accessible and the bucket exists
func (s *minioService) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, healthCheckTimeout)
	defer cancel()

	// Check if bucket exists and is accessible
	exists, err := s.client.BucketExists(ctx, s.bucketName)
	if err != nil {
		return apperror.Wrap(
			apperror.Internal,
			errStorageHealthCheck.Code,
			errStorageHealthCheck.Message,
			err,
		)
	}

	if !exists {
		return apperror.New(
			apperror.Internal,
			errStorageBucketMissing.Code,
			errStorageBucketMissing.Message,
		)
	}

	return nil
}
