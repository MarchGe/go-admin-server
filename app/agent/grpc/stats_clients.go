package grpc

import (
	"github.com/MarchGe/go-admin-server/agent/config"
	"github.com/MarchGe/go-admin-server/agent/grpc/pb/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"strconv"
	"sync"
)

var sysStatsServiceClient service.SysStatsServiceClient
var once sync.Once

func Initialize(c *config.Config) {
	once.Do(func() {
		conn, err := grpc.NewClient(c.Host+":"+strconv.Itoa(c.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			slog.Error("connect to grpc server error,", err)
			return
		}
		slog.Info("connected to server on addr: " + c.Host + ":" + strconv.Itoa(c.Port))
		sysStatsServiceClient = service.NewSysStatsServiceClient(conn)
	})
}

func GetSysStatsServiceClient() service.SysStatsServiceClient {
	return sysStatsServiceClient
}
