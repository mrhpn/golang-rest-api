package apperror

var (
	// ErrInternal represents an unexpected internal server error occurring during authentication or authorization operations.
	ErrInternal = New(
		Internal,
		"INTERNAL_ERROR",
		"internal server error",
	)

	// ErrDatabaseError represents an unexpected internal database error occurring during database operations.
	ErrDatabaseError = New(
		Internal,
		"DATABASE_ERROR",
		"failed to perform database operation",
	)

	// ErrInvalidID indicates that the provided identifier does not match the expected format.
	ErrInvalidID = New(
		InvalidInput,
		"INVALID_ID_FORMAT",
		"invalid id format",
	)
)
