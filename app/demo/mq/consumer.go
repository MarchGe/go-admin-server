package mq

import (
	"fmt"
	"github.com/MarchGe/go-admin-server/app/common/rabbitmq"
	"github.com/rabbitmq/amqp091-go"
	"log/slog"
	"time"
)

func RunConsumer(consumer *rabbitmq.Rabbit) {
	var channel *amqp091.Channel
	var err error
	for {
		channel, err = consumer.CreateChannel()
		if err != nil {
			slog.Error("rabbit get consumer channel error", slog.Any("err", err))
			reCreateInterval := rabbitmq.DefaultReconnectInterval
			if consumer.GetConfig().ReconnectInterval > 0 {
				reCreateInterval = consumer.GetConfig().ReconnectInterval
			}
			time.Sleep(time.Duration(reCreateInterval) * time.Second)
			continue
		}
		break
	}
	defer consumer.CloseChannel(channel)
	consumer.SetChannelCloseListener(channel, func() {
		RunConsumer(consumer)
	})
	err = channel.Qos(10, 0, false)
	if err != nil {
		slog.Error("set consumer channel qos error", slog.Any("err", err))
		return
	}
	consume(channel)
}

func consume(channel *amqp091.Channel) {
	deliveries, err := channel.Consume(LogQueue1, LogConsumer, false, false, true, false, nil)
	if err != nil {
		slog.Error("consume message error: ", slog.Any("err", err))
		return
	}
	for delivery := range deliveries {
		fmt.Println("RabbitMQ收到的消息为：", string(delivery.Body))
		if err = delivery.Ack(false); err != nil {
			slog.Error("ack message error: ", slog.Any("err", err))
			return
		}
	}
}
