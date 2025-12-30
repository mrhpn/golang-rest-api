package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mrhpn/go-rest-api/internal/config"
)

func setupHTTPServer(cfg *config.Config, router *gin.Engine) *http.Server {
	// Server configuration
	readTimeout := 20 * time.Second
	writeTimeout := 30 * time.Second // Longer for file uploads
	idleTimeout := 120 * time.Second

	return &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  readTimeout,  // Maximum duration for reading the entire request
		WriteTimeout: writeTimeout, // Maximum duration before timing out writes
		IdleTimeout:  idleTimeout,  // Maximum amount of time to wait for the next request
		// MaxHeaderBytes: 1 << 20, // 1 MB - uncomment if needed
	}
}
