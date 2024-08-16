package app

import (
	"github.com/go-redis/redis/v8"
	"github.com/mirjalilova/api-gateway/api"
	"github.com/mirjalilova/api-gateway/api/handlers"
	"github.com/mirjalilova/api-gateway/config"
	_ "github.com/mirjalilova/api-gateway/docs"
	prd "github.com/mirjalilova/api-gateway/pkg/kafka/producer"
	"golang.org/x/exp/slog"
)

func Run(cfg *config.Config) {

	// Kafka
	brokers := []string{"kafka:29092"}
	pr, err := prd.NewKafkaProducer(brokers)
	if err != nil {
		slog.Error("Failed to create Kafka producer:", err)
		return
	}

	// Redis
	rd := redis.NewClient(&redis.Options{
		Addr:     "redis:6379", 
		Password: "",               
		DB:       0,                
	})

	h := handlers.NewHandlerStruct(&pr, rd)

	router := api.Engine(h)
	if err := router.Run(cfg.API_GATEWAY); err != nil {
		slog.Error("Failed to start API Gateway: %v", err)
	}
}
