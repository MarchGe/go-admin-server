package config

import (
	"fmt"
	"github.com/spf13/viper"
)

var _ Source = (*FileSource)(nil)

type FileSource struct {
	FilePath string
}

func (s *FileSource) GetConfig() (*Config, error) {
	c := &Config{}
	if s.FilePath == "" {
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
	} else {
		viper.SetConfigFile(s.FilePath)
	}
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config file error: %w", err)
	}
	if err := viper.Unmarshal(c); err != nil {
		return nil, fmt.Errorf("unmarshal config error: %w", err)
	}
	return c, nil
}

func (s *FileSource) Close() error {
	return nil
}
