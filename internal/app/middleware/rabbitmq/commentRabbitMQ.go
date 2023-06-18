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

type VideoCommentRabbitMQ struct {
	RabbitMQ     *RabbitMQ
	QueueName    string
	ExchangeName string
	RoutingKey   string
	Rdb          *redis.Client
}

type CommentVideoRabbitMQ struct {
	RabbitMQ     *RabbitMQ
	QueueName    string
	ExchangeName string
	RoutingKey   string
	Rdb          *redis.Client
}

type CommentMessage struct {
	CommentDealType int64
	VideoId         int64
	CommentId       int64
}

func NewVideoCommentRabbitMQ(rabbitMQ *RabbitMQ, rdb *redis.Client) *VideoCommentRabbitMQ {
	var videoCommentRabbitMQ VideoCommentRabbitMQ
	videoCommentRabbitMQ.RabbitMQ = rabbitMQ
	videoCommentRabbitMQ.QueueName = consts.VideoCommentQueue
	videoCommentRabbitMQ.Rdb = rdb
	go videoCommentRabbitMQ.Consume()
	return &videoCommentRabbitMQ
}

func NewCommentVideoRabbitMQ(rabbitMQ *RabbitMQ, rdb *redis.Client) *CommentVideoRabbitMQ {
	var commentVideoRabbitMQ CommentVideoRabbitMQ
	commentVideoRabbitMQ.RabbitMQ = rabbitMQ
	commentVideoRabbitMQ.QueueName = consts.CommentVideoQueue
	commentVideoRabbitMQ.Rdb = rdb
	go commentVideoRabbitMQ.Consume()
	return &commentVideoRabbitMQ
}

func (vcr *VideoCommentRabbitMQ) Publish(message CommentMessage) {
	channel, err := vcr.RabbitMQ.Conn.Channel()
	if err != nil {
		log.AppLogger.Error(err.Error())
	}
	_, err = channel.QueueDeclare(vcr.QueueName, true, false, false, false, nil)
	if err != nil {
		log.AppLogger.Error(err.Error())
	}
	messageByte, err := json.Marshal(message)
	if err != nil {
		log.AppLogger.Error(err.Error())
	}
	err = channel.Publish(vcr.ExchangeName, vcr.QueueName, false, false, amqp.Publishing{
		ContentType:  "text/plain",
		DeliveryMode: amqp.Persistent,
		Body:         messageByte,
	})
	if err != nil {
		log.AppLogger.Error(err.Error())
	}
}

func (vcr *VideoCommentRabbitMQ) Consume() {
	channel, err := vcr.RabbitMQ.Conn.Channel()
	if err != nil {
		log.AppLogger.Error(err.Error())
	}
	_, err = channel.QueueDeclare(vcr.QueueName, true, false, false, false, nil)
	if err != nil {
		log.AppLogger.Error(err.Error())
	}
	deliveries, err := channel.Consume(vcr.QueueName, vcr.QueueName, false, false, false, false, nil)
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
			result, err := vcr.Rdb.Exists(context.Background(), dao.GetRedisKeyByPrefix(consts.RedisVideoCommentPrefix, message.VideoId)).Result()
			if err != nil {
				// 发生错误，Nack并且重新入队
				err = channel.Reject(0, true)
				if err != nil {
					log.AppLogger.Error(err.Error())
				}
			}
			// key存在，进行补偿业务处理
			if result > 0 {
				_, err = vcr.Rdb.SRem(context.Background(), dao.GetRedisKeyByPrefix(consts.RedisVideoCommentPrefix, message.VideoId), message.CommentId).Result()
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
			_, err = vcr.Rdb.SAdd(context.Background(), dao.GetRedisKeyByPrefix(consts.RedisVideoCommentPrefix, message.VideoId), message.CommentId).Result()
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

func (cvr *CommentVideoRabbitMQ) Publish(message CommentMessage) {
	channel, err := cvr.RabbitMQ.Conn.Channel()
	if err != nil {
		log.AppLogger.Error(err.Error())
	}
	_, err = channel.QueueDeclare(cvr.QueueName, true, false, false, false, nil)
	if err != nil {
		log.AppLogger.Error(err.Error())
	}
	messageByte, err := json.Marshal(message)
	if err != nil {
		log.AppLogger.Error(err.Error())
	}
	err = channel.Publish(cvr.ExchangeName, cvr.QueueName, false, false, amqp.Publishing{
		ContentType:  "text/plain",
		DeliveryMode: amqp.Persistent,
		Body:         messageByte,
	})
	if err != nil {
		log.AppLogger.Error(err.Error())
	}
}

func (cvr *CommentVideoRabbitMQ) Consume() {
	channel, err := cvr.RabbitMQ.Conn.Channel()
	if err != nil {
		log.AppLogger.Error(err.Error())
	}
	_, err = channel.QueueDeclare(cvr.QueueName, true, false, false, false, nil)
	if err != nil {
		log.AppLogger.Error(err.Error())
	}
	deliveries, err := channel.Consume(cvr.QueueName, cvr.QueueName, false, false, false, false, nil)
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
			result, err := cvr.Rdb.Exists(context.Background(), dao.GetRedisKeyByPrefix(consts.RedisCommentVideoPrefix, message.CommentId)).Result()
			if err != nil {
				// 发生错误，Nack并且重新入队
				err = channel.Reject(0, true)
				if err != nil {
					log.AppLogger.Error(err.Error())
				}
			}
			// key存在，进行补偿业务处理
			if result > 0 {
				_, err = cvr.Rdb.SRem(context.Background(), dao.GetRedisKeyByPrefix(consts.RedisCommentVideoPrefix, message.CommentId), message.VideoId).Result()
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
			_, err = cvr.Rdb.SAdd(context.Background(), dao.GetRedisKeyByPrefix(consts.RedisCommentVideoPrefix, message.CommentId), message.VideoId).Result()
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
