package utils

import (
	"fmt"
	"net"
	"strconv"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/roman-kart/go-initial-project/v2/project/config"
	"github.com/roman-kart/go-initial-project/v2/project/tools"
)

// RabbitMQ manipulates RabbitMQ connections.
type RabbitMQ struct {
	Config              *config.Config
	Logger              *Logger
	ErrorWrapperCreator tools.ErrorWrapperCreator
}

// NewRabbitMQ creates a new instance of RabbitMQ.
// Using for configuring with wire.
func NewRabbitMQ(
	config *config.Config,
	logger *Logger,
	errorWrapperCreator tools.ErrorWrapperCreator,
) *RabbitMQ {
	return &RabbitMQ{
		Config:              config,
		Logger:              logger,
		ErrorWrapperCreator: errorWrapperCreator.AppendToPrefix("RabbitMQ"),
	}
}

// GetConnectionString returns formated connection string.
func (r *RabbitMQ) GetConnectionString(vhost string) string {
	hostAndPort := net.JoinHostPort(
		r.Config.RabbitMQ.Host,
		strconv.Itoa(r.Config.RabbitMQ.Port),
	)

	return fmt.Sprintf(
		"amqp://%s:%s@%s/%s",
		r.Config.RabbitMQ.User,
		r.Config.RabbitMQ.Password,
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
