// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v4.25.0
// source: grpc/service/sys_stats_service.proto

package service

import (
	model "github.com/MarchGe/go-admin-server/app/admin/grpc/pb/model"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_grpc_service_sys_stats_service_proto protoreflect.FileDescriptor

var file_grpc_service_sys_stats_service_proto_rawDesc = []byte{
	0x0a, 0x24, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x73,
	0x79, 0x73, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x73, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x1a,
	0x1a, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2f, 0x73, 0x79, 0x73, 0x5f,
	0x73, 0x74, 0x61, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1a, 0x67, 0x72, 0x70,
	0x63, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2f, 0x68, 0x6f, 0x73, 0x74, 0x5f, 0x69, 0x6e, 0x66,
	0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x32, 0x91, 0x01, 0x0a, 0x0f, 0x53, 0x79, 0x73, 0x53, 0x74, 0x61, 0x74,
	0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x3c, 0x0a, 0x11, 0x72, 0x65, 0x70, 0x6f,
	0x72, 0x74, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x53, 0x74, 0x61, 0x74, 0x73, 0x12, 0x0f, 0x2e,
	0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x53, 0x79, 0x73, 0x53, 0x74, 0x61, 0x74, 0x73, 0x1a, 0x16,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x40, 0x0a, 0x15, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74,
	0x48, 0x6f, 0x73, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x0f, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x48, 0x6f, 0x73, 0x74, 0x49, 0x6e, 0x66, 0x6f,
	0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x42, 0x3e, 0x5a, 0x3c, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x4d, 0x61, 0x72, 0x63, 0x68, 0x47, 0x65, 0x2f, 0x67,
	0x6f, 0x2d, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x61,
	0x70, 0x70, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x70, 0x62,
	0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_grpc_service_sys_stats_service_proto_goTypes = []interface{}{
	(*model.SysStats)(nil), // 0: model.SysStats
	(*model.HostInfo)(nil), // 1: model.HostInfo
	(*emptypb.Empty)(nil),  // 2: google.protobuf.Empty
}
var file_grpc_service_sys_stats_service_proto_depIdxs = []int32{
	0, // 0: service.SysStatsService.reportSystemStats:input_type -> model.SysStats
	1, // 1: service.SysStatsService.reportHostInformation:input_type -> model.HostInfo
	2, // 2: service.SysStatsService.reportSystemStats:output_type -> google.protobuf.Empty
	2, // 3: service.SysStatsService.reportHostInformation:output_type -> google.protobuf.Empty
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_grpc_service_sys_stats_service_proto_init() }
func file_grpc_service_sys_stats_service_proto_init() {
	if File_grpc_service_sys_stats_service_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_grpc_service_sys_stats_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_grpc_service_sys_stats_service_proto_goTypes,
		DependencyIndexes: file_grpc_service_sys_stats_service_proto_depIdxs,
	}.Build()
	File_grpc_service_sys_stats_service_proto = out.File
	file_grpc_service_sys_stats_service_proto_rawDesc = nil
	file_grpc_service_sys_stats_service_proto_goTypes = nil
	file_grpc_service_sys_stats_service_proto_depIdxs = nil
}
