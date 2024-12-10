package config

type Source interface {
	GetConfig() (*Config, error)
	Close() error
}
