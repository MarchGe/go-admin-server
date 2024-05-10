package nacos

type Config struct { // 其他配置信息参见nacos-sdk-go的constant.ClientConfig
	TimeoutMs           uint64         `mapstructure:"timeoutMs"` // 连接nacos服务器的超时时间，默认10000ms
	Username            string         `mapstructure:"username"`
	Password            string         `mapstructure:"password"`
	NamespaceId         string         `mapstructure:"namespaceId"`
	DataId              string         `mapstructure:"dataId"`
	Group               string         `mapstructure:"group"`
	Type                string         `mapstructure:"type"` // nacos存储配置信息的格式
	Tag                 string         `mapstructure:"tag"`
	Servers             []ServerConfig `mapstructure:"servers"`
	ClusterName         string         `mapstructure:"clusterName"`
	CacheDir            string         `mapstructure:"cacheDir"`            // 缓存nacos服务信息的目录，默认当前目录
	NotLoadCacheAtStart bool           `mapstructure:"notLoadCacheAtStart"` // 应用每次启动时，是否从缓存中加载nacos服务信息，true：不从缓存中加载
	LogLevel            string         `mapstructure:"logLevel"`            // debug | info | warn | error，默认：info
	LogDir              string         `mapstructure:"logDir"`              // 默认当前目录
	BeatInterval        uint64         `mapstructure:"beatInterval"`        // 与nacos服务器保持心跳的时间间隔，默认5000ms
	ServiceInfo         ServiceInfo    `mapstructure:"serviceInfo"`         // 服务注册
}

type ServerConfig struct {
	IpAddr   string `mapstructure:"ipAddr"`
	Port     uint64 `mapstructure:"port"`
	GrpcPort uint64 `mapstructure:"grpcPort"` // 通常是port + 1000
}

type ServiceInfo struct {
	ServiceName string  `mapstructure:"serviceName"`
	Weight      float64 `mapstructure:"weight"`    // 权重
	Enable      bool    `mapstructure:"enable"`    // 服务实例是否可用
	Healthy     bool    `mapstructure:"healthy"`   // 服务实例是否健康
	Ephemeral   bool    `mapstructure:"ephemeral"` // 是否为临时实例，实例下线时自动删除
}
