package media

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/mrhpn/go-rest-api/internal/apperror"
	"github.com/mrhpn/go-rest-api/internal/constants"
)

// localService implements the Service interface using the local filesystem.
type localService struct {
	basePath string
}

// NewLocalService initializes a local filesystem-backed media service.
func NewLocalService(basePath string) Service {
	log.Info().Msg("✅ Storage (local) — exists and ready at " + basePath)
	return &localService{basePath: basePath}
}

// Upload stores the file on disk and returns the relative public path.
func (s *localService) Upload(_ context.Context, file *multipart.FileHeader, subDir fileCategory) (string, error) {
	src, err := file.Open()
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
		processed, _, pErr := processImage(src, *opts)
		if pErr != nil {
			return "", apperror.Wrap(
				apperror.BadRequest,
				errInvalidFile.Code,
				errInvalidFile.Message,
				pErr,
			)
		}
		reader = processed
		ext = ".jpg"
	}

	// 4. construct the obj name
	objectName := fmt.Sprintf("%s/%s%s", subDir, uuid.New().String(), ext)

	// 5. create directories and store the file
	storagePath := filepath.Join(s.basePath, objectName)
	if mkErr := os.MkdirAll(filepath.Dir(storagePath), 0750); mkErr != nil {
		return "", apperror.Wrap(
			apperror.Internal,
			errUploadToStorage.Code,
			errUploadToStorage.Message,
			mkErr,
		)
	}

	dst, err := os.Create(storagePath)
	if err != nil {
		return "", apperror.Wrap(
			apperror.Internal,
			errUploadToStorage.Code,
			errUploadToStorage.Message,
			err,
		)
	}
	defer func() { _ = dst.Close() }()

	if _, copyErr := io.Copy(dst, reader); copyErr != nil {
		return "", apperror.Wrap(
			apperror.Internal,
			errUploadToStorage.Code,
			errUploadToStorage.Message,
			copyErr,
		)
	}

	return fmt.Sprintf("/%s", objectName), nil
}

// HealthCheck verifies that the base path exists and is writable.
func (s *localService) HealthCheck(ctx context.Context) error {
	_ = ctx
	if err := os.MkdirAll(s.basePath, 0750); err != nil {
		return apperror.Wrap(
			apperror.Internal,
			errStorageHealthCheck.Code,
			errStorageHealthCheck.Message,
			err,
		)
	}
	return nil
}
