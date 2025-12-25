package app

import (
	"github.com/mrhpn/go-rest-api/internal/config"
	"github.com/mrhpn/go-rest-api/internal/security"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type AppContext struct {
	DB              *gorm.DB
	Cfg             *config.Config
	Logger          zerolog.Logger
	SecurityHandler *security.JWTHandler
}
