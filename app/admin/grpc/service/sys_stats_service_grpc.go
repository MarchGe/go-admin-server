package service

import (
	"context"
	"github.com/MarchGe/go-admin-server/app/admin/apis/devops/monitor"
	"github.com/MarchGe/go-admin-server/app/admin/grpc/pb/model"
	pbService "github.com/MarchGe/go-admin-server/app/admin/grpc/pb/service"
	"google.golang.org/protobuf/types/known/emptypb"
	"log/slog"
)

var _ pbService.SysStatsServiceServer = (*SysStatsServiceGrpc)(nil)

type SysStatsServiceGrpc struct {
	pbService.UnimplementedSysStatsServiceServer
}

func (s SysStatsServiceGrpc) ReportHostInformation(ctx context.Context, info *model.HostInfo) (*emptypb.Empty, error) {
	slog.Debug("received rpc message: ", slog.Any("info", info))
	if err := monitor.ProcessHostInformation(info); err != nil {
		slog.Error("-", slog.Any("err", err))
		return &emptypb.Empty{}, nil
	}
	return &emptypb.Empty{}, nil

}

func (s SysStatsServiceGrpc) ReportSystemStats(ctx context.Context, stats *model.SysStats) (*emptypb.Empty, error) {
	slog.Debug("received rpc message: ", slog.Any("stats", stats))
	if err := monitor.ProcessPerformanceStats(stats); err != nil {
		slog.Error("-", slog.Any("err", err))
		return &emptypb.Empty{}, nil
	}
	return &emptypb.Empty{}, nil
}
