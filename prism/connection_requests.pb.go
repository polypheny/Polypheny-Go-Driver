//
//Messages related to establishing and maintaining a connection.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v5.26.0
// source: org/polypheny/prism/connection_requests.proto

package prism

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

//
//The ConnectionRequest message is designed to initiate a connection request from the client to the server.
//It contains information regarding the API version, client identity, and optional credentials, as well as properties associated with the connection.
type ConnectionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Specifies the major version of the API that the client is using.
	MajorApiVersion int32 `protobuf:"varint,1,opt,name=major_api_version,json=majorApiVersion,proto3" json:"major_api_version,omitempty"`
	// Represents the minor version of the API in use.
	MinorApiVersion int32 `protobuf:"varint,2,opt,name=minor_api_version,json=minorApiVersion,proto3" json:"minor_api_version,omitempty"`
	// (Optional) The username for authentication when establishing the connection.
	Username *string `protobuf:"bytes,5,opt,name=username,proto3,oneof" json:"username,omitempty"`
	// (Optional) The password associated with the specified username for authentication purposes.
	Password *string `protobuf:"bytes,6,opt,name=password,proto3,oneof" json:"password,omitempty"`
	// (Optional) Contains specific properties related to the connection, such as timeout settings or database preferences.
	ConnectionProperties *ConnectionProperties `protobuf:"bytes,4,opt,name=connection_properties,json=connectionProperties,proto3,oneof" json:"connection_properties,omitempty"`
}

func (x *ConnectionRequest) Reset() {
	*x = ConnectionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_org_polypheny_prism_connection_requests_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConnectionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConnectionRequest) ProtoMessage() {}

func (x *ConnectionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_org_polypheny_prism_connection_requests_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConnectionRequest.ProtoReflect.Descriptor instead.
func (*ConnectionRequest) Descriptor() ([]byte, []int) {
	return file_org_polypheny_prism_connection_requests_proto_rawDescGZIP(), []int{0}
}

func (x *ConnectionRequest) GetMajorApiVersion() int32 {
	if x != nil {
		return x.MajorApiVersion
	}
	return 0
}

func (x *ConnectionRequest) GetMinorApiVersion() int32 {
	if x != nil {
		return x.MinorApiVersion
	}
	return 0
}

func (x *ConnectionRequest) GetUsername() string {
	if x != nil && x.Username != nil {
		return *x.Username
	}
	return ""
}

func (x *ConnectionRequest) GetPassword() string {
	if x != nil && x.Password != nil {
		return *x.Password
	}
	return ""
}

func (x *ConnectionRequest) GetConnectionProperties() *ConnectionProperties {
	if x != nil {
		return x.ConnectionProperties
	}
	return nil
}

//
//The ConnectionProperties message defines specific properties related to the client-server connection.
//It allows clients to specify certain behaviors and settings for the connection, such as transaction auto-commit status and target namespace preference.
type ConnectionPropertiesUpdateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	//
	//Specifies the new settings for the connection.
	//Contains properties such as transaction auto-commit status and target namespace preference.
	ConnectionProperties *ConnectionProperties `protobuf:"bytes,4,opt,name=connection_properties,json=connectionProperties,proto3" json:"connection_properties,omitempty"`
}

func (x *ConnectionPropertiesUpdateRequest) Reset() {
	*x = ConnectionPropertiesUpdateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_org_polypheny_prism_connection_requests_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConnectionPropertiesUpdateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConnectionPropertiesUpdateRequest) ProtoMessage() {}

func (x *ConnectionPropertiesUpdateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_org_polypheny_prism_connection_requests_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConnectionPropertiesUpdateRequest.ProtoReflect.Descriptor instead.
func (*ConnectionPropertiesUpdateRequest) Descriptor() ([]byte, []int) {
	return file_org_polypheny_prism_connection_requests_proto_rawDescGZIP(), []int{1}
}

func (x *ConnectionPropertiesUpdateRequest) GetConnectionProperties() *ConnectionProperties {
	if x != nil {
		return x.ConnectionProperties
	}
	return nil
}

//
//The ConnectionProperties message defines specific properties related to the client-server connection.
//It allows clients to specify certain behaviors and settings for the connection, such as transaction auto-commit status and target namespace preference.
type ConnectionProperties struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// (Optional) Indicates whether transactions are automatically committed upon execution or require manual commitment.
	IsAutoCommit *bool `protobuf:"varint,1,opt,name=is_auto_commit,json=isAutoCommit,proto3,oneof" json:"is_auto_commit,omitempty"`
	// (Optional) Specifies the preferred namespace within the database or system that the client wants to interact with.
	NamespaceName *string `protobuf:"bytes,2,opt,name=namespace_name,json=namespaceName,proto3,oneof" json:"namespace_name,omitempty"`
}

