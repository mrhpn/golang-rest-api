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
)
