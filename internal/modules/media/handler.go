package media

import (
	"context"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/mrhpn/go-rest-api/internal/httpx"
)

type mediaService interface {
	Upload(ctx context.Context, file *multipart.FileHeader, subDir fileCategory) (string, error)
}

// Handler handles media-related HTTP endpoints such as uploads, retrieval, and media management operations.
type Handler struct {
	service  mediaService
	policies map[fileType]filePolicy
}

// NewHandler constructs a media Handler with its required service dependency.
func NewHandler(service mediaService) *Handler {
	return &Handler{
		service:  service,
		policies: getDefaultPolicies(),
	}
}

// UploadProfilePicture godoc
//
//	@Summary		Upload profile picture
//	@Description	Upload an image file to be used as a profile picture. Max 5MB.
//	@Tags			Media
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			file	formData	file	true	"Image file (jpg, jpeg, png)"
//	@Success		201		{object}	httpx.SuccessResponse{data=media.Response}
//	@Failure		400		{object}	httpx.ErrorResponse
//	@Security		BearerAuth
//	@Router			/media/upload/profile [post]
func (h *Handler) UploadProfilePicture(c *gin.Context) {
	h.handleUpload(c, fileCategoryProfile, fileTypeImage)
}

func (h *Handler) handleUpload(c *gin.Context, subDir fileCategory, fileType fileType) {
	// 1. get policy for type
	policy, exists := h.policies[fileType]
	if !exists {
		httpx.FailWithError(c, errInvalidFileType)
		return
	}

	// 2. early check for file too large error
	if c.Request.ContentLength > policy.MaxSize {
		httpx.FailWithError(c, errFileTooLarge)
		return
	}

	// 3. parse file
	file, err := c.FormFile("file")
	if err != nil {
		httpx.FailWithError(c, errNoFileUploaded)
		return
	}

	// 4. validate size - check for empty or invalid size
	if file.Size <= 0 {
		httpx.FailWithError(c, errFileEmpty)
		return
	}

	// 5. validate maximum size
	if file.Size > policy.MaxSize {
		httpx.FailWithError(c, errFileTooLarge)
		return
	}

	// 6. validate extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !policy.AllowedExtensions[ext] {
		httpx.FailWithError(c, errInvalidFile)
		return
	}

	// 7. upload
	url, err := h.service.Upload(httpx.ReqCtx(c), file, subDir)
	if err != nil {
		httpx.FailWithError(c, err)
		return
	}

	httpx.OK(c, http.StatusCreated, ToResponse(url))
}
