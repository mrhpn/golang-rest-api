package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mrhpn/go-rest-api/internal/config"
)

func setupHTTPServer(cfg *config.Config, router *gin.Engine) *http.Server {
	return &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}
