package app

import (
	"github.com/mirjalilova/api-gateway/api"
	"github.com/mirjalilova/api-gateway/api/handlers"
	"github.com/mirjalilova/api-gateway/config"
	_ "github.com/mirjalilova/api-gateway/docs"
	"golang.org/x/exp/slog"
)

func Run(cfg *config.Config) {

	h := handlers.NewHandlerStruct()

	router := api.Engine(h)
	if err := router.Run(cfg.API_GATEWAY); err != nil {
		slog.Error("Failed to start API Gateway: %v", err)
	}
}
