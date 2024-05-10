package rabbitmq

import (
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"log/slog"
	"sync"
	"time"
)

const DefaultReconnectInterval = 2

type Rabbit struct {
	conn *amqp091.Connection
	cfg  *Config
}

func CreateRabbit(c Config) (*Rabbit, error) {
	r := &Rabbit{
		cfg: &c,
	}
	r.connect()
	return r, nil
}

func (r *Rabbit) GetConfig() *Config {
	return r.cfg
}

func newConnection(c *Config) (*amqp091.Connection, error) {
	amqpCfg := amqp091.Config{
		Vhost:  c.Vhost,
		Locale: c.Locale,
	}
	connection, err := amqp091.DialConfig(c.Url, amqpCfg)
	if err != nil {
		return nil, fmt.Errorf("dial rabbitmq server error: %w", err)
	}
	return connection, nil
}

var mutex sync.Mutex

func (r *Rabbit) connect() {
	mutex.Lock()
	defer mutex.Unlock()
	if r.conn != nil && !r.conn.IsClosed() {
		return
	}
	for {
		connection, err := newConnection(r.cfg)
		if err != nil {
			slog.Error("connect rabbitmq server error", slog.Any("err", err))
			if r.cfg.ReconnectInterval <= 0 {
				time.Sleep(DefaultReconnectInterval * time.Second)
			} else {
				time.Sleep(time.Duration(r.cfg.ReconnectInterval) * time.Second)
			}
			continue
		}
		r.conn = connection
		slog.Info("connect rabbitmq server success.")
		break
	}
	closeChan := r.conn.NotifyClose(make(chan *amqp091.Error, 1))
	go func() {
		err := <-closeChan
		if err != nil {
			slog.Error("rabbitmq connection error", slog.Any("err", err))
			r.connect()
		}
	}()
}

func (r *Rabbit) CreateChannel() (*amqp091.Channel, error) {
	channel, err := r.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("open rabbitmq channel error: %w", err)
	}
	return channel, nil
}

type ChannelCloseListener func()

func (r *Rabbit) SetChannelCloseListener(channel *amqp091.Channel, listener ChannelCloseListener) {
	closeChan := channel.NotifyClose(make(chan *amqp091.Error))
	go func() {
		err := <-closeChan
		if err != nil {
			slog.Error("rabbitmq channel error", slog.Any("err", err))
			listener()
		}
	}()
}

func (r *Rabbit) CloseChannel(channel *amqp091.Channel) {
	if channel != nil && !channel.IsClosed() {
		if err := channel.Close(); err != nil {
			slog.Error("close rabbitmq channel failed", slog.Any("err", err))
		} else {
			slog.Debug("close rabbitmq channel success")
		}
	}
}

func (r *Rabbit) Close() {
	if r.conn != nil && !r.conn.IsClosed() {
		if err := r.conn.Close(); err != nil {
			slog.Error("close rabbitmq connection failed", slog.Any("err", err))
		} else {
			slog.Debug("close rabbitmq connection success")
		}
	}
}

func (r *Rabbit) Disconnected() bool {
	return r.conn.IsClosed()
}