func (x *ConnectionProperties) Reset() {
	*x = ConnectionProperties{}
	if protoimpl.UnsafeEnabled {
		mi := &file_org_polypheny_prism_connection_requests_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConnectionProperties) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConnectionProperties) ProtoMessage() {}

func (x *ConnectionProperties) ProtoReflect() protoreflect.Message {
	mi := &file_org_polypheny_prism_connection_requests_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConnectionProperties.ProtoReflect.Descriptor instead.
func (*ConnectionProperties) Descriptor() ([]byte, []int) {
	return file_org_polypheny_prism_connection_requests_proto_rawDescGZIP(), []int{2}
}

func (x *ConnectionProperties) GetIsAutoCommit() bool {
	if x != nil && x.IsAutoCommit != nil {
		return *x.IsAutoCommit
	}
	return false
}

func (x *ConnectionProperties) GetNamespaceName() string {
	if x != nil && x.NamespaceName != nil {
		return *x.NamespaceName
	}
	return ""
}

//
//The ConnectionCheckRequest message in combination with the corresponding remote procedure call is utilized to verify the current state of an established connection.
//It acts as a simple “ping” request, enabling clients to ascertain if the server is responsive and if the connection is still valid.
//This message does not contain any fields. It simply acts as an indicator to prompt the server for a ConnectionCheckResponse.
type DisconnectRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DisconnectRequest) Reset() {
	*x = DisconnectRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_org_polypheny_prism_connection_requests_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DisconnectRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DisconnectRequest) ProtoMessage() {}

func (x *DisconnectRequest) ProtoReflect() protoreflect.Message {
	mi := &file_org_polypheny_prism_connection_requests_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DisconnectRequest.ProtoReflect.Descriptor instead.
func (*DisconnectRequest) Descriptor() ([]byte, []int) {
	return file_org_polypheny_prism_connection_requests_proto_rawDescGZIP(), []int{3}
}

//
//The ConnectionCheckRequest message in combination with the corresponding remote procedure call is utilized to verify the current state of an established connection.
//It acts as a simple “ping” request, enabling clients to ascertain if the server is responsive and if the connection is still valid.
type ConnectionCheckRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ConnectionCheckRequest) Reset() {
	*x = ConnectionCheckRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_org_polypheny_prism_connection_requests_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConnectionCheckRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConnectionCheckRequest) ProtoMessage() {}

