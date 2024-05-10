package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/MarchGe/go-admin-server/app"
	_ "github.com/MarchGe/go-admin-server/app/admin/apis/routes"
	_ "github.com/MarchGe/go-admin-server/app/admin/apis/routes/dvroutes"
	"github.com/MarchGe/go-admin-server/app/admin/grpc"
	"github.com/MarchGe/go-admin-server/app/admin/model"
	"github.com/MarchGe/go-admin-server/app/admin/service"
	"github.com/MarchGe/go-admin-server/app/common/database"
	"github.com/MarchGe/go-admin-server/app/common/middleware"
	"github.com/MarchGe/go-admin-server/app/common/nacos"
	"github.com/MarchGe/go-admin-server/app/common/rabbitmq"
	"github.com/MarchGe/go-admin-server/app/common/utils"
	_ "github.com/MarchGe/go-admin-server/app/demo"
	"github.com/MarchGe/go-admin-server/app/demo/mq"
	"github.com/MarchGe/go-admin-server/config"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"log"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var configFile string
var nacosConfigFile string
var cfg = &config.Config{}

var Server = &cobra.Command{
	Use:   "server",
	Short: "start go-admin-server server.",
	Long:  "This command is the bootstrap of go-admin-server.",
	PreRun: func(cmd *cobra.Command, args []string) {
		if configFile != "" && nacosConfigFile != "" {
			log.Fatal("command options --config(or -c)) and --nacosConfig(or -C) cannot be specified together.")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		loadConfig()

		if cfg.Log.Level == "" {
			cfg.Log.Level = "info"
		}
		logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: parseSlogLevel(cfg.Log.Level),
		})
		slog.SetDefault(slog.New(logHandler))

		closeCache := utils.InitializeCache()
		defer closeCache()
		closeMysql := database.InitializeMysql(&cfg.Mysql)
		defer closeMysql()
		closeRedis := database.InitializeRedis(&cfg.Redis)
		defer closeRedis()
		closeMongo := database.InitializeMongo(&cfg.Mongo)
		defer closeMongo()

		if cfg.RabbitMQ.Enable {
			producer, err := createRabbit(cfg.RabbitMQ)
			if err != nil {
				slog.Error("-", err)
				os.Exit(1)
			}
			defer producer.Close()
			if err = mq.InitExchangeAndQueues(producer); err != nil {
				slog.Error("-", slog.Any("err", err))
				os.Exit(1)
			}

			consumer, err := createRabbit(cfg.RabbitMQ)
			if err != nil {
				slog.Error("-", err)
				os.Exit(1)
			}
			defer consumer.Close()
			go mq.RunConsumer(consumer)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if cfg.Grpc.Enable {
			go func() {
				defer func() {
					if r := recover(); r != nil {
						slog.Error("grpc error", slog.Any("err", r))
					}
				}()
				grpc.Run(ctx, cfg.Grpc.Addr)
			}()
		}

		middleware.Initialize(cfg.ContextPath)
		go runWebServer(ctx)

		if cfg.Pprof.Enable {
			go runPprofAnalysis()
		}

		var shutdown = make(chan os.Signal, 1)
		defer close(shutdown)
		signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
		<-shutdown
	},
}

func runPprofAnalysis() {
	slog.Info(fmt.Sprintf("pprof: listening on address :%d", cfg.Pprof.Port))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Pprof.Port), nil); err != nil {
		slog.Error(fmt.Sprintf("pprof: listening on address :%d error", cfg.Pprof.Port), slog.Any("err", err))
	}
}

func parseSlogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		slog.Error("unknown log level: " + level)
		return slog.LevelInfo
	}
}

func loadConfig() {
	var nacoser *nacos.Nacoser
	if nacosConfigFile != "" {
		nc := getNacosConfigFromFile(nacosConfigFile)
		nacoser = nacos.CreateNacoser(nc)
		cfg = parseConfigFromString(nacoser.GetConfig(), nc.Type)
		err := nacoser.RegisterService(cfg.Listen)
		if err != nil {
			log.Panicf("register service to nacos error: %v", err)
		}
	} else {
		cfg = loadConfigFromFile()
	}
	config.Setup(cfg)
}

func init() {
	Server.PersistentFlags().StringVarP(&configFile, "config", "c", "", "Specify config file.(This cannot be specified together with --nacosConfig(or -C))")
	Server.PersistentFlags().StringVarP(&nacosConfigFile, "nacosConfig", "C", "", "Specify nacos config file.(This cannot be specified together with --config(or -c))")
}

func createRabbit(c rabbitmq.Config) (*rabbitmq.Rabbit, error) {
	rabbit, err := rabbitmq.CreateRabbit(c)
	if err != nil {
		return nil, fmt.Errorf("create rabbit error, %w", err)
	}
	return rabbit, nil
}

