package health

import "time"

// Response represents the health check response
type Response struct {
	Status    string            `json:"status"`
	Checks    map[string]string `json:"checks,omitempty"`
	Timestamp string            `json:"timestamp"`
}

type RateLimitKeysDetails struct {
	Count int    `json:"count"`
	TTL   string `json:"ttl"`
}

type RateLimitResponse struct {
	TotalKeys int                             `json:"total_keys"`
	Keys      map[string]RateLimitKeysDetails `json:"keys"`
}

type RedisRateLimitResetResponse struct {
	Message     string `json:"message"`
	KeysDeleted int    `json:"keys_deleted"`
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

func ToRedisRateLimitResponse(totalKeys int, keys map[string]RateLimitKeysDetails) RateLimitResponse {
	return RateLimitResponse{
		TotalKeys: totalKeys,
		Keys:      keys,
	}
}

func ToRedisRateLimitResetResponse(message string, keysDeleted int) RedisRateLimitResetResponse {
	return RedisRateLimitResetResponse{
		Message:     message,
		KeysDeleted: keysDeleted,
	}
}
