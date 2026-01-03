package app

import (
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/mrhpn/go-rest-api/internal/config"
	"github.com/mrhpn/go-rest-api/internal/modules/media"
	"github.com/mrhpn/go-rest-api/internal/security"
)

// Context is the application context containing all required dependencies.
type Context struct {
	DB              *gorm.DB
	Redis           *redis.Client
	Cfg             *config.Config
	Logger          zerolog.Logger
	SecurityHandler *security.JWTHandler
	MediaService    media.Service
}
