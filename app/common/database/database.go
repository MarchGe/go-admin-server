package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/MarchGe/go-admin-server/config"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"log/slog"
	"os"
	"time"
)

var mdb *gorm.DB
var rdb *redis.Client

var mongodb *mongo.Database

func GetMysql() *gorm.DB {
	return mdb
}

func GetRedis() *redis.Client {
	return rdb
}

func GetMongodb() *mongo.Database {
	return mongodb
}

type CloseFunc func()

func InitializeMysql(c *config.MysqlConfig) CloseFunc {
	var dsn = fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8mb4&parseTime=true&loc=Local&collation=utf8mb4_0900_ai_ci", c.Username, c.Password, c.Host, c.Port, c.Database)
	gormConfig := &gorm.Config{}
	if c.ShowSql {
		gormConfig.Logger = gormLogger.Default.LogMode(gormLogger.Info)
	} else {
		gormConfig.Logger = gormLogger.Default.LogMode(gormLogger.Error)
	}
	var err error
	mdb, err = gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		slog.Error("connect mysql server error: ", slog.Any("err", err))
		os.Exit(1)
	}

	var sqlDB *sql.DB
	sqlDB, err = mdb.DB()
	if err != nil {
		slog.Error("get sql.DB error: ", slog.Any("err", err))
		os.Exit(1)
	}

	sqlDB.SetMaxOpenConns(c.ConnPool.MaxOpenConns)
	sqlDB.SetMaxIdleConns(c.ConnPool.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(c.ConnPool.ConnMaxLifetime))
	sqlDB.SetConnMaxIdleTime(time.Second * time.Duration(c.ConnPool.ConnMaxIdleTime))

	if c.AutoMigrate {
		err = TableAutoMigrate(mdb)
		if err != nil {
			slog.Error("mysql table autoMigrate error: ", slog.Any("err", err))
			os.Exit(1)
		}
	}

	return func() {
		if err = sqlDB.Close(); err != nil {
			slog.Error("close mysql connection error", slog.Any("err", err))
		} else {
			slog.Debug("mysql connection closed success.")
		}
	}
}

func InitializeRedis(c *config.RedisConfig) CloseFunc {
	if !c.Enable {
		return func() {}
	}
	rdb = redis.NewClient(&redis.Options{ // 其他配置可以查看redis.Options结构体，包括连接池配置（默认已经具备连接池功能）、tls配置、slave库的只读配置等。
		Addr:     c.Addr,
		Password: c.Password,
		DB:       c.DB,
	})
	return func() {
		if err := rdb.Close(); err != nil {
			slog.Error("close redis connection error", slog.Any("err", err))
		} else {
			slog.Debug("redis connection closed success.")
		}
	}
}

func InitializeMongo(c *config.MongoConfig) CloseFunc {
	if !c.Enable {
		return func() {}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	loggerOptions := &options.LoggerOptions{}
	var logLevel options.LogLevel
	if c.ShowCommandLog {
		logLevel = options.LogLevelDebug
	} else {
		logLevel = options.LogLevelInfo
	}
	loggerOptions.ComponentLevels = map[options.LogComponent]options.LogLevel{
		options.LogComponentCommand: logLevel,
	}
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(c.Uri).SetLoggerOptions(loggerOptions))
	if err != nil {
		slog.Error("connect mongodb server error: ", slog.Any("err", err))
		os.Exit(1)
	}
	mongodb = mongoClient.Database(c.Database) // 第二个参数&options.DatabaseOptions{}可以做一些相关配置，包括配置读关注、写关注等。

	return func() {
		if err = mongoClient.Disconnect(ctx); err != nil {
			slog.Error("close Mongodb connection error", slog.Any("err", err))
		} else {
			slog.Debug("mongodb connection closed success.")
		}
	}
}
