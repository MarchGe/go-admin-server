package config

var cfg *Config

type Config struct {
	Host                       string `mapstructure:"host"`
	Port                       int    `mapstructure:"port"`
	PerformanceReportFrequency int    `mapstructure:"performanceReportFrequency"` // 性能数据上报频率，单位：秒
}

func Setup(c *Config) {
	cfg = c
}

func GetConfig() *Config {
	return cfg
}
