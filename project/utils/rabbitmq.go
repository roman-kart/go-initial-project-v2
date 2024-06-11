package utils

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/roman-kart/go-initial-project/project/config"
)

type RabbitMQ struct {
	Config *config.Config
	Logger *Logger
}

func NewRabbitMQ(config *config.Config, logger *Logger) *RabbitMQ {
	return &RabbitMQ{
		Config: config,
		Logger: logger,
	}
}

func (r *RabbitMQ) GetConnectionString(vhost string) string {
	return fmt.Sprintf(
		"amqp://%s:%s@%s:%d/%s",
		r.Config.RabbitMQ.User,
		r.Config.RabbitMQ.Password,
		r.Config.RabbitMQ.Host,
		r.Config.RabbitMQ.Port,
		vhost,
	)
}

func (r *RabbitMQ) GetConnection(vhost string) (*amqp.Connection, error) {
	return amqp.Dial(r.GetConnectionString(vhost))
}
