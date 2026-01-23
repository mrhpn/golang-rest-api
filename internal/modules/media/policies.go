package media

import "github.com/mrhpn/go-rest-api/internal/constants"

type fileType string     // fileType represents image, video, document, etc
type fileCategory string // fileCategory represents profile, thumbnail, etc

const (
	fileCategoryProfile   fileCategory = "profile"
	fileCategoryThumbnail fileCategory = "thumbnail"
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

func getDefaultPolicies() map[fileType]filePolicy {
	return map[fileType]filePolicy{
		fileTypeImage: {
			AllowedExtensions: map[string]bool{
				".jpg":  true,
				".jpeg": true,
				".png":  true,
			},
			MaxSize: constants.MaxImageSize,
		},
		fileTypeVideo: {
			AllowedExtensions: map[string]bool{
				".mp4": true,
				".mov": true,
				".avi": true,
			},
			MaxSize: constants.MaxVideoSize,
		},
		fileTypeDoc: {
			AllowedExtensions: map[string]bool{
				".pdf":  true,
				".docx": true,
				".txt":  true,
			},
			MaxSize: constants.MaxDocumentSize,
		},
	}
}
