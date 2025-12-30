package media

import "github.com/mrhpn/go-rest-api/internal/errors"

var (
	ErrInternal = errors.New(
		errors.Internal,
		"INTERNAL_ERROR",
		"internal server error",
	)

	ErrNoFileUploaded = errors.New(
		errors.BadRequest,
		"NO_FILE_UPLOADED",
		"no file uploaded",
	)

	// error if none of these: image, video, doc
	ErrInvalidFileTypeCategory = errors.New(
		errors.BadRequest,
		"INVALID_FILE_TYPE_CATEGORY",
		"invalid file type category",
	)

	ErrInvalidFile = errors.New(
		errors.BadRequest,
		"INVALID_FILE",
		"invalid file",
	)

	ErrFileTooLarge = errors.New(
		errors.BadRequest,
		"FILE_TOO_LARGE",
		"file too large",
	)

	ErrFileOpen = errors.New(
		errors.Internal,
		"FILE_OPEN_ERROR",
		"failed to process uploaded file",
	)

	ErrUploadToStorage = errors.New(
		errors.Internal,
		"STORAGE_UPLOAD_ERROR",
		"failed to upload file to storage",
	)

	ErrStorageHealthCheck = errors.New(
		errors.Internal,
		"STORAGE_HEALTH_CHECK_ERROR",
		"storage health check failed",
	)

	ErrStorageBucketMissing = errors.New(
		errors.Internal,
		"STORAGE_BUCKET_MISSING",
		"storage bucket missing",
	)
)
