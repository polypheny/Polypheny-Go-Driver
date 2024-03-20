//
//Messages that are used for querying general metadata.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.24.3
// source: polyprism/meta_requests.proto

package protos

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

// The DbmsVersionRequest message is utilized to request the version information of the database management system (DBMS) in use.
// It acts as a trigger for the server to respond with the specific version details of the DBMS.
// This message does not contain any fields. It simply acts as an indicator to prompt the server for DbmsVersionResponse.
type DbmsVersionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DbmsVersionRequest) Reset() {
	*x = DbmsVersionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_polyprism_meta_requests_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DbmsVersionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DbmsVersionRequest) ProtoMessage() {}

func (x *DbmsVersionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_polyprism_meta_requests_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DbmsVersionRequest.ProtoReflect.Descriptor instead.
func (*DbmsVersionRequest) Descriptor() ([]byte, []int) {
	return file_polyprism_meta_requests_proto_rawDescGZIP(), []int{0}
}

// The LanguageRequest message facilitates the retrieval of supported query languages that can be used for constructing statements.
// This message does not contain any fields. It acts as a request for the server to return the list of supported query languages in a LanguageResponse.
type LanguageRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *LanguageRequest) Reset() {
	*x = LanguageRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_polyprism_meta_requests_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LanguageRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LanguageRequest) ProtoMessage() {}

func (x *LanguageRequest) ProtoReflect() protoreflect.Message {
	mi := &file_polyprism_meta_requests_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LanguageRequest.ProtoReflect.Descriptor instead.
func (*LanguageRequest) Descriptor() ([]byte, []int) {
	return file_polyprism_meta_requests_proto_rawDescGZIP(), []int{1}
}

// The DatabasesRequest message is employed to solicit a list of databases that are currently managed by the DBMS.
// This is pivotal for clients aiming to interact with a specific database or to understand the landscape of databases under management.
// This message does not have any fields. It signals the server to provide a list of databases in a DatabasesResponse.
type DatabasesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DatabasesRequest) Reset() {
	*x = DatabasesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_polyprism_meta_requests_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DatabasesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DatabasesRequest) ProtoMessage() {}

func (x *DatabasesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_polyprism_meta_requests_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DatabasesRequest.ProtoReflect.Descriptor instead.
func (*DatabasesRequest) Descriptor() ([]byte, []int) {
	return file_polyprism_meta_requests_proto_rawDescGZIP(), []int{2}
}

// The TableTypesRequest message is designed to request information about the different types of tables that are supported or recognized by the DBMS.
// This message does not contain any fields. It’s a prompt for the server to respond with details about the supported table types in a TableTypesResponse.
type TableTypesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *TableTypesRequest) Reset() {
	*x = TableTypesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_polyprism_meta_requests_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TableTypesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TableTypesRequest) ProtoMessage() {}

func (x *TableTypesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_polyprism_meta_requests_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TableTypesRequest.ProtoReflect.Descriptor instead.
func (*TableTypesRequest) Descriptor() ([]byte, []int) {
	return file_polyprism_meta_requests_proto_rawDescGZIP(), []int{3}
}

// The TypesRequest message is deployed to obtain a list of data types supported by the database management system.
// This helps clients understand the range of data types they can utilize when defining or querying tables.
// This message does not contain any fields. It simply prompts the server to return a list of supported data types in a TypesResponse.
type TypesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *TypesRequest) Reset() {
	*x = TypesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_polyprism_meta_requests_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TypesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TypesRequest) ProtoMessage() {}

func (x *TypesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_polyprism_meta_requests_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TypesRequest.ProtoReflect.Descriptor instead.
func (*TypesRequest) Descriptor() ([]byte, []int) {
	return file_polyprism_meta_requests_proto_rawDescGZIP(), []int{4}
}

// The UserDefinedTypesRequest message facilitates the retrieval of user-defined data types present in the database.
// These custom data types can be crafted to meet specific needs not met by standard types.
// This message is devoid of fields. It acts as a signal for the server to provide details of user-defined data types in a UserDefinedTypesResponse.
type UserDefinedTypesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *UserDefinedTypesRequest) Reset() {
	*x = UserDefinedTypesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_polyprism_meta_requests_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserDefinedTypesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserDefinedTypesRequest) ProtoMessage() {}

