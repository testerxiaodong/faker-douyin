package rabbitmq

import (
	"faker-douyin/internal/app/config"
	"testing"
)

func TestNewRabbitMQ(t *testing.T) {
	t.Parallel()
	c := config.NewConfig()
	NewRabbitMQ(c)
}
