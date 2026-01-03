package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/mrhpn/go-rest-api/internal/config"
	"github.com/mrhpn/go-rest-api/internal/constants"
)

func setupHTTPServer(cfg *config.Config, router *gin.Engine) *http.Server {
	return &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  constants.ServerReadTimeoutSecond * time.Second,  // Maximum duration for reading the entire request
		WriteTimeout: constants.ServerWriteTimeoutSecond * time.Second, // Maximum duration before timing out writes
		IdleTimeout:  constants.ServerIdleTimeoutSecond * time.Second,  // Maximum amount of time to wait for the next request
		// MaxHeaderBytes: 1 << 20, // 1 MB - uncomment if needed
	}
}
