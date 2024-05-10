package rabbitmq

type Config struct {
	Enable            bool   `mapstructure:"enable"`
	Url               string `mapstructure:"url"`
	Vhost             string `mapstructure:"vhost"`
	Locale            string `mapstructure:"locale"`
	ReconnectInterval int    `mapstructure:"reconnectInterval"` // 断线重连的频率，单位：秒，默认2秒
}
