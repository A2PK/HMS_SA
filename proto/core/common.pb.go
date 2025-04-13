// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: proto/core/common.proto

package core

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	structpb "google.golang.org/protobuf/types/known/structpb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Represents common filtering, pagination, and sorting options.
// Based on pkg/core/types/common.go FilterOptions struct.
type FilterOptions struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Maximum number of items to return per page.
	Limit *int32 `protobuf:"varint,1,opt,name=limit,proto3,oneof" json:"limit,omitempty"`
	// Number of items to skip before starting to collect the result set.
	Offset *int32 `protobuf:"varint,2,opt,name=offset,proto3,oneof" json:"offset,omitempty"`
	// Field name to sort the results by.
	SortBy *string `protobuf:"bytes,3,opt,name=sort_by,json=sortBy,proto3,oneof" json:"sort_by,omitempty"`
	// Whether to sort in descending order. Defaults to false (ascending).
	SortDesc *bool `protobuf:"varint,4,opt,name=sort_desc,json=sortDesc,proto3,oneof" json:"sort_desc,omitempty"`
	// Key-value pairs for specific field filtering.
	// Uses google.protobuf.Value to allow various types (string, number, bool, null).
	Filters map[string]*structpb.Value `protobuf:"bytes,5,rep,name=filters,proto3" json:"filters,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	// Whether to include soft-deleted records in the results.
	IncludeDeleted *bool `protobuf:"varint,8,opt,name=include_deleted,json=includeDeleted,proto3,oneof" json:"include_deleted,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *FilterOptions) Reset() {
	*x = FilterOptions{}
	mi := &file_proto_core_common_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FilterOptions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FilterOptions) ProtoMessage() {}

func (x *FilterOptions) ProtoReflect() protoreflect.Message {
	mi := &file_proto_core_common_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FilterOptions.ProtoReflect.Descriptor instead.
func (*FilterOptions) Descriptor() ([]byte, []int) {
	return file_proto_core_common_proto_rawDescGZIP(), []int{0}
}

func (x *FilterOptions) GetLimit() int32 {
	if x != nil && x.Limit != nil {
		return *x.Limit
	}
	return 0
}

func (x *FilterOptions) GetOffset() int32 {
	if x != nil && x.Offset != nil {
		return *x.Offset
	}
	return 0
}

func (x *FilterOptions) GetSortBy() string {
	if x != nil && x.SortBy != nil {
		return *x.SortBy
	}
	return ""
}

func (x *FilterOptions) GetSortDesc() bool {
	if x != nil && x.SortDesc != nil {
		return *x.SortDesc
	}
	return false
}

func (x *FilterOptions) GetFilters() map[string]*structpb.Value {
	if x != nil {
		return x.Filters
	}
	return nil
}

func (x *FilterOptions) GetIncludeDeleted() bool {
	if x != nil && x.IncludeDeleted != nil {
		return *x.IncludeDeleted
	}
	return false
}

// Represents common pagination metadata included in list responses.
// Based on pkg/core/types/common.go PaginationResult struct (metadata fields only).
// Specific list responses should include this alongside their repeated items field.
type PaginationInfo struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Total number of items matching the query criteria across all pages.
	TotalItems int64 `protobuf:"varint,1,opt,name=total_items,json=totalItems,proto3" json:"total_items,omitempty"`
	// The limit (page size) used for the current response.
	Limit int32 `protobuf:"varint,2,opt,name=limit,proto3" json:"limit,omitempty"`
	// The offset (number of items skipped) used for the current response.
	Offset        int32 `protobuf:"varint,3,opt,name=offset,proto3" json:"offset,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PaginationInfo) Reset() {
	*x = PaginationInfo{}
	mi := &file_proto_core_common_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PaginationInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PaginationInfo) ProtoMessage() {}

func (x *PaginationInfo) ProtoReflect() protoreflect.Message {
	mi := &file_proto_core_common_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PaginationInfo.ProtoReflect.Descriptor instead.
func (*PaginationInfo) Descriptor() ([]byte, []int) {
	return file_proto_core_common_proto_rawDescGZIP(), []int{1}
}

func (x *PaginationInfo) GetTotalItems() int64 {
	if x != nil {
		return x.TotalItems
	}
	return 0
}

func (x *PaginationInfo) GetLimit() int32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *PaginationInfo) GetOffset() int32 {
	if x != nil {
		return x.Offset
	}
	return 0
}

var File_proto_core_common_proto protoreflect.FileDescriptor

const file_proto_core_common_proto_rawDesc = "" +
	"\n" +
	"\x17proto/core/common.proto\x12\x04core\x1a\x1cgoogle/protobuf/struct.proto\x1a.protoc-gen-openapiv2/options/annotations.proto\"\xc5\a\n" +
	"\rFilterOptions\x12S\n" +
	"\x05limit\x18\x01 \x01(\x05B8\x92A52+Maximum number of items to return per page.:\x0250J\x0250H\x00R\x05limit\x88\x01\x01\x12{\n" +
	"\x06offset\x18\x02 \x01(\x05B^\x92A[2SNumber of items to skip before starting to collect the result set (for pagination).:\x010J\x010H\x01R\x06offset\x88\x01\x01\x12~\n" +
	"\asort_by\x18\x03 \x01(\tB`\x92A]2?Field name to sort the results by (e.g., 'created_at', 'name').:\f\"created_at\"J\f\"created_at\"H\x02R\x06sortBy\x88\x01\x01\x12[\n" +
	"\tsort_desc\x18\x04 \x01(\bB9\x92A62(Set to true to sort in descending order.:\x04trueJ\x04trueH\x03R\bsortDesc\x88\x01\x01\x12\xef\x01\n" +
	"\afilters\x18\x05 \x03(\v2 .core.FilterOptions.FiltersEntryB\xb2\x01\x92A\xae\x012\x8e\x01Key-value pairs for specific field filtering. Values should correspond to google.protobuf.Value structure (e.g., {\"email\": \"user@gmail.com\"}).J\x1b{\"email\": \"user@gmail.com\"}R\afilters\x12|\n" +
	"\x0finclude_deleted\x18\b \x01(\bBN\x92AK2;Set to true to include soft-deleted records in the results.:\x05falseJ\x05falseH\x04R\x0eincludeDeleted\x88\x01\x01\x1aR\n" +
	"\fFiltersEntry\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x12,\n" +
	"\x05value\x18\x02 \x01(\v2\x16.google.protobuf.ValueR\x05value:\x028\x01B\b\n" +
	"\x06_limitB\t\n" +
	"\a_offsetB\n" +
	"\n" +
	"\b_sort_byB\f\n" +
	"\n" +
	"_sort_descB\x12\n" +
	"\x10_include_deleted\"\xbb\x02\n" +
	"\x0ePaginationInfo\x12o\n" +
	"\vtotal_items\x18\x01 \x01(\x03BN\x92AK2CTotal number of items matching the query criteria across all pages.J\x041234R\n" +
	"totalItems\x12S\n" +
	"\x05limit\x18\x02 \x01(\x05B=\x92A:24The limit (page size) used for the current response.J\x0250R\x05limit\x12c\n" +
	"\x06offset\x18\x03 \x01(\x05BK\x92AH2CThe offset (number of items skipped) used for the current response.J\x010R\x06offsetB\xba\x01\x92A\x89\x01\x12_\n" +
	"\x17Core Common Definitions\x12?Commonly used Protobuf messages for filtering, pagination, etc.2\x031.0*\x02\x01\x022\x10application/json:\x10application/jsonZ+golang-microservices-boilerplate/proto/coreb\x06proto3"

