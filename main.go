package main

import (
	"github.com/mirjalilova/api-gateway/config"
	"github.com/mirjalilova/api-gateway/pkg/app"
)

func main() {
	// Load configuration and start the API Gateway
	cfg := config.Load()

	// Set up runtime profiling if enabled
	app.Run(&cfg)
}
