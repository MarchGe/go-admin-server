package cmd

import (
	"context"
	"github.com/MarchGe/go-admin-server/agent/config"
	myGrpc "github.com/MarchGe/go-admin-server/agent/grpc"
	"github.com/MarchGe/go-admin-server/agent/scheduler"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var cfg = &config.Config{}
var rootCmd = &cobra.Command{
	Use: "agent",
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		myGrpc.Initialize(cfg)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		tasks := getTasks()
		go scheduler.Run(ctx, tasks)

		shutdown := make(chan os.Signal, 1)
		defer close(shutdown)
		signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
		<-shutdown
	},
	Version: config.Version,
}

func getTasks() []scheduler.Task {
	d, _ := time.ParseDuration(strconv.Itoa(cfg.PerformanceReportFrequency) + "s")
	tasks := make([]scheduler.Task, 1)
	tasks[0] = &scheduler.StatsTask{
		TaskI: scheduler.TaskI{
			Id:   uuid.NewString(),
			Name: "上报性能数据",
			Cron: "@every " + d.String(),
		},
	}
	return tasks
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfg.Host, "host", "H", "", "Specify admin server host")
	rootCmd.PersistentFlags().IntVarP(&cfg.Port, "port", "p", 9080, "Specify admin server port")
	rootCmd.PersistentFlags().IntVarP(&cfg.PerformanceReportFrequency, "performanceReportFrequency", "f", 2, "the frequency in seconds of reporting performance data to admin manage server")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		slog.Error("", err)
		os.Exit(1)
	}
}
