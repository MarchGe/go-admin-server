// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v4.25.0
// source: grpc/model/host_info.proto

package model

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type HostInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ip              string                 `protobuf:"bytes,1,opt,name=ip,proto3" json:"ip,omitempty"`
	HostName        string                 `protobuf:"bytes,2,opt,name=hostName,proto3" json:"hostName,omitempty"`
	UpTime          string                 `protobuf:"bytes,3,opt,name=upTime,proto3" json:"upTime,omitempty"`
	Platform        string                 `protobuf:"bytes,4,opt,name=platform,proto3" json:"platform,omitempty"`
	PlatformVersion string                 `protobuf:"bytes,5,opt,name=platformVersion,proto3" json:"platformVersion,omitempty"`
	KernelVersion   string                 `protobuf:"bytes,6,opt,name=kernelVersion,proto3" json:"kernelVersion,omitempty"`
	KernelArch      string                 `protobuf:"bytes,7,opt,name=kernelArch,proto3" json:"kernelArch,omitempty"`
	CpuInfos        []*CpuInfo             `protobuf:"bytes,8,rep,name=cpuInfos,proto3" json:"cpuInfos,omitempty"`
	Timestamp       *timestamppb.Timestamp `protobuf:"bytes,9,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

func (x *HostInfo) Reset() {
	*x = HostInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_model_host_info_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HostInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HostInfo) ProtoMessage() {}

func (x *HostInfo) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_model_host_info_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HostInfo.ProtoReflect.Descriptor instead.
func (*HostInfo) Descriptor() ([]byte, []int) {
	return file_grpc_model_host_info_proto_rawDescGZIP(), []int{0}
}

func (x *HostInfo) GetIp() string {
	if x != nil {
		return x.Ip
	}
	return ""
}

func (x *HostInfo) GetHostName() string {
	if x != nil {
		return x.HostName
	}
	return ""
}

func (x *HostInfo) GetUpTime() string {
	if x != nil {
		return x.UpTime
	}
	return ""
}

func (x *HostInfo) GetPlatform() string {
	if x != nil {
		return x.Platform
	}
	return ""
}

func (x *HostInfo) GetPlatformVersion() string {
	if x != nil {
		return x.PlatformVersion
	}
	return ""
}

func (x *HostInfo) GetKernelVersion() string {
	if x != nil {
		return x.KernelVersion
	}
	return ""
}

func (x *HostInfo) GetKernelArch() string {
	if x != nil {
		return x.KernelArch
	}
	return ""
}

func (x *HostInfo) GetCpuInfos() []*CpuInfo {
	if x != nil {
		return x.CpuInfos
	}
	return nil
}

func (x *HostInfo) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

type CpuInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Num        int32   `protobuf:"varint,1,opt,name=num,proto3" json:"num,omitempty"`
	VendorId   string  `protobuf:"bytes,2,opt,name=vendorId,proto3" json:"vendorId,omitempty"`
	Family     string  `protobuf:"bytes,3,opt,name=family,proto3" json:"family,omitempty"`
	PhysicalId string  `protobuf:"bytes,4,opt,name=physicalId,proto3" json:"physicalId,omitempty"`
	Cores      int32   `protobuf:"varint,5,opt,name=cores,proto3" json:"cores,omitempty"`
	ModelName  string  `protobuf:"bytes,6,opt,name=modelName,proto3" json:"modelName,omitempty"`
	Mhz        float32 `protobuf:"fixed32,7,opt,name=mhz,proto3" json:"mhz,omitempty"`
}

func (x *CpuInfo) Reset() {
	*x = CpuInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_model_host_info_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CpuInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CpuInfo) ProtoMessage() {}

func (x *CpuInfo) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_model_host_info_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CpuInfo.ProtoReflect.Descriptor instead.
func (*CpuInfo) Descriptor() ([]byte, []int) {
	return file_grpc_model_host_info_proto_rawDescGZIP(), []int{1}
}

func (x *CpuInfo) GetNum() int32 {
	if x != nil {
		return x.Num
	}
	return 0
}

func (x *CpuInfo) GetVendorId() string {
	if x != nil {
		return x.VendorId
	}
	return ""
}

func (x *CpuInfo) GetFamily() string {
	if x != nil {
		return x.Family
	}
	return ""
}

func (x *CpuInfo) GetPhysicalId() string {
	if x != nil {
		return x.PhysicalId
	}
	return ""
}

func (x *CpuInfo) GetCores() int32 {
	if x != nil {
		return x.Cores
	}
	return 0
}

func (x *CpuInfo) GetModelName() string {
	if x != nil {
		return x.ModelName
	}
	return ""
}

func (x *CpuInfo) GetMhz() float32 {
	if x != nil {
		return x.Mhz
	}
	return 0
}

var File_grpc_model_host_info_proto protoreflect.FileDescriptor

var file_grpc_model_host_info_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2f, 0x68, 0x6f, 0x73,
	0x74, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x6d, 0x6f,
	0x64, 0x65, 0x6c, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc0, 0x02, 0x0a, 0x08, 0x48, 0x6f, 0x73, 0x74, 0x49, 0x6e, 0x66,
	0x6f, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
	0x70, 0x12, 0x1a, 0x0a, 0x08, 0x68, 0x6f, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x68, 0x6f, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a,
	0x06, 0x75, 0x70, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75,
	0x70, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72,
	0x6d, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72,
	0x6d, 0x12, 0x28, 0x0a, 0x0f, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x56, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x70, 0x6c, 0x61, 0x74,
	0x66, 0x6f, 0x72, 0x6d, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x24, 0x0a, 0x0d, 0x6b,
	0x65, 0x72, 0x6e, 0x65, 0x6c, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0d, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x12, 0x1e, 0x0a, 0x0a, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x41, 0x72, 0x63, 0x68, 0x18,
	0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x41, 0x72, 0x63,
	0x68, 0x12, 0x2a, 0x0a, 0x08, 0x63, 0x70, 0x75, 0x49, 0x6e, 0x66, 0x6f, 0x73, 0x18, 0x08, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x43, 0x70, 0x75, 0x49,
	0x6e, 0x66, 0x6f, 0x52, 0x08, 0x63, 0x70, 0x75, 0x49, 0x6e, 0x66, 0x6f, 0x73, 0x12, 0x38, 0x0a,
	0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x22, 0xb5, 0x01, 0x0a, 0x07, 0x43, 0x70, 0x75, 0x49,
	0x6e, 0x66, 0x6f, 0x12, 0x10, 0x0a, 0x03, 0x6e, 0x75, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x03, 0x6e, 0x75, 0x6d, 0x12, 0x1a, 0x0a, 0x08, 0x76, 0x65, 0x6e, 0x64, 0x6f, 0x72, 0x49,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x76, 0x65, 0x6e, 0x64, 0x6f, 0x72, 0x49,
	0x64, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x61, 0x6d, 0x69, 0x6c, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x66, 0x61, 0x6d, 0x69, 0x6c, 0x79, 0x12, 0x1e, 0x0a, 0x0a, 0x70, 0x68, 0x79,
	0x73, 0x69, 0x63, 0x61, 0x6c, 0x49, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x70,
	0x68, 0x79, 0x73, 0x69, 0x63, 0x61, 0x6c, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x72,
	0x65, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x63, 0x6f, 0x72, 0x65, 0x73, 0x12,
	0x1c, 0x0a, 0x09, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x10, 0x0a,
	0x03, 0x6d, 0x68, 0x7a, 0x18, 0x07, 0x20, 0x01, 0x28, 0x02, 0x52, 0x03, 0x6d, 0x68, 0x7a, 0x42,
	0x3c, 0x5a, 0x3a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x4d, 0x61,
	0x72, 0x63, 0x68, 0x47, 0x65, 0x2f, 0x67, 0x6f, 0x2d, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2d, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x61, 0x70, 0x70, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f,
	0x67, 0x72, 0x70, 0x63, 0x2f, 0x70, 0x62, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_grpc_model_host_info_proto_rawDescOnce sync.Once
	file_grpc_model_host_info_proto_rawDescData = file_grpc_model_host_info_proto_rawDesc
)

func file_grpc_model_host_info_proto_rawDescGZIP() []byte {
	file_grpc_model_host_info_proto_rawDescOnce.Do(func() {
		file_grpc_model_host_info_proto_rawDescData = protoimpl.X.CompressGZIP(file_grpc_model_host_info_proto_rawDescData)
	})
	return file_grpc_model_host_info_proto_rawDescData
}

var file_grpc_model_host_info_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_grpc_model_host_info_proto_goTypes = []interface{}{
	(*HostInfo)(nil),              // 0: model.HostInfo
	(*CpuInfo)(nil),               // 1: model.CpuInfo
	(*timestamppb.Timestamp)(nil), // 2: google.protobuf.Timestamp
}
var file_grpc_model_host_info_proto_depIdxs = []int32{
	1, // 0: model.HostInfo.cpuInfos:type_name -> model.CpuInfo
	2, // 1: model.HostInfo.timestamp:type_name -> google.protobuf.Timestamp
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_grpc_model_host_info_proto_init() }
func file_grpc_model_host_info_proto_init() {
	if File_grpc_model_host_info_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_grpc_model_host_info_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HostInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_grpc_model_host_info_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CpuInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_grpc_model_host_info_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_grpc_model_host_info_proto_goTypes,
		DependencyIndexes: file_grpc_model_host_info_proto_depIdxs,
		MessageInfos:      file_grpc_model_host_info_proto_msgTypes,
	}.Build()
	File_grpc_model_host_info_proto = out.File
	file_grpc_model_host_info_proto_rawDesc = nil
	file_grpc_model_host_info_proto_goTypes = nil
	file_grpc_model_host_info_proto_depIdxs = nil
}