func (x *UserDefinedTypesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_polyprism_meta_requests_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserDefinedTypesRequest.ProtoReflect.Descriptor instead.
func (*UserDefinedTypesRequest) Descriptor() ([]byte, []int) {
	return file_polyprism_meta_requests_proto_rawDescGZIP(), []int{5}
}

// The SqlStringFunctionsRequest message is used to solicit information about the string functions supported by the SQL implementation of the DBMS.
// This message is field less and prompts the server to return details of the available string functions in an SqlStringFunctionsResponse.
type SqlStringFunctionsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SqlStringFunctionsRequest) Reset() {
	*x = SqlStringFunctionsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_polyprism_meta_requests_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SqlStringFunctionsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SqlStringFunctionsRequest) ProtoMessage() {}

func (x *SqlStringFunctionsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_polyprism_meta_requests_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SqlStringFunctionsRequest.ProtoReflect.Descriptor instead.
func (*SqlStringFunctionsRequest) Descriptor() ([]byte, []int) {
	return file_polyprism_meta_requests_proto_rawDescGZIP(), []int{6}
}

// The SqlSystemFunctionsRequest message aims to retrieve the list of system functions provided by the SQL implementation of the DBMS.
// Without any fields, this message indicates the server to respond with details about the system functions in an SqlSystemFunctionsResponse.
type SqlSystemFunctionsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SqlSystemFunctionsRequest) Reset() {
	*x = SqlSystemFunctionsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_polyprism_meta_requests_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SqlSystemFunctionsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SqlSystemFunctionsRequest) ProtoMessage() {}

func (x *SqlSystemFunctionsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_polyprism_meta_requests_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SqlSystemFunctionsRequest.ProtoReflect.Descriptor instead.
func (*SqlSystemFunctionsRequest) Descriptor() ([]byte, []int) {
	return file_polyprism_meta_requests_proto_rawDescGZIP(), []int{7}
}

// The SqlTimeDateFunctionsRequest message is dispatched to fetch a list of time and date functions supported by the SQL implementation of the DBMS.
// This message, being field less, acts as a request for the server to list time and date functions in a SqlTimeDateFunctionsResponse.
type SqlTimeDateFunctionsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SqlTimeDateFunctionsRequest) Reset() {
	*x = SqlTimeDateFunctionsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_polyprism_meta_requests_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SqlTimeDateFunctionsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SqlTimeDateFunctionsRequest) ProtoMessage() {}

func (x *SqlTimeDateFunctionsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_polyprism_meta_requests_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SqlTimeDateFunctionsRequest.ProtoReflect.Descriptor instead.
func (*SqlTimeDateFunctionsRequest) Descriptor() ([]byte, []int) {
	return file_polyprism_meta_requests_proto_rawDescGZIP(), []int{8}
}

// The SqlNumericFunctionsRequest message endeavors to obtain details about the numeric functions provided by the SQL implementation of the DBMS.
// This message, devoid of fields, prompts the server to respond with information about numeric functions in an SqlNumericFunctionsResponse.
type SqlNumericFunctionsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SqlNumericFunctionsRequest) Reset() {
	*x = SqlNumericFunctionsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_polyprism_meta_requests_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SqlNumericFunctionsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SqlNumericFunctionsRequest) ProtoMessage() {}

func (x *SqlNumericFunctionsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_polyprism_meta_requests_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SqlNumericFunctionsRequest.ProtoReflect.Descriptor instead.
func (*SqlNumericFunctionsRequest) Descriptor() ([]byte, []int) {
	return file_polyprism_meta_requests_proto_rawDescGZIP(), []int{9}
}

// The SqlKeywordsRequest message is designed to request a list of reserved keywords utilized by the SQL implementation of the DBMS.
// With no fields, this message acts as an indicator for the server to provide the list of SQL keywords in an SqlKeywordsResponse.
type SqlKeywordsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SqlKeywordsRequest) Reset() {
	*x = SqlKeywordsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_polyprism_meta_requests_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SqlKeywordsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SqlKeywordsRequest) ProtoMessage() {}

func (x *SqlKeywordsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_polyprism_meta_requests_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SqlKeywordsRequest.ProtoReflect.Descriptor instead.
func (*SqlKeywordsRequest) Descriptor() ([]byte, []int) {
	return file_polyprism_meta_requests_proto_rawDescGZIP(), []int{10}
}

// The ProceduresRequest message is employed to retrieve a list of stored procedures in the database for a specified query language.
// The client can also narrow down the results by specifying a procedure name pattern.
type ProceduresRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The query language for which the procedures are defined. This determines the set of procedures that will be retrieved.
	Language string `protobuf:"bytes,1,opt,name=language,proto3" json:"language,omitempty"`
	// A pattern to filter the names of procedures. For example, providing “get%” might retrieve procedures named getUser, getDetails, etc.
	// If not specified, all procedures for the provided language will be returned. Like in sql, the symbol _ can be used to match a single character.
	ProcedureNamePattern *string `protobuf:"bytes,3,opt,name=procedure_name_pattern,json=procedureNamePattern,proto3,oneof" json:"procedure_name_pattern,omitempty"`
}

