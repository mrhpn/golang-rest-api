package media

import "mime/multipart"

type Service interface {
	Upload(file *multipart.FileHeader, subDir FileCategory) (string, error)
}
