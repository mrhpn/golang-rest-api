package health

import "time"

// Response represents the health check response
type Response struct {
	Status    string            `json:"status"`
	Checks    map[string]string `json:"checks,omitempty"`
	Timestamp string            `json:"timestamp"`
}

func ToResponse(status string, checks ...map[string]string) Response {
	checksMap := make(map[string]string)
	if len(checks) > 0 && checks[0] != nil {
		checksMap = checks[0]
	}

	return Response{
		Status:    status,
		Checks:    checksMap,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}