func (x *ProceduresRequest) Reset() {
	*x = ProceduresRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_polyprism_meta_requests_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProceduresRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProceduresRequest) ProtoMessage() {}

func (x *ProceduresRequest) ProtoReflect() protoreflect.Message {
	mi := &file_polyprism_meta_requests_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProceduresRequest.ProtoReflect.Descriptor instead.
func (*ProceduresRequest) Descriptor() ([]byte, []int) {
	return file_polyprism_meta_requests_proto_rawDescGZIP(), []int{11}
}

func (x *ProceduresRequest) GetLanguage() string {
	if x != nil {
		return x.Language
	}
	return ""
}

func (x *ProceduresRequest) GetProcedureNamePattern() string {
	if x != nil && x.ProcedureNamePattern != nil {
		return *x.ProcedureNamePattern
	}
	return ""
}

// The ClientInfoPropertiesRequest message facilitates the acquisition of client information properties stored in the database.
// These properties can offer additional context about the connected client.
// This message doesn’t possess any fields.
// It acts as a directive for the server to provide the associated client information properties in a ClientInfoPropertiesResponse.
type ClientInfoPropertiesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ClientInfoPropertiesRequest) Reset() {
	*x = ClientInfoPropertiesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_polyprism_meta_requests_proto_msgTypes[12]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClientInfoPropertiesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClientInfoPropertiesRequest) ProtoMessage() {}

func (x *ClientInfoPropertiesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_polyprism_meta_requests_proto_msgTypes[12]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClientInfoPropertiesRequest.ProtoReflect.Descriptor instead.
func (*ClientInfoPropertiesRequest) Descriptor() ([]byte, []int) {
	return file_polyprism_meta_requests_proto_rawDescGZIP(), []int{12}
}

// The ClientInfoPropertyMetaRequest message aids in extracting metadata about the client information properties present in the database.
// This fieldless message prompts the server to detail metadata concerning client information properties in a ClientInfoPropertyMetaResponse.
type ClientInfoPropertyMetaRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ClientInfoPropertyMetaRequest) Reset() {
	*x = ClientInfoPropertyMetaRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_polyprism_meta_requests_proto_msgTypes[13]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClientInfoPropertyMetaRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClientInfoPropertyMetaRequest) ProtoMessage() {}

func (x *ClientInfoPropertyMetaRequest) ProtoReflect() protoreflect.Message {
	mi := &file_polyprism_meta_requests_proto_msgTypes[13]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClientInfoPropertyMetaRequest.ProtoReflect.Descriptor instead.
func (*ClientInfoPropertyMetaRequest) Descriptor() ([]byte, []int) {
	return file_polyprism_meta_requests_proto_rawDescGZIP(), []int{13}
}

// The FunctionsRequest message is wielded to obtain a list of functions from the database based on the specified query language and function category.
type FunctionsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Defines the query language for which the functions are sought. This field determines the range of functions that will be retrieved.
	QueryLanguage string `protobuf:"bytes,1,opt,name=query_language,json=queryLanguage,proto3" json:"query_language,omitempty"`
	// Categorizes the function, allowing clients to filter results based on specific categories, such as “NUMERIC”, “STRING”, or “SYSTEM”.
	// This helps in refining the search for specific types of functions.
	FunctionCategory string `protobuf:"bytes,2,opt,name=function_category,json=functionCategory,proto3" json:"function_category,omitempty"`
}

func (x *FunctionsRequest) Reset() {
	*x = FunctionsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_polyprism_meta_requests_proto_msgTypes[14]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FunctionsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FunctionsRequest) ProtoMessage() {}

