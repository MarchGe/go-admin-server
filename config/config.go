package config

import (
	"github.com/MarchGe/go-admin-server/app/common/rabbitmq"
	"github.com/MarchGe/go-admin-server/app/common/utils/email"
	"log"
)

type Config struct {
	Environment          string          `mapstructure:"environment"` // 环境 dev | test | prod
	Listen               string          `mapstructure:"listen"`
	ContextPath          string          `mapstructure:"contextPath"`
	Log                  LogConfig       `mapstructure:"log"`
	EncryptKey           string          `mapstructure:"encryptKey"` // 为了某些数据的安全考虑，使用的对称加密的密钥，16 bytes 或 24 bytes 或 32 bytes
	Cookie               CookieConfig    `mapstructure:"cookie"`
	TrustedProxies       string          `mapstructure:"trustedProxies"`
	Pprof                Pprof           `mapstructure:"pprof"` // pprof性能分析工具，只对开发人员有用
	Mysql                MysqlConfig     `mapstructure:"mysql"`
	Redis                RedisConfig     `mapstructure:"redis"` // 该配置项暂时没用到，可忽略
	Mongo                MongoConfig     `mapstructure:"mongo"` // 该配置项暂时没用到，可忽略
	Email                EmailConfig     `mapstructure:"email"`
	RabbitMQ             rabbitmq.Config `mapstructure:"rabbitmq"` // 该配置项没用到，可忽略
	Grpc                 GrpcConfig      `mapstructure:"grpc"`
	UploadPkgSizeLimit   int32           `mapstructure:"uploadPkgSizeLimit"`   // 上传的应用部署包的大小限制，单位MB
	UploadPkgPath        string          `mapstructure:"uploadPkgPath"`        // 应用部署包上传路径，路径分隔符必须是"/"
	WorkDir              string          `mapstructure:"workDir"`              // 指定工作目录，路径分隔符必须是"/"
	ScriptExecuteTimeout int32           `mapstructure:"scriptExecuteTimeout"` // 执行script脚本时的超时时间（单位：秒），超过该时间，如果脚本还没执行完，ssh连接将自动断开
}

type LogConfig struct {
	Level        string `mapstructure:"level"`        // 日志级别 debug | info | warning | error，全局有效
	StackTrace   bool   `mapstructure:"stackTrace"`   // 是否打印错误堆栈信息，dev环境始终都会向stderr打印错误堆栈，不管该选项开启还是关闭
	LoginLog     bool   `mapstructure:"loginLog"`     // 是否记录登录日志
	OpLog        bool   `mapstructure:"opLog"`        // 是否记录操作日志
	ExceptionLog bool   `mapstructure:"exceptionLog"` // 是否记录异常日志
}

type CookieConfig struct {
	AuthenticationKey string `mapstructure:"authenticationKey"` // 32 bytes 或 64 bytes
	SecretKey         string `mapstructure:"secretKey"`         // 16 bytes 或 24 bytes 或 32 bytes
	Path              string `mapstructure:"path"`
	MaxAge            int    `mapstructure:"maxAge"`
	Secure            bool   `mapstructure:"secure"`
	HttpOnly          bool   `mapstructure:"httpOnly"`
}

type Pprof struct {
	Enable bool  `mapstructure:"enable"` // 是否开启pprof性能分析
	Port   int32 `mapstructure:"port"`   // pprof性能分析服务监听的端口
}
type MysqlConfig struct {
	Host     string `mapstructure:"host"`
	Port     int32  `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	ConnPool struct {
		MaxOpenConns    int `mapstructure:"maxOpenConns"`
		MaxIdleConns    int `mapstructure:"maxIdleConns"`
		ConnMaxLifetime int `mapstructure:"connMaxLifetime"` // 单位：秒
		ConnMaxIdleTime int `mapstructure:"connMaxIdleTime"` // 单位：秒
	} `mapstructure:"connPool"`
	ShowSql     bool `mapstructure:"showSql"`
	AutoMigrate bool `mapstructure:"autoMigrate"`
}
type RedisConfig struct {
	Enable   bool   `mapstructure:"enable"`
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type MongoConfig struct {
	Enable         bool   `mapstructure:"enable"`
	Uri            string `mapstructure:"uri"`
	Database       string `mapstructure:"database"`
	ShowCommandLog bool   `mapstructure:"showCommandLog"` // 是否打印出mongo语句
}

type EmailConfig struct {
	email.MailConfig `mapstructure:"mc"`
	SystemName       string `mapstructure:"systemName"` // 系统名称，发送邮件时会使用
	AccessUrl        string `mapstructure:"accessUrl"`  // 系统访问的url，发送邮件时会使用
}

type GrpcConfig struct {
	Enable            bool   `mapstructure:"enable"`
	Addr              string `mapstructure:"addr"`
	ConnectionTimeout int    `mapstructure:"connectionTimeout"` // 单位：秒
}

const (
	DEV  = "dev"
	PROD = "prod"
	TEST = "test"
)

var cfg *Config

func Setup(c *Config) {
	verifyCfg(c)
	cfg = c
}

func GetConfig() *Config {
	return cfg
}

func verifyCfg(cfg *Config) {
	if cfg.Environment != DEV && cfg.Environment != PROD && cfg.Environment != TEST {
		log.Panicf("The config value of 'environment' is incorrect: %v", cfg.Environment)
	}
}
