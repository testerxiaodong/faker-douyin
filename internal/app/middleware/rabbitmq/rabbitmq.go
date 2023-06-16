package rabbitmq

import (
	"faker-douyin/internal/app/config"
	"faker-douyin/internal/app/log"
	"fmt"
	"github.com/google/wire"
	"github.com/streadway/amqp"
)

var ProviderSet = wire.NewSet(NewRabbitMQ, NewCommentRabbitMQ)

type RabbitMQ struct {
	Conn  *amqp.Connection
	mqurl string
}

func NewRabbitMQ(config *config.Config) *RabbitMQ {
	mqurl := fmt.Sprintf("amqp://%s:%s@%s:%s/", config.RabbitMq.User, config.RabbitMq.Password, config.RabbitMq.Host, config.RabbitMq.Port)
	connection, err := amqp.Dial(mqurl)
	if err != nil {
		log.AppLogger.Error(err.Error())
		panic(err)
	}
	fmt.Println("rabbitmq连接成功")
	return &RabbitMQ{
		Conn:  connection,
		mqurl: mqurl,
	}
}