func (x *FunctionsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_polyprism_meta_requests_proto_msgTypes[14]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FunctionsRequest.ProtoReflect.Descriptor instead.
func (*FunctionsRequest) Descriptor() ([]byte, []int) {
	return file_polyprism_meta_requests_proto_rawDescGZIP(), []int{14}
}

func (x *FunctionsRequest) GetQueryLanguage() string {
	if x != nil {
		return x.QueryLanguage
	}
	return ""
}

func (x *FunctionsRequest) GetFunctionCategory() string {
	if x != nil {
		return x.FunctionCategory
	}
	return ""
}

var File_polyprism_meta_requests_proto protoreflect.FileDescriptor

var file_polyprism_meta_requests_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x70, 0x6f, 0x6c, 0x79, 0x70, 0x72, 0x69, 0x73, 0x6d, 0x2f, 0x6d, 0x65, 0x74, 0x61,
	0x5f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x18, 0x70, 0x6f, 0x6c, 0x79, 0x70, 0x68, 0x65, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x69, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x22, 0x14, 0x0a, 0x12, 0x44, 0x62, 0x6d,
	0x73, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22,
	0x11, 0x0a, 0x0f, 0x4c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x22, 0x12, 0x0a, 0x10, 0x44, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x13, 0x0a, 0x11, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x54,
	0x79, 0x70, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x0e, 0x0a, 0x0c, 0x54,
	0x79, 0x70, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x19, 0x0a, 0x17, 0x55,
	0x73, 0x65, 0x72, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x65, 0x64, 0x54, 0x79, 0x70, 0x65, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x1b, 0x0a, 0x19, 0x53, 0x71, 0x6c, 0x53, 0x74, 0x72,
	0x69, 0x6e, 0x67, 0x46, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x22, 0x1b, 0x0a, 0x19, 0x53, 0x71, 0x6c, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d,
	0x46, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x22, 0x1d, 0x0a, 0x1b, 0x53, 0x71, 0x6c, 0x54, 0x69, 0x6d, 0x65, 0x44, 0x61, 0x74, 0x65, 0x46,
	0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22,
	0x1c, 0x0a, 0x1a, 0x53, 0x71, 0x6c, 0x4e, 0x75, 0x6d, 0x65, 0x72, 0x69, 0x63, 0x46, 0x75, 0x6e,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x14, 0x0a,
	0x12, 0x53, 0x71, 0x6c, 0x4b, 0x65, 0x79, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x22, 0x85, 0x01, 0x0a, 0x11, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x64, 0x75, 0x72,
	0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x61, 0x6e,
	0x67, 0x75, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6c, 0x61, 0x6e,
	0x67, 0x75, 0x61, 0x67, 0x65, 0x12, 0x39, 0x0a, 0x16, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x64, 0x75,
	0x72, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x5f, 0x70, 0x61, 0x74, 0x74, 0x65, 0x72, 0x6e, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x14, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x64, 0x75,
	0x72, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x50, 0x61, 0x74, 0x74, 0x65, 0x72, 0x6e, 0x88, 0x01, 0x01,
	0x42, 0x19, 0x0a, 0x17, 0x5f, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x64, 0x75, 0x72, 0x65, 0x5f, 0x6e,
	0x61, 0x6d, 0x65, 0x5f, 0x70, 0x61, 0x74, 0x74, 0x65, 0x72, 0x6e, 0x22, 0x1d, 0x0a, 0x1b, 0x43,
	0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x50, 0x72, 0x6f, 0x70, 0x65, 0x72, 0x74,
	0x69, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x1f, 0x0a, 0x1d, 0x43, 0x6c,
	0x69, 0x65, 0x6e, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x50, 0x72, 0x6f, 0x70, 0x65, 0x72, 0x74, 0x79,
	0x4d, 0x65, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x66, 0x0a, 0x10, 0x46,
	0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x25, 0x0a, 0x0e, 0x71, 0x75, 0x65, 0x72, 0x79, 0x5f, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x71, 0x75, 0x65, 0x72, 0x79, 0x4c, 0x61,
	0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x12, 0x2b, 0x0a, 0x11, 0x66, 0x75, 0x6e, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x5f, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x10, 0x66, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x61, 0x74, 0x65, 0x67,
	0x6f, 0x72, 0x79, 0x42, 0x37, 0x0a, 0x25, 0x6f, 0x72, 0x67, 0x2e, 0x70, 0x6f, 0x6c, 0x79, 0x70,
	0x68, 0x65, 0x6e, 0x79, 0x2e, 0x64, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x69, 0x6e, 0x74,
	0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x42, 0x0c, 0x4d, 0x65,
	0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x73, 0x50, 0x01, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_polyprism_meta_requests_proto_rawDescOnce sync.Once
	file_polyprism_meta_requests_proto_rawDescData = file_polyprism_meta_requests_proto_rawDesc
)

