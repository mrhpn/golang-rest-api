package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mrhpn/go-rest-api/internal/modules/health"
	"github.com/mrhpn/go-rest-api/internal/modules/users"
	"gorm.io/gorm"
)

func Register(router *gin.Engine, db *gorm.DB) {
	api := router.Group("/api")

	// health
	health.Register(api)

	// users
	userRepo := users.NewRepository(db)
	userService := users.NewService(userRepo)
	users.NewHandler(api, userService)
}
