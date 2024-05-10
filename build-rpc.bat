REM protoc version v25.0
protoc --proto_path=C:/Workspace/git/GO-ADMIN/go-admin-server/proto ^
--go_out=rpc_build/go --go-grpc_out=rpc_build/go ^
grpc/model/sys_stats.proto ^
grpc/model/host_info.proto ^
grpc/service/sys_stats_service.proto
