package main

import (
	"github.com/joho/godotenv"

	"github.com/mrhpn/go-rest-api/internal/config"
)

func setupConfig() (*config.Config, error) {
	_ = godotenv.Load()
	return config.Load()
}
