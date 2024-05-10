package mq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/MarchGe/go-admin-server/app/common/rabbitmq"
	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
	"log/slog"
	"strconv"
	"time"
)

var p *rabbitmq.Rabbit

func InitExchangeAndQueues(producer *rabbitmq.Rabbit) error {
	p = producer
	channel, err := producer.CreateChannel()
	if err != nil {
		return fmt.Errorf("rabbit get producer channel error, %w", err)
	}
	defer p.CloseChannel(channel)

	err1 := channel.ExchangeDeclare(AlternateExchange, rabbitmq.Fanout, false, false, false, false, nil)
	err2 := channel.ExchangeDeclare(DLXExchange, rabbitmq.Direct, false, false, false, false, nil)
	args := make(amqp091.Table)
	args["alternate-exchange"] = AlternateExchange
	err3 := channel.ExchangeDeclare(LogExchange, rabbitmq.Direct, false, false, false, false, args)

	qArgs := make(amqp091.Table)
	qArgs["x-dead-letter-exchange"] = DLXExchange
	qArgs["x-dead-letter-routing-key"] = DLXRK
	qArgs["x-message-ttl"] = 10000
	_, err4 := channel.QueueDeclare(LogQueue1, true, false, false, false, qArgs)
	_, err5 := channel.QueueDeclare(LogQueue2, true, false, false, false, nil)
	_, err6 := channel.QueueDeclare(AlternateQueue, true, false, false, false, nil)
	_, err7 := channel.QueueDeclare(DLXQueue, true, false, false, false, nil)

	err8 := channel.QueueBind(LogQueue1, LogRK1, LogExchange, false, nil)
	err9 := channel.QueueBind(LogQueue2, LogRK2, LogExchange, false, nil)
	err10 := channel.QueueBind(AlternateQueue, AlternateRK, AlternateExchange, false, nil)
	err11 := channel.QueueBind(DLXQueue, DLXRK, DLXExchange, false, nil)
	if err = errors.Join(err1, err2, err3, err4, err5, err6, err7, err8, err9, err10, err11); err != nil {
		return fmt.Errorf("init exchange and queues error, %w", err)
	}
	return nil
}

// SendMessageConfirm 使用生产者确认模式发送消息
func SendMessageConfirm(msg *Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	channel, err := p.CreateChannel()
	if err != nil {
		return fmt.Errorf("create rabbitmq channel error: %w", err)
	}
	defer p.CloseChannel(channel)
	if err = channel.Confirm(false); err != nil {
		return fmt.Errorf("set channel to confirm mode error: %w", err)
	}

	returnChan := channel.NotifyReturn(make(chan amqp091.Return))
	go func() {
		v := <-returnChan // channel关闭时，这些通道会被自动关闭
		fmt.Println("return 回来的消息：", v)
	}()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	confirmation, err := channel.PublishWithDeferredConfirmWithContext(ctx, LogExchange, LogRK1, true, false, amqp091.Publishing{
		MessageId:    uuid.New().String(),
		ContentType:  "application/json",
		DeliveryMode: amqp091.Transient,
		Timestamp:    time.Now(),
		Type:         "Log",
		Body:         data,
	})
	if err != nil {
		return fmt.Errorf("publish message to rabbitmq error: %w", err)
	}
	if confirmation.Wait() {
		slog.Debug("rabbitmq receive message succeed.")
		return nil
	} else {
		return fmt.Errorf("rabbitmq receive message failed, deliveryTag=%d", confirmation.DeliveryTag)
	}
}

// SendMessageTx 使用事务模式发送消息，性能较低
func SendMessageTx(msg *Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	channel, err := p.CreateChannel()
	if err != nil {
		return fmt.Errorf("create rabbitmq channel error: %w", err)
	}
	defer p.CloseChannel(channel)

	if err = channel.Tx(); err != nil {
		return fmt.Errorf("set channel to transaction mode error: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = channel.PublishWithContext(ctx, LogExchange, LogRK1, true, false, amqp091.Publishing{
		MessageId:    uuid.New().String(),
		ContentType:  "application/json",
		DeliveryMode: amqp091.Transient,
		Timestamp:    time.Now(),
		Type:         "Log",
		Body:         data,
	})
	if err != nil {
		return fmt.Errorf("publish message to rabbitmq error: %w", err)
	}

	if ok := doSomething(); ok {
		if err = channel.TxCommit(); err != nil {
			return fmt.Errorf("commit message to rabbitmq error: %w", err)
		}
	} else {
		if err = channel.TxRollback(); err != nil {
			return fmt.Errorf("rollback message from rabbitmq error: %w", err)
		}
		return fmt.Errorf("message has been rolled back")
	}
	return nil
}

// doSomething 模式业务逻辑
func doSomething() bool {
	return true
}

// SendMessageTTL 发送TTL消息
func SendMessageTTL(msg *Message, ttlInMillis int) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	channel, err := p.CreateChannel()
	if err != nil {
		return fmt.Errorf("create rabbitmq channel error: %w", err)
	}
	defer p.CloseChannel(channel)
	if err = channel.Confirm(false); err != nil {
		return fmt.Errorf("set channel to confirm mode error: %w", err)
	}

	returnChan := channel.NotifyReturn(make(chan amqp091.Return))
	go func() {
		v := <-returnChan
		fmt.Println("return 回来的消息：", v)
	}()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	confirmation, err := channel.PublishWithDeferredConfirmWithContext(ctx, LogExchange, LogRK1, true, false, amqp091.Publishing{
		MessageId:    uuid.New().String(),
		ContentType:  "application/json",
		DeliveryMode: amqp091.Transient,
		Timestamp:    time.Now(),
		Type:         "Log",
		Body:         data,
		Expiration:   strconv.Itoa(ttlInMillis),
	})
	if err != nil {
		return fmt.Errorf("publish message to rabbitmq error: %w", err)
	}
	if confirmation.Wait() {
		slog.Debug("rabbitmq receive message succeed.")
		return nil
	} else {
		return fmt.Errorf("rabbitmq receive message failed, deliveryTag=%d", confirmation.DeliveryTag)
	}
}
