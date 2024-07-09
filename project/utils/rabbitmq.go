package utils

import (
	"fmt"
	"go.uber.org/zap"
	"net"
	"strconv"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/roman-kart/go-initial-project/v2/project/tools"
)

type RabbitMQConfig struct {
	Host     string
	Port     int
	User     string
	Password string
}

// RabbitMQ manipulates RabbitMQ connections.
type RabbitMQ struct {
	Config              *RabbitMQConfig
	ErrorWrapperCreator tools.ErrorWrapperCreator
	logger              *zap.Logger
}

// NewRabbitMQ creates a new instance of RabbitMQ.
// Using for configuring with wire.
func NewRabbitMQ(
	config *RabbitMQConfig,
	logger *zap.Logger,
	errorWrapperCreator tools.ErrorWrapperCreator,
) *RabbitMQ {
	return &RabbitMQ{
		Config:              config,
		logger:              logger.Named("RabbitMQ"),
		ErrorWrapperCreator: errorWrapperCreator.AppendToPrefix("RabbitMQ"),
	}
}

// GetConnectionString returns formated connection string.
func (r *RabbitMQ) GetConnectionString(vhost string) string {
	hostAndPort := net.JoinHostPort(
		r.Config.Host,
		strconv.Itoa(r.Config.Port),
	)

	return fmt.Sprintf(
		"amqp://%s:%s@%s/%s",
		r.Config.User,
		r.Config.Password,
		hostAndPort,
		vhost,
	)
}

// GetConnection create connection to RabbitMQ without caching.
func (r *RabbitMQ) GetConnection(vhost string) (*amqp.Connection, error) {
	ew := r.ErrorWrapperCreator.GetMethodWrapper("GetConnection")
	c, err := amqp.Dial(r.GetConnectionString(vhost))

	return c, ew(err)
}
