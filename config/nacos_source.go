package config

import (
	"fmt"
	"github.com/MarchGe/go-admin-server/app/common/nacos"
	"github.com/spf13/viper"
	"strings"
)

var _ Source = (*NacosSource)(nil)

type NacosSource struct {
	Nacoser *nacos.Nacoser
}

func (s *NacosSource) GetConfig() (*Config, error) {
	sConfig, formatType := s.Nacoser.GetConfig()
	return s.parseConfigFromString(sConfig, formatType)
}

func (s *NacosSource) Close() error {
	s.Nacoser.Close()
	return nil
}

func (s *NacosSource) parseConfigFromString(sConfig, sType string) (*Config, error) {
	c := &Config{}
	reader := strings.NewReader(sConfig)
	viper.SetConfigType(sType)
	if err := viper.ReadConfig(reader); err != nil {
		return nil, fmt.Errorf("viper read config from byte buffer error: %w", err)
	}
	if err := viper.Unmarshal(c); err != nil {
		return nil, fmt.Errorf("unmarshal config error: %w", err)
	}
	return c, nil
}