func loadConfigFromFile() *config.Config {
	c := &config.Config{}
	if configFile == "" {
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
	} else {
		viper.SetConfigFile(configFile)
	}
	if err := viper.ReadInConfig(); err != nil {
		log.Panicf("read config file error: %v", err)
	}
	if err := viper.Unmarshal(c); err != nil {
		log.Panicf("unmarshal config error: %v", err)
	}
	return c
}

func parseConfigFromString(sConfig, sType string) *config.Config {
	c := &config.Config{}
	reader := strings.NewReader(sConfig)
	viper.SetConfigType(sType)
	err := viper.ReadConfig(reader)
	if err != nil {
		log.Panicf("viper read config from byte buffer error: %v", err)
	}
	if err = viper.Unmarshal(c); err != nil {
		log.Panicf("unmarshal config error: %v", err)
	}
	return c
}

func getNacosConfigFromFile(nacosConfigFile string) *nacos.Config {
	c := &nacos.Config{}
	viper.SetConfigFile(nacosConfigFile)
	if err := viper.ReadInConfig(); err != nil {
		log.Panicf("read config file error: %v", err)
	}
	if err := viper.Unmarshal(c); err != nil {
		log.Panicf("unmarshal config error: %v", err)
	}
	return c
}

var server *http.Server

func runWebServer(ctx context.Context) {
	engine := getEngine()
	server = &http.Server{
		Addr:    cfg.Listen,
		Handler: engine,
	}
	initRootUser()
	go func() {
		slog.Info("Listening and serving HTTP on " + cfg.Listen)
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("http server run error: ", slog.Any("err", err))
			os.Exit(1)
		}
	}()
	<-ctx.Done()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("http server shutdown error", slog.Any("err", err))
	}
	slog.Info("http server shutdown success.")
}

const (
	initRootEmail    = "root@example.com"
	initRootPassword = "123456"
)

func initRootUser() {
	db := database.GetMysql()
	u := &model.User{}
	err := db.Where("root = 1").First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = db.Transaction(func(tx *gorm.DB) error {
				user, e := createRootUser(tx)
				if e != nil {
					return e
				}
				if e = createRootPassword(tx, user.Id); e != nil {
					return e
				}
				return nil
			})
			if err != nil {
				slog.Error("init root user error", slog.Any("err", err))
				os.Exit(1)
			}
			slog.Info("*************************************")
			slog.Info("** The initial root user is created: ")
			slog.Info("** Email: " + initRootEmail)
			slog.Info("** Password: " + initRootPassword)
			slog.Info("** WARNING: This information only displayed on first startup, ")
			slog.Info("** please remember account information and login immediately ")
			slog.Info("** to change root password!")
			slog.Info("*************************************")
		} else {
			slog.Error("find root user error", slog.Any("err", err))
			os.Exit(1)
		}
	}
}

func createRootPassword(tx *gorm.DB, userId int64) error {
	passwordHash := service.GetAuthService().PasswordHash(initRootPassword)
	up := &model.UserPassword{
		UserId:   userId,
		Password: passwordHash,
	}
	return tx.Save(up).Error
}

func createRootUser(tx *gorm.DB) (*model.User, error) {
	u := &model.User{
		Name:     "超级用户",
		Nickname: "root",
		Email:    initRootEmail,
		Root:     true,
		Sex:      model.UserSexMan,
		Status:   model.UserStatusNormal,
		DeptId:   0,
		Base: model.Base{
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		},
	}
	e := tx.Save(u).Error
	if e != nil {
		return nil, e
	}
	return u, nil
}

func getEngine() *gin.Engine {
	if cfg.Environment == config.DEV {
		gin.SetMode(gin.DebugMode)
	} else if cfg.Environment == config.TEST {
		gin.SetMode(gin.TestMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()

	engine.
		Use(middleware.ApiDebugLogger()).
		Use(middleware.GlobalErrHandler()).
		Use(middleware.SetRequestId()).
		Use(middleware.SetSession(&cfg.Cookie)).
		Use(middleware.AuthenticationHandler()).
		Use(middleware.BindingValidateTranslator())

	trustedProxies := strings.TrimSpace(cfg.TrustedProxies)
	if trustedProxies != "" {
		proxies := strings.Split(trustedProxies, ",")
		for i := range proxies {
			proxies[i] = strings.TrimSpace(proxies[i])
		}
		err := engine.SetTrustedProxies(proxies)
		if err != nil {
			panic(err)
		}
	}

	app.InitRoutes(engine.Group(cfg.ContextPath))
	return engine
}
