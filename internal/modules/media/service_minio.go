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
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/mrhpn/go-rest-api/internal/apperror"
	"github.com/mrhpn/go-rest-api/internal/constants"
)

const healthCheckTimeout = 5 * time.Second
const bucketCreateTimeout = 10 * time.Second

// minioService implements the Service interface using MinIO as the underlying object storage backend.
type minioService struct {
	client     *minio.Client
	bucketName string
}

// NewMinioService initializes a MinIO-backed media service with the provided connection credentials and bucket configuration.
func NewMinioService(host, accessKey, secretKey, bucketName string, useSSL bool) (Service, error) {
	client, err := minio.New(host, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize minio client: %w", err)
	}

	// 1. auto-create bucket if it doesn't exist
	ctx, cancel := context.WithTimeout(context.Background(), bucketCreateTimeout)
	defer cancel()

	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to check if minio bucket exists: %w", err)
	}

	if !exists {
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}

		// 2. set public policy so images can be viewed in browser via direct lin[k
		policy := fmt.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [{
					"Action": ["s3:GetObject"],
					"Effect":"Allow",
					"Principal":"*",
					"Resource":["arn:aws:s3:::%s/*"]
				}]
			}`, bucketName)

		err = client.SetBucketPolicy(ctx, bucketName, policy)
		if err != nil {
			return nil, fmt.Errorf("failed to set bucket policy: %w", err)
		}
	}

	return &minioService{
		client:     client,
		bucketName: bucketName,
	}, nil
}

// Upload streams the file to MinIO and returns the path
func (s *minioService) Upload(file *multipart.FileHeader, subDir fileCategory) (string, error) {
	ctx := context.Background()

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
	case fileCategoryProfiles:
		opts = &imageOptions{
			MaxWidth:  constants.MaxProfileImageWidth,
			MaxHeight: constants.MaxProfileImageHeight,
			Quality:   constants.MaxProfileImageQuality,
		}
	case fileCategoryThumbnails:
		opts = &imageOptions{
			MaxWidth:  constants.MaxThumbnailImageWidth,
			MaxHeight: constants.MaxThumbnailImageHeight,
			Quality:   constants.MaxThumbnailImageQuality,
		}
	}

	// 3. apply processing if options were found
	if opts != nil {
		processed, newSize, pErr := processImage(src, *opts) // 116
		if pErr == nil {
			reader = processed
			size = newSize
			contentType = "image/jpeg"
			ext = ".jpg"
		}
		// Note: If image processing fails, we continue with the original file
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
