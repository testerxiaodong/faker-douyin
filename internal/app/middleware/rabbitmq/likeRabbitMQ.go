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

type UserLikeVideoRabbitMQ struct {
	RabbitMQ     *RabbitMQ
	QueueName    string
	ExchangeName string
	RoutingKey   string
	Rdb          *redis.Client
}

type VideoLikedByUserRabbitMQ struct {
	RabbitMQ     *RabbitMQ
	QueueName    string
	ExchangeName string
	RoutingKey   string
	Rdb          *redis.Client
}

type LikeMessage struct {
	LikeDealType int64
	UserId       int64
	VideoId      int64
}

func NewUserLikeVideoRabbitMQ(rabbitMQ *RabbitMQ, rdb *redis.Client) *UserLikeVideoRabbitMQ {
	var userLikeVideoRabbitMQ UserLikeVideoRabbitMQ
	userLikeVideoRabbitMQ.RabbitMQ = rabbitMQ
	userLikeVideoRabbitMQ.QueueName = consts.UserLikeVideoQueue
	userLikeVideoRabbitMQ.Rdb = rdb
	go userLikeVideoRabbitMQ.Consume()
	return &userLikeVideoRabbitMQ
}

func NewVideoLikedByUserRabbitMQ(rabbitMQ *RabbitMQ, rdb *redis.Client) *VideoLikedByUserRabbitMQ {
	var videoLikedByUserRabbitMQ VideoLikedByUserRabbitMQ
	videoLikedByUserRabbitMQ.RabbitMQ = rabbitMQ
	videoLikedByUserRabbitMQ.QueueName = consts.VideoLikedByUserQueue
	videoLikedByUserRabbitMQ.Rdb = rdb
	go videoLikedByUserRabbitMQ.Consume()
	return &videoLikedByUserRabbitMQ
}

func (ulvr *UserLikeVideoRabbitMQ) Publish(message LikeMessage) {
	channel, err := ulvr.RabbitMQ.Conn.Channel()
	if err != nil {
		log.AppLogger.Error(err.Error())
	}
	_, err = channel.QueueDeclare(ulvr.QueueName, true, false, false, false, nil)
	if err != nil {
		log.AppLogger.Error(err.Error())
	}
	messageByte, err := json.Marshal(message)
	if err != nil {
		log.AppLogger.Error(err.Error())
	}
	err = channel.Publish(ulvr.ExchangeName, ulvr.QueueName, false, false, amqp.Publishing{
		ContentType:  "text/plain",
		DeliveryMode: amqp.Persistent,
		Body:         messageByte,
	})
	if err != nil {
		log.AppLogger.Error(err.Error())
	}
}

func (ulvr *UserLikeVideoRabbitMQ) Consume() {
	channel, err := ulvr.RabbitMQ.Conn.Channel()
	if err != nil {
		log.AppLogger.Error(err.Error())
	}
	_, err = channel.QueueDeclare(ulvr.QueueName, true, false, false, false, nil)
	if err != nil {
		log.AppLogger.Error(err.Error())
	}
	deliveries, err := channel.Consume(ulvr.QueueName, ulvr.QueueName, false, false, false, false, nil)
	if err != nil {
		log.AppLogger.Error(err.Error())
	}
	for delivery := range deliveries {
		var message *LikeMessage
		err := json.Unmarshal(delivery.Body, message)
		if err != nil {
			log.AppLogger.Error(err.Error())
		}
		// 取消点赞消息补偿
		if message.LikeDealType == 0 {
			result, err := ulvr.Rdb.Exists(context.Background(), dao.GetRedisKeyByPrefix(consts.RedisUserLikeVideoPrefix, message.UserId)).Result()
			if err != nil {
				// 发生错误，Nack并且重新入队
				err = channel.Reject(0, true)
				if err != nil {
					log.AppLogger.Error(err.Error())
				}
			}
			// key存在，进行补偿业务处理，不存在则可能是过期了
			if result > 0 {
				_, err = ulvr.Rdb.SRem(context.Background(), dao.GetRedisKeyByPrefix(consts.RedisUserLikeVideoPrefix, message.UserId), message.VideoId).Result()
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
		} else if message.LikeDealType == 1 {
			_, err = ulvr.Rdb.SAdd(context.Background(), dao.GetRedisKeyByPrefix(consts.RedisUserLikeVideoPrefix, message.UserId), message.VideoId).Result()
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

func (vlur *VideoLikedByUserRabbitMQ) Publish(message LikeMessage) {
	channel, err := vlur.RabbitMQ.Conn.Channel()
	if err != nil {
		log.AppLogger.Error(err.Error())
	}
	_, err = channel.QueueDeclare(vlur.QueueName, true, false, false, false, nil)
	if err != nil {
		log.AppLogger.Error(err.Error())
	}
	messageByte, err := json.Marshal(message)
	if err != nil {
		log.AppLogger.Error(err.Error())
	}
	err = channel.Publish(vlur.ExchangeName, vlur.QueueName, false, false, amqp.Publishing{
		ContentType:  "text/plain",
		DeliveryMode: amqp.Persistent,
		Body:         messageByte,
	})
	if err != nil {
		log.AppLogger.Error(err.Error())
	}
}

func (vlur *VideoLikedByUserRabbitMQ) Consume() {
	channel, err := vlur.RabbitMQ.Conn.Channel()
	if err != nil {
		log.AppLogger.Error(err.Error())
	}
	_, err = channel.QueueDeclare(vlur.QueueName, true, false, false, false, nil)
	if err != nil {
		log.AppLogger.Error(err.Error())
	}
	deliveries, err := channel.Consume(vlur.QueueName, vlur.QueueName, false, false, false, false, nil)
	if err != nil {
		log.AppLogger.Error(err.Error())
	}
	for delivery := range deliveries {
		var message *LikeMessage
		err := json.Unmarshal(delivery.Body, message)
		if err != nil {
			log.AppLogger.Error(err.Error())
		}
		// 取消点赞消息补偿
		if message.LikeDealType == 0 {
			result, err := vlur.Rdb.Exists(context.Background(), dao.GetRedisKeyByPrefix(consts.RedisVideoLikedByUserPrefix, message.VideoId)).Result()
			if err != nil {
				// 发生错误，Nack并且重新入队
				err = channel.Reject(0, true)
				if err != nil {
					log.AppLogger.Error(err.Error())
				}
			}
			// key存在，进行补偿业务处理，不存在则可能是过期了
			if result > 0 {
				_, err = vlur.Rdb.SRem(context.Background(), dao.GetRedisKeyByPrefix(consts.RedisVideoLikedByUserPrefix, message.VideoId), message.UserId).Result()
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
		} else if message.LikeDealType == 1 {
			_, err = vlur.Rdb.SAdd(context.Background(), dao.GetRedisKeyByPrefix(consts.RedisVideoLikedByUserPrefix, message.VideoId), message.UserId).Result()
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
