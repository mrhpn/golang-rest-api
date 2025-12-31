package media

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mrhpn/go-rest-api/internal/httpx"
)

type fileType string
type fileCategory string

const (
	fileCategoryProfiles   fileCategory = "profiles"
	fileCategoryThumbnails fileCategory = "thumbnails"
)

const (
	fileTypeImage fileType = "image"
	fileTypeVideo fileType = "video"
	fileTypeDoc   fileType = "document"
)

type filePolicy struct {
	AllowedExtensions map[string]bool
	MaxSize           int64
}

var policies = map[fileType]filePolicy{
	fileTypeImage: {
		AllowedExtensions: map[string]bool{
			".jpg":  true,
			".jpeg": true,
			".png":  true,
		},
		MaxSize: 5 * 1024 * 1024, // 5MB
	},
	fileTypeVideo: {
		AllowedExtensions: map[string]bool{
			".mp4": true,
			".mov": true,
			".avi": true,
		},
		MaxSize: 50 * 1024 * 1024, // 50MB
	},
	fileTypeDoc: {
		AllowedExtensions: map[string]bool{
			".pdf":  true,
			".docx": true,
			".txt":  true,
		},
		MaxSize: 10 * 1024 * 1024, // 10MB
	},
}

// Handler handles media-related HTTP endpoints such as uploads, retrieval, and media management operations.
type Handler struct {
	mediaService Service
}

// NewHandler constructs a media Handler with its required service dependency.
func NewHandler(mediaService Service) *Handler {
	return &Handler{mediaService: mediaService}
}

// UploadProfilePicture godoc
//
//	@Summary		Upload profile picture
//	@Description	Upload an image file to be used as a profile picture. Max 5MB.
//	@Tags			Media
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			file	formData	file	true	"Image file (jpg, jpeg, png)"
//	@Success		201		{object}	httpx.SuccessResponse{data=map[string]string}
//	@Failure		400		{object}	httpx.ErrorResponse
//	@Security		BearerAuth
//	@Router			/media/upload/profile [post]
func (h *Handler) UploadProfilePicture(c *gin.Context) {
	h.handleUpload(c, fileCategoryProfiles, fileTypeImage)
}

func (h *Handler) handleUpload(c *gin.Context, subDir fileCategory, category fileType) {
	// 1. get policy for type
	policy, exists := policies[category]
	if !exists {
		httpx.FailWithError(c, errInvalidFileTypeCategory)
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

	// 4. validate size
	if file.Size > policy.MaxSize {
		httpx.FailWithError(c, errFileTooLarge)
		return
	}

	// 5. validate extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !policy.AllowedExtensions[ext] {
		httpx.FailWithError(c, errInvalidFile)
		return
	}

	// 6. upload
	url, err := h.mediaService.Upload(file, subDir)
	if err != nil {
		httpx.FailWithError(c, err)
		return
	}

	httpx.OK(c, http.StatusCreated, gin.H{"url": url})
}
