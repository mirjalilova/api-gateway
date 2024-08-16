package handlers

import (
	"github.com/go-redis/redis/v8"
	"github.com/mirjalilova/api-gateway/client"
	kafka "github.com/mirjalilova/api-gateway/pkg/kafka/producer"
)

type Handlers struct {
	Clients  client.Clients
	Producer kafka.KafkaProducer
	Redis    *redis.Client
}

func NewHandlerStruct(pr *kafka.KafkaProducer, rd *redis.Client) *Handlers {
	return &Handlers{
		Clients:  *client.NewClients(),
		Producer: *pr,
		Redis:    rd,
	}
}
