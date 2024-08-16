package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer interface {
	ProduceMessages(topic string, message []byte) error
	Close() error
}

type Producer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(brokers []string) (KafkaProducer, error) {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		AllowAutoTopicCreation: true,
	}
	fmt.Println("lll")
	return &Producer{writer: writer}, nil
}

func (p *Producer) ProduceMessages(topic string, message []byte) error {
	return p.writer.WriteMessages(context.Background(), kafka.Message{
		Topic: topic,
		Value: message,
	})
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
