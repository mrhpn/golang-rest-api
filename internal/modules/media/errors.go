package media

import "github.com/mrhpn/go-rest-api/internal/apperror"

var (
	errNoFileUploaded = apperror.New(
		apperror.BadRequest,
		"NO_FILE_UPLOADED",
		"no file uploaded",
	)

	// error if none of these: image, video, doc
	errInvalidFileType = apperror.New(
		apperror.BadRequest,
		"INVALID_FILE_TYPE",
		"invalid file type",
	)

	errInvalidFile = apperror.New(
		apperror.BadRequest,
		"INVALID_FILE",
		"invalid file",
	)

	errFileEmpty = apperror.New(
		apperror.BadRequest,
		"FILE_EMPTY",
		"file is empty or invalid",
	)

	errFileTooLarge = apperror.New(
		apperror.BadRequest,
		"FILE_TOO_LARGE",
		"file too large",
	)

	errFileOpen = apperror.New(
		apperror.Internal,
		"FILE_OPEN_ERROR",
		"failed to process uploaded file",
	)

	errUploadToStorage = apperror.New(
		apperror.Internal,
		"STORAGE_UPLOAD_ERROR",
		"failed to upload file to storage",
	)

	errStorageHealthCheck = apperror.New(
		apperror.Internal,
		"STORAGE_HEALTH_CHECK_ERROR",
		"storage health check failed",
	)

	errStorageBucketMissing = apperror.New(
		apperror.Internal,
		"STORAGE_BUCKET_MISSING",
		"storage bucket missing",
	)
)