func (x *ConnectionCheckRequest) ProtoReflect() protoreflect.Message {
	mi := &file_org_polypheny_prism_connection_requests_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConnectionCheckRequest.ProtoReflect.Descriptor instead.
func (*ConnectionCheckRequest) Descriptor() ([]byte, []int) {
	return file_org_polypheny_prism_connection_requests_proto_rawDescGZIP(), []int{4}
}

var File_org_polypheny_prism_connection_requests_proto protoreflect.FileDescriptor

var file_org_polypheny_prism_connection_requests_proto_rawDesc = []byte{
	0x0a, 0x2d, 0x6f, 0x72, 0x67, 0x2f, 0x70, 0x6f, 0x6c, 0x79, 0x70, 0x68, 0x65, 0x6e, 0x79, 0x2f,
	0x70, 0x72, 0x69, 0x73, 0x6d, 0x2f, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x5f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x13, 0x6f, 0x72, 0x67, 0x2e, 0x70, 0x6f, 0x6c, 0x79, 0x70, 0x68, 0x65, 0x6e, 0x79, 0x2e, 0x70,
	0x72, 0x69, 0x73, 0x6d, 0x22, 0xc6, 0x02, 0x0a, 0x11, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2a, 0x0a, 0x11, 0x6d, 0x61,
	0x6a, 0x6f, 0x72, 0x5f, 0x61, 0x70, 0x69, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0f, 0x6d, 0x61, 0x6a, 0x6f, 0x72, 0x41, 0x70, 0x69, 0x56,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x2a, 0x0a, 0x11, 0x6d, 0x69, 0x6e, 0x6f, 0x72, 0x5f,
	0x61, 0x70, 0x69, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x0f, 0x6d, 0x69, 0x6e, 0x6f, 0x72, 0x41, 0x70, 0x69, 0x56, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x12, 0x1f, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65,
	0x88, 0x01, 0x01, 0x12, 0x1f, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x09, 0x48, 0x01, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72,
	0x64, 0x88, 0x01, 0x01, 0x12, 0x63, 0x0a, 0x15, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x5f, 0x70, 0x72, 0x6f, 0x70, 0x65, 0x72, 0x74, 0x69, 0x65, 0x73, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x6f, 0x72, 0x67, 0x2e, 0x70, 0x6f, 0x6c, 0x79, 0x70, 0x68,
	0x65, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x69, 0x73, 0x6d, 0x2e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x50, 0x72, 0x6f, 0x70, 0x65, 0x72, 0x74, 0x69, 0x65, 0x73, 0x48, 0x02,
	0x52, 0x14, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x50, 0x72, 0x6f, 0x70,
	0x65, 0x72, 0x74, 0x69, 0x65, 0x73, 0x88, 0x01, 0x01, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x75, 0x73,
	0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x70, 0x61, 0x73, 0x73, 0x77,
	0x6f, 0x72, 0x64, 0x42, 0x18, 0x0a, 0x16, 0x5f, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x5f, 0x70, 0x72, 0x6f, 0x70, 0x65, 0x72, 0x74, 0x69, 0x65, 0x73, 0x22, 0x83, 0x01,
	0x0a, 0x21, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x50, 0x72, 0x6f, 0x70,
	0x65, 0x72, 0x74, 0x69, 0x65, 0x73, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x5e, 0x0a, 0x15, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x5f, 0x70, 0x72, 0x6f, 0x70, 0x65, 0x72, 0x74, 0x69, 0x65, 0x73, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x29, 0x2e, 0x6f, 0x72, 0x67, 0x2e, 0x70, 0x6f, 0x6c, 0x79, 0x70, 0x68, 0x65,
	0x6e, 0x79, 0x2e, 0x70, 0x72, 0x69, 0x73, 0x6d, 0x2e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x50, 0x72, 0x6f, 0x70, 0x65, 0x72, 0x74, 0x69, 0x65, 0x73, 0x52, 0x14, 0x63,
	0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x50, 0x72, 0x6f, 0x70, 0x65, 0x72, 0x74,
	0x69, 0x65, 0x73, 0x22, 0x93, 0x01, 0x0a, 0x14, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x50, 0x72, 0x6f, 0x70, 0x65, 0x72, 0x74, 0x69, 0x65, 0x73, 0x12, 0x29, 0x0a, 0x0e,
	0x69, 0x73, 0x5f, 0x61, 0x75, 0x74, 0x6f, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x08, 0x48, 0x00, 0x52, 0x0c, 0x69, 0x73, 0x41, 0x75, 0x74, 0x6f, 0x43, 0x6f,
	0x6d, 0x6d, 0x69, 0x74, 0x88, 0x01, 0x01, 0x12, 0x2a, 0x0a, 0x0e, 0x6e, 0x61, 0x6d, 0x65, 0x73,
	0x70, 0x61, 0x63, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48,
	0x01, 0x52, 0x0d, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65,
	0x88, 0x01, 0x01, 0x42, 0x11, 0x0a, 0x0f, 0x5f, 0x69, 0x73, 0x5f, 0x61, 0x75, 0x74, 0x6f, 0x5f,
	0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x42, 0x11, 0x0a, 0x0f, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x73,
	0x70, 0x61, 0x63, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x13, 0x0a, 0x11, 0x44, 0x69, 0x73,
	0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x18,
	0x0a, 0x16, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x68, 0x65, 0x63,
	0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x42, 0x58, 0x0a, 0x13, 0x6f, 0x72, 0x67, 0x2e,
	0x70, 0x6f, 0x6c, 0x79, 0x70, 0x68, 0x65, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x69, 0x73, 0x6d, 0x42,
	0x12, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x73, 0x50, 0x01, 0x5a, 0x19, 0x6f, 0x72, 0x67, 0x2f, 0x70, 0x6f, 0x6c, 0x79, 0x70,
	0x68, 0x65, 0x6e, 0x79, 0x2f, 0x70, 0x72, 0x69, 0x73, 0x6d, 0x3b, 0x70, 0x72, 0x69, 0x73, 0x6d,
	0xaa, 0x02, 0x0f, 0x50, 0x6f, 0x6c, 0x79, 0x70, 0x68, 0x65, 0x6e, 0x79, 0x2e, 0x50, 0x72, 0x69,
	0x73, 0x6d, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_org_polypheny_prism_connection_requests_proto_rawDescOnce sync.Once
	file_org_polypheny_prism_connection_requests_proto_rawDescData = file_org_polypheny_prism_connection_requests_proto_rawDesc
)

func file_org_polypheny_prism_connection_requests_proto_rawDescGZIP() []byte {
	file_org_polypheny_prism_connection_requests_proto_rawDescOnce.Do(func() {
		file_org_polypheny_prism_connection_requests_proto_rawDescData = protoimpl.X.CompressGZIP(file_org_polypheny_prism_connection_requests_proto_rawDescData)
	})
	return file_org_polypheny_prism_connection_requests_proto_rawDescData
}

var file_org_polypheny_prism_connection_requests_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_org_polypheny_prism_connection_requests_proto_goTypes = []interface{}{
	(*ConnectionRequest)(nil),                 // 0: org.polypheny.prism.ConnectionRequest
	(*ConnectionPropertiesUpdateRequest)(nil), // 1: org.polypheny.prism.ConnectionPropertiesUpdateRequest
	(*ConnectionProperties)(nil),              // 2: org.polypheny.prism.ConnectionProperties
	(*DisconnectRequest)(nil),                 // 3: org.polypheny.prism.DisconnectRequest
	(*ConnectionCheckRequest)(nil),            // 4: org.polypheny.prism.ConnectionCheckRequest
}
var file_org_polypheny_prism_connection_requests_proto_depIdxs = []int32{
	2, // 0: org.polypheny.prism.ConnectionRequest.connection_properties:type_name -> org.polypheny.prism.ConnectionProperties
	2, // 1: org.polypheny.prism.ConnectionPropertiesUpdateRequest.connection_properties:type_name -> org.polypheny.prism.ConnectionProperties
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_org_polypheny_prism_connection_requests_proto_init() }
func file_org_polypheny_prism_connection_requests_proto_init() {
	if File_org_polypheny_prism_connection_requests_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_org_polypheny_prism_connection_requests_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConnectionRequest); i {
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
		file_org_polypheny_prism_connection_requests_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConnectionPropertiesUpdateRequest); i {
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
		file_org_polypheny_prism_connection_requests_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConnectionProperties); i {
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
		file_org_polypheny_prism_connection_requests_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DisconnectRequest); i {
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
		file_org_polypheny_prism_connection_requests_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConnectionCheckRequest); i {
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
	file_org_polypheny_prism_connection_requests_proto_msgTypes[0].OneofWrappers = []interface{}{}
	file_org_polypheny_prism_connection_requests_proto_msgTypes[2].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_org_polypheny_prism_connection_requests_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_org_polypheny_prism_connection_requests_proto_goTypes,
		DependencyIndexes: file_org_polypheny_prism_connection_requests_proto_depIdxs,
		MessageInfos:      file_org_polypheny_prism_connection_requests_proto_msgTypes,
	}.Build()
	File_org_polypheny_prism_connection_requests_proto = out.File
	file_org_polypheny_prism_connection_requests_proto_rawDesc = nil
	file_org_polypheny_prism_connection_requests_proto_goTypes = nil
	file_org_polypheny_prism_connection_requests_proto_depIdxs = nil
}
