package rabbitmq

import (
	"context"
	"encoding/json"
	"faker-douyin/internal/app/consts"
	"faker-douyin/internal/app/dao"
	"faker-douyin/internal/app/log"
	"github.com/redis/go-redis/v9"
	"github.com/streadway/amqp"
)

type CommentRabbitMQ struct {
	RabbitMQ     *RabbitMQ
	QueueName    string
	ExchangeName string
	RoutingKey   string
	RDB          *redis.Client
}

type CommentMessage struct {
	CommentDealType int64
	VideoId         int64
	CommentId       int64
}

func NewCommentRabbitMQ(rabbitMQ *RabbitMQ, rdb *redis.Client) *CommentRabbitMQ {
	var commentRabbitMQ CommentRabbitMQ
	commentRabbitMQ.RabbitMQ = rabbitMQ
	commentRabbitMQ.QueueName = consts.VideoCommentQueue
	commentRabbitMQ.RDB = rdb
	go commentRabbitMQ.Consume()
	return &commentRabbitMQ
}

func (cr *CommentRabbitMQ) Publish(message CommentMessage) {
	channel, err := cr.RabbitMQ.Conn.Channel()
	if err != nil {
		log.AppLogger.Error(err.Error())
	}
	_, err = channel.QueueDeclare(cr.QueueName, true, false, false, false, nil)
	if err != nil {
		log.AppLogger.Error(err.Error())
	}
	messageByte, err := json.Marshal(message)
	if err != nil {
		log.AppLogger.Error(err.Error())
	}
	err = channel.Publish(cr.ExchangeName, cr.QueueName, false, false, amqp.Publishing{
		ContentType:  "text/plain",
		DeliveryMode: amqp.Persistent,
		Body:         messageByte,
	})
	if err != nil {
		log.AppLogger.Error(err.Error())
	}
}

func (cr *CommentRabbitMQ) Consume() {
	channel, err := cr.RabbitMQ.Conn.Channel()
	if err != nil {
		log.AppLogger.Error(err.Error())
	}
	_, err = channel.QueueDeclare(cr.QueueName, true, false, false, false, nil)
	if err != nil {
		log.AppLogger.Error(err.Error())
	}
	deliveries, err := channel.Consume(cr.QueueName, cr.QueueName, false, false, false, false, nil)
	if err != nil {
		log.AppLogger.Error(err.Error())
	}
	for delivery := range deliveries {
		var message *CommentMessage
		err := json.Unmarshal(delivery.Body, message)
		if err != nil {
			log.AppLogger.Error(err.Error())
		}
		// 删除评论补偿消息
		if message.CommentDealType == 0 {
			result, err := cr.RDB.Exists(context.Background(), dao.GetRedisKeyByPrefix(consts.RedisVideoCommentPrefix, message.VideoId)).Result()
			if err != nil {
				// 发生错误，Nack并且重新入队
				err = channel.Reject(0, true)
				if err != nil {
					log.AppLogger.Error(err.Error())
				}
			}
			// key存在，进行补偿业务处理
			if result > 0 {
				_, err = cr.RDB.SRem(context.Background(), dao.GetRedisKeyByPrefix(consts.RedisVideoCommentPrefix, message.VideoId), message.CommentId).Result()
				// 删除失败，重新入队
				if err != nil {
					err = channel.Reject(0, true)
					if err != nil {
						log.AppLogger.Error(err.Error())
					}
				}
			}
			// 手动ack
			err = channel.Ack(0, false)
			if err != nil {
				log.AppLogger.Error(err.Error())
			}
		} else if message.CommentDealType == 1 {
			_, err = cr.RDB.SAdd(context.Background(), dao.GetRedisKeyByPrefix(consts.RedisVideoCommentPrefix, message.VideoId), message.CommentId).Result()
			if err != nil {
				err = channel.Reject(0, true)
				log.AppLogger.Error(err.Error())
			}
			err = channel.Ack(0, false)
			if err != nil {
				log.AppLogger.Error(err.Error())
			}
		}
	}
}