var (
	file_proto_core_common_proto_rawDescOnce sync.Once
	file_proto_core_common_proto_rawDescData []byte
)

func file_proto_core_common_proto_rawDescGZIP() []byte {
	file_proto_core_common_proto_rawDescOnce.Do(func() {
		file_proto_core_common_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_core_common_proto_rawDesc), len(file_proto_core_common_proto_rawDesc)))
	})
	return file_proto_core_common_proto_rawDescData
}

var file_proto_core_common_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_proto_core_common_proto_goTypes = []any{
	(*FilterOptions)(nil),  // 0: core.FilterOptions
	(*PaginationInfo)(nil), // 1: core.PaginationInfo
	nil,                    // 2: core.FilterOptions.FiltersEntry
	(*structpb.Value)(nil), // 3: google.protobuf.Value
}
var file_proto_core_common_proto_depIdxs = []int32{
	2, // 0: core.FilterOptions.filters:type_name -> core.FilterOptions.FiltersEntry
	3, // 1: core.FilterOptions.FiltersEntry.value:type_name -> google.protobuf.Value
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_proto_core_common_proto_init() }
func file_proto_core_common_proto_init() {
	if File_proto_core_common_proto != nil {
		return
	}
	file_proto_core_common_proto_msgTypes[0].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_core_common_proto_rawDesc), len(file_proto_core_common_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_core_common_proto_goTypes,
		DependencyIndexes: file_proto_core_common_proto_depIdxs,
		MessageInfos:      file_proto_core_common_proto_msgTypes,
	}.Build()
	File_proto_core_common_proto = out.File
	file_proto_core_common_proto_goTypes = nil
	file_proto_core_common_proto_depIdxs = nil
}