func file_polyprism_meta_requests_proto_rawDescGZIP() []byte {
	file_polyprism_meta_requests_proto_rawDescOnce.Do(func() {
		file_polyprism_meta_requests_proto_rawDescData = protoimpl.X.CompressGZIP(file_polyprism_meta_requests_proto_rawDescData)
	})
	return file_polyprism_meta_requests_proto_rawDescData
}

var file_polyprism_meta_requests_proto_msgTypes = make([]protoimpl.MessageInfo, 15)
var file_polyprism_meta_requests_proto_goTypes = []interface{}{
	(*DbmsVersionRequest)(nil),            // 0: polypheny.protointerface.DbmsVersionRequest
	(*LanguageRequest)(nil),               // 1: polypheny.protointerface.LanguageRequest
	(*DatabasesRequest)(nil),              // 2: polypheny.protointerface.DatabasesRequest
	(*TableTypesRequest)(nil),             // 3: polypheny.protointerface.TableTypesRequest
	(*TypesRequest)(nil),                  // 4: polypheny.protointerface.TypesRequest
	(*UserDefinedTypesRequest)(nil),       // 5: polypheny.protointerface.UserDefinedTypesRequest
	(*SqlStringFunctionsRequest)(nil),     // 6: polypheny.protointerface.SqlStringFunctionsRequest
	(*SqlSystemFunctionsRequest)(nil),     // 7: polypheny.protointerface.SqlSystemFunctionsRequest
	(*SqlTimeDateFunctionsRequest)(nil),   // 8: polypheny.protointerface.SqlTimeDateFunctionsRequest
	(*SqlNumericFunctionsRequest)(nil),    // 9: polypheny.protointerface.SqlNumericFunctionsRequest
	(*SqlKeywordsRequest)(nil),            // 10: polypheny.protointerface.SqlKeywordsRequest
	(*ProceduresRequest)(nil),             // 11: polypheny.protointerface.ProceduresRequest
	(*ClientInfoPropertiesRequest)(nil),   // 12: polypheny.protointerface.ClientInfoPropertiesRequest
	(*ClientInfoPropertyMetaRequest)(nil), // 13: polypheny.protointerface.ClientInfoPropertyMetaRequest
	(*FunctionsRequest)(nil),              // 14: polypheny.protointerface.FunctionsRequest
}
var file_polyprism_meta_requests_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_polyprism_meta_requests_proto_init() }
func file_polyprism_meta_requests_proto_init() {
	if File_polyprism_meta_requests_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_polyprism_meta_requests_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DbmsVersionRequest); i {
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
		file_polyprism_meta_requests_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LanguageRequest); i {
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
		file_polyprism_meta_requests_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DatabasesRequest); i {
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
		file_polyprism_meta_requests_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TableTypesRequest); i {
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
		file_polyprism_meta_requests_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TypesRequest); i {
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
		file_polyprism_meta_requests_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserDefinedTypesRequest); i {
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
		file_polyprism_meta_requests_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SqlStringFunctionsRequest); i {
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
		file_polyprism_meta_requests_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SqlSystemFunctionsRequest); i {
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
		file_polyprism_meta_requests_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SqlTimeDateFunctionsRequest); i {
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
		file_polyprism_meta_requests_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SqlNumericFunctionsRequest); i {
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
		file_polyprism_meta_requests_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SqlKeywordsRequest); i {
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
		file_polyprism_meta_requests_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProceduresRequest); i {
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
		file_polyprism_meta_requests_proto_msgTypes[12].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClientInfoPropertiesRequest); i {
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
		file_polyprism_meta_requests_proto_msgTypes[13].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClientInfoPropertyMetaRequest); i {
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
		file_polyprism_meta_requests_proto_msgTypes[14].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FunctionsRequest); i {
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
	file_polyprism_meta_requests_proto_msgTypes[11].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_polyprism_meta_requests_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   15,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_polyprism_meta_requests_proto_goTypes,
		DependencyIndexes: file_polyprism_meta_requests_proto_depIdxs,
		MessageInfos:      file_polyprism_meta_requests_proto_msgTypes,
	}.Build()
	File_polyprism_meta_requests_proto = out.File
	file_polyprism_meta_requests_proto_rawDesc = nil
	file_polyprism_meta_requests_proto_goTypes = nil
	file_polyprism_meta_requests_proto_depIdxs = nil
}
