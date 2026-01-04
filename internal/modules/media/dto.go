package media

// Response returns the URL of an uploaded media file
type Response struct {
	URL string `json:"url"`
}

// ToResponse converts a URL string to MediaResponse DTO
func ToResponse(url string) Response {
	return Response{
		URL: url,
	}
}
