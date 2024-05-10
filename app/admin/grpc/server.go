package grpc

import (
	"context"
	pbService "github.com/MarchGe/go-admin-server/app/admin/grpc/pb/service"
	"github.com/MarchGe/go-admin-server/app/admin/grpc/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"net"
	"os"
)

// TODO grpc服务集成到微服务中...
func Run(ctx context.Context, address string) {
	slog.Info("starting grpc server on address: " + address)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		slog.Error("net listen on address fail", slog.Any("err", err))
		os.Exit(1)
	}
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(addUnaryInterceptors()...))
	registerGrpcService(server)

	go func() {
		if err = server.Serve(listener); err != nil {
			slog.Error("grpc server run fail", slog.Any("err", err))
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	server.GracefulStop()
	slog.Info("grpc server shutdown success.")
}

func addUnaryInterceptors() []grpc.UnaryServerInterceptor {
	interceptors := make([]grpc.UnaryServerInterceptor, 1)
	interceptors[0] = globalErrHandle()
	return interceptors
}

func globalErrHandle() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				resp = nil
				err = status.Errorf(codes.Internal, "RPC服务出错了! err: %v", r)
			}
		}()
		resp, err = handler(ctx, req)
		return resp, err
	}
}

func registerGrpcService(register grpc.ServiceRegistrar) {
	pbService.RegisterSysStatsServiceServer(register, &service.SysStatsServiceGrpc{})
}
