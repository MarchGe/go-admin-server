package scheduler

import (
	"context"
	"errors"
	"fmt"
	grpcService "github.com/MarchGe/go-admin-server/agent/grpc"
	"github.com/MarchGe/go-admin-server/agent/grpc/pb/model"
	"github.com/MarchGe/go-admin-server/agent/grpc/pb/service"
	"github.com/MarchGe/go-admin-server/agent/monitor"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log/slog"
	"time"
)

var _ Task = (*StatsTask)(nil)

type StatsTask struct {
	TaskI
}

func (s *StatsTask) Execute() {
	sysStatsServiceClient := grpcService.GetSysStatsServiceClient()
	if err := reportPerformanceStats(sysStatsServiceClient); err != nil {
		slog.Error("-", slog.Any("err", err))
	}
	if err := reportHostInfo(sysStatsServiceClient); err != nil {
		slog.Error("-", slog.Any("err", err))
	}
}

func reportPerformanceStats(client service.SysStatsServiceClient) error {
	cpuStat, e1 := monitor.GetCpuStat()
	virtualStat, swapStat, e2 := monitor.GetMemoryStat()
	diskStats, e3 := monitor.GetDiskStats()
	if err := errors.Join(e1, e2, e3); err != nil {
		return fmt.Errorf("get system stats error, %w", err)
	}
	sysStats := &model.SysStats{
		Cpu:           cpuStat,
		VirtualMemory: virtualStat,
		SwapMemory:    swapStat,
		Disk:          diskStats,
	}
	ip, err := monitor.GetIP()
	if err != nil {
		return fmt.Errorf("get machine ip error, %w", err)
	}
	sysStats.Ip = ip
	now := time.Now()
	sysStats.Timestamp = &timestamppb.Timestamp{
		Seconds: now.UnixMilli() / 1000,
		Nanos:   int32(now.Nanosecond()),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if _, err = client.ReportSystemStats(ctx, sysStats); err != nil {
		return fmt.Errorf("invoke grpc api 'ReportSystemStats' error, %w", err)
	}
	return nil
}

func reportHostInfo(client service.SysStatsServiceClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	info, err := monitor.GetHostInfo()
	if err != nil {
		return fmt.Errorf("get host information error, %w", err)
	}
	now := time.Now()
	info.Timestamp = &timestamppb.Timestamp{
		Seconds: now.UnixMilli() / 1000,
		Nanos:   int32(now.Nanosecond()),
	}
	if _, err = client.ReportHostInformation(ctx, info); err != nil {
		return fmt.Errorf("invoke grpc api 'ReportHostInformation' error, %w", err)
	}
	return nil
}
