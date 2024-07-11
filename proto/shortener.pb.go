// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.1
// source: proto/shortener.proto

package proto

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

type ShortenLink struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ShortenLink string `protobuf:"bytes,1,opt,name=shortenLink,proto3" json:"shortenLink,omitempty"`
}

func (x *ShortenLink) Reset() {
	*x = ShortenLink{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_shortener_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ShortenLink) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShortenLink) ProtoMessage() {}

func (x *ShortenLink) ProtoReflect() protoreflect.Message {
	mi := &file_proto_shortener_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShortenLink.ProtoReflect.Descriptor instead.
func (*ShortenLink) Descriptor() ([]byte, []int) {
	return file_proto_shortener_proto_rawDescGZIP(), []int{0}
}

func (x *ShortenLink) GetShortenLink() string {
	if x != nil {
		return x.ShortenLink
	}
	return ""
}

type LongLink struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LongLink string `protobuf:"bytes,1,opt,name=longLink,proto3" json:"longLink,omitempty"`
}

func (x *LongLink) Reset() {
	*x = LongLink{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_shortener_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LongLink) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LongLink) ProtoMessage() {}

func (x *LongLink) ProtoReflect() protoreflect.Message {
	mi := &file_proto_shortener_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LongLink.ProtoReflect.Descriptor instead.
func (*LongLink) Descriptor() ([]byte, []int) {
	return file_proto_shortener_proto_rawDescGZIP(), []int{1}
}

func (x *LongLink) GetLongLink() string {
	if x != nil {
		return x.LongLink
	}
	return ""
}

type UserLinks struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LongLink  string `protobuf:"bytes,1,opt,name=LongLink,proto3" json:"LongLink,omitempty"`
	ShortLink string `protobuf:"bytes,2,opt,name=ShortLink,proto3" json:"ShortLink,omitempty"`
}

func (x *UserLinks) Reset() {
	*x = UserLinks{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_shortener_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserLinks) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserLinks) ProtoMessage() {}

func (x *UserLinks) ProtoReflect() protoreflect.Message {
	mi := &file_proto_shortener_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserLinks.ProtoReflect.Descriptor instead.
func (*UserLinks) Descriptor() ([]byte, []int) {
	return file_proto_shortener_proto_rawDescGZIP(), []int{2}
}

func (x *UserLinks) GetLongLink() string {
	if x != nil {
		return x.LongLink
	}
	return ""
}

func (x *UserLinks) GetShortLink() string {
	if x != nil {
		return x.ShortLink
	}
	return ""
}

type ListShortenLinks struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserLinks string `protobuf:"bytes,1,opt,name=userLinks,proto3" json:"userLinks,omitempty"`
}

func (x *ListShortenLinks) Reset() {
	*x = ListShortenLinks{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_shortener_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListShortenLinks) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListShortenLinks) ProtoMessage() {}

func (x *ListShortenLinks) ProtoReflect() protoreflect.Message {
	mi := &file_proto_shortener_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListShortenLinks.ProtoReflect.Descriptor instead.
func (*ListShortenLinks) Descriptor() ([]byte, []int) {
	return file_proto_shortener_proto_rawDescGZIP(), []int{3}
}

func (x *ListShortenLinks) GetUserLinks() string {
	if x != nil {
		return x.UserLinks
	}
	return ""
}

type ListShortenLinksToDelete struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserLinks []string `protobuf:"bytes,1,rep,name=userLinks,proto3" json:"userLinks,omitempty"`
}

func (x *ListShortenLinksToDelete) Reset() {
	*x = ListShortenLinksToDelete{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_shortener_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListShortenLinksToDelete) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListShortenLinksToDelete) ProtoMessage() {}

func (x *ListShortenLinksToDelete) ProtoReflect() protoreflect.Message {
	mi := &file_proto_shortener_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListShortenLinksToDelete.ProtoReflect.Descriptor instead.
func (*ListShortenLinksToDelete) Descriptor() ([]byte, []int) {
	return file_proto_shortener_proto_rawDescGZIP(), []int{4}
}

func (x *ListShortenLinksToDelete) GetUserLinks() []string {
	if x != nil {
		return x.UserLinks
	}
	return nil
}

type FindShortenLinkResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LongLink     string `protobuf:"bytes,1,opt,name=longLink,proto3" json:"longLink,omitempty"`
	DeleteStatus bool   `protobuf:"varint,2,opt,name=deleteStatus,proto3" json:"deleteStatus,omitempty"`
}

func (x *FindShortenLinkResponse) Reset() {
	*x = FindShortenLinkResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_shortener_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FindShortenLinkResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FindShortenLinkResponse) ProtoMessage() {}

func (x *FindShortenLinkResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_shortener_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FindShortenLinkResponse.ProtoReflect.Descriptor instead.
func (*FindShortenLinkResponse) Descriptor() ([]byte, []int) {
	return file_proto_shortener_proto_rawDescGZIP(), []int{5}
}

func (x *FindShortenLinkResponse) GetLongLink() string {
	if x != nil {
		return x.LongLink
	}
	return ""
}

func (x *FindShortenLinkResponse) GetDeleteStatus() bool {
	if x != nil {
		return x.DeleteStatus
	}
	return false
}

type ShortenLinkResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LongLink     string `protobuf:"bytes,1,opt,name=longLink,proto3" json:"longLink,omitempty"`
	DeleteStatus bool   `protobuf:"varint,2,opt,name=deleteStatus,proto3" json:"deleteStatus,omitempty"`
}

func (x *ShortenLinkResponse) Reset() {
	*x = ShortenLinkResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_shortener_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ShortenLinkResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShortenLinkResponse) ProtoMessage() {}

func (x *ShortenLinkResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_shortener_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShortenLinkResponse.ProtoReflect.Descriptor instead.
func (*ShortenLinkResponse) Descriptor() ([]byte, []int) {
	return file_proto_shortener_proto_rawDescGZIP(), []int{6}
}

func (x *ShortenLinkResponse) GetLongLink() string {
	if x != nil {
		return x.LongLink
	}
	return ""
}

func (x *ShortenLinkResponse) GetDeleteStatus() bool {
	if x != nil {
		return x.DeleteStatus
	}
	return false
}

type LongLinkResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ShortenLink string `protobuf:"bytes,1,opt,name=shortenLink,proto3" json:"shortenLink,omitempty"`
}

func (x *LongLinkResponse) Reset() {
	*x = LongLinkResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_shortener_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LongLinkResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LongLinkResponse) ProtoMessage() {}

func (x *LongLinkResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_shortener_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LongLinkResponse.ProtoReflect.Descriptor instead.
func (*LongLinkResponse) Descriptor() ([]byte, []int) {
	return file_proto_shortener_proto_rawDescGZIP(), []int{7}
}

func (x *LongLinkResponse) GetShortenLink() string {
	if x != nil {
		return x.ShortenLink
	}
	return ""
}

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_shortener_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_proto_shortener_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_proto_shortener_proto_rawDescGZIP(), []int{8}
}

var File_proto_shortener_proto protoreflect.FileDescriptor

var file_proto_shortener_proto_rawDesc = []byte{
	0x0a, 0x15, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65,
	0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x2f,
	0x0a, 0x0b, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x4c, 0x69, 0x6e, 0x6b, 0x12, 0x20, 0x0a,
	0x0b, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x4c, 0x69, 0x6e, 0x6b, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0b, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x4c, 0x69, 0x6e, 0x6b, 0x22,
	0x26, 0x0a, 0x08, 0x4c, 0x6f, 0x6e, 0x67, 0x4c, 0x69, 0x6e, 0x6b, 0x12, 0x1a, 0x0a, 0x08, 0x6c,
	0x6f, 0x6e, 0x67, 0x4c, 0x69, 0x6e, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6c,
	0x6f, 0x6e, 0x67, 0x4c, 0x69, 0x6e, 0x6b, 0x22, 0x45, 0x0a, 0x09, 0x55, 0x73, 0x65, 0x72, 0x4c,
	0x69, 0x6e, 0x6b, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x4c, 0x6f, 0x6e, 0x67, 0x4c, 0x69, 0x6e, 0x6b,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x4c, 0x6f, 0x6e, 0x67, 0x4c, 0x69, 0x6e, 0x6b,
	0x12, 0x1c, 0x0a, 0x09, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x4c, 0x69, 0x6e, 0x6b, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x4c, 0x69, 0x6e, 0x6b, 0x22, 0x30,
	0x0a, 0x10, 0x4c, 0x69, 0x73, 0x74, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x4c, 0x69, 0x6e,
	0x6b, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x75, 0x73, 0x65, 0x72, 0x4c, 0x69, 0x6e, 0x6b, 0x73, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x75, 0x73, 0x65, 0x72, 0x4c, 0x69, 0x6e, 0x6b, 0x73,
	0x22, 0x38, 0x0a, 0x18, 0x4c, 0x69, 0x73, 0x74, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x4c,
	0x69, 0x6e, 0x6b, 0x73, 0x54, 0x6f, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x1c, 0x0a, 0x09,
	0x75, 0x73, 0x65, 0x72, 0x4c, 0x69, 0x6e, 0x6b, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52,
	0x09, 0x75, 0x73, 0x65, 0x72, 0x4c, 0x69, 0x6e, 0x6b, 0x73, 0x22, 0x59, 0x0a, 0x17, 0x46, 0x69,
	0x6e, 0x64, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x4c, 0x69, 0x6e, 0x6b, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x6f, 0x6e, 0x67, 0x4c, 0x69, 0x6e,
	0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6c, 0x6f, 0x6e, 0x67, 0x4c, 0x69, 0x6e,
	0x6b, 0x12, 0x22, 0x0a, 0x0c, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0c, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x55, 0x0a, 0x13, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e,
	0x4c, 0x69, 0x6e, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1a, 0x0a, 0x08,
	0x6c, 0x6f, 0x6e, 0x67, 0x4c, 0x69, 0x6e, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x6c, 0x6f, 0x6e, 0x67, 0x4c, 0x69, 0x6e, 0x6b, 0x12, 0x22, 0x0a, 0x0c, 0x64, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0c,
	0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x34, 0x0a, 0x10,
	0x4c, 0x6f, 0x6e, 0x67, 0x4c, 0x69, 0x6e, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x20, 0x0a, 0x0b, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x4c, 0x69, 0x6e, 0x6b, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x4c, 0x69,
	0x6e, 0x6b, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x32, 0xd7, 0x01, 0x0a, 0x05,
	0x4c, 0x69, 0x6e, 0x6b, 0x73, 0x12, 0x35, 0x0a, 0x03, 0x47, 0x65, 0x74, 0x12, 0x12, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x4c, 0x69, 0x6e, 0x6b,
	0x1a, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e,
	0x4c, 0x69, 0x6e, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x30, 0x0a, 0x04,
	0x53, 0x61, 0x76, 0x65, 0x12, 0x0f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4c, 0x6f, 0x6e,
	0x67, 0x4c, 0x69, 0x6e, 0x6b, 0x1a, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4c, 0x6f,
	0x6e, 0x67, 0x4c, 0x69, 0x6e, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2f,
	0x0a, 0x06, 0x47, 0x65, 0x74, 0x41, 0x6c, 0x6c, 0x12, 0x0c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4c,
	0x69, 0x73, 0x74, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x4c, 0x69, 0x6e, 0x6b, 0x73, 0x12,
	0x34, 0x0a, 0x03, 0x44, 0x65, 0x6c, 0x12, 0x1f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4c,
	0x69, 0x73, 0x74, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x4c, 0x69, 0x6e, 0x6b, 0x73, 0x54,
	0x6f, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x1a, 0x0c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x42, 0x1a, 0x5a, 0x18, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x6e, 0x65, 0x78, 0x74, 0x6c, 0x61, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_shortener_proto_rawDescOnce sync.Once
	file_proto_shortener_proto_rawDescData = file_proto_shortener_proto_rawDesc
)

func file_proto_shortener_proto_rawDescGZIP() []byte {
	file_proto_shortener_proto_rawDescOnce.Do(func() {
		file_proto_shortener_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_shortener_proto_rawDescData)
	})
	return file_proto_shortener_proto_rawDescData
}

var file_proto_shortener_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_proto_shortener_proto_goTypes = []any{
	(*ShortenLink)(nil),              // 0: proto.ShortenLink
	(*LongLink)(nil),                 // 1: proto.LongLink
	(*UserLinks)(nil),                // 2: proto.UserLinks
	(*ListShortenLinks)(nil),         // 3: proto.ListShortenLinks
	(*ListShortenLinksToDelete)(nil), // 4: proto.ListShortenLinksToDelete
	(*FindShortenLinkResponse)(nil),  // 5: proto.FindShortenLinkResponse
	(*ShortenLinkResponse)(nil),      // 6: proto.ShortenLinkResponse
	(*LongLinkResponse)(nil),         // 7: proto.LongLinkResponse
	(*Empty)(nil),                    // 8: proto.Empty
}
var file_proto_shortener_proto_depIdxs = []int32{
	0, // 0: proto.Links.Get:input_type -> proto.ShortenLink
	1, // 1: proto.Links.Save:input_type -> proto.LongLink
	8, // 2: proto.Links.GetAll:input_type -> proto.Empty
	4, // 3: proto.Links.Del:input_type -> proto.ListShortenLinksToDelete
	6, // 4: proto.Links.Get:output_type -> proto.ShortenLinkResponse
	7, // 5: proto.Links.Save:output_type -> proto.LongLinkResponse
	3, // 6: proto.Links.GetAll:output_type -> proto.ListShortenLinks
	8, // 7: proto.Links.Del:output_type -> proto.Empty
	4, // [4:8] is the sub-list for method output_type
	0, // [0:4] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_shortener_proto_init() }
func file_proto_shortener_proto_init() {
	if File_proto_shortener_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_shortener_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*ShortenLink); i {
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
		file_proto_shortener_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*LongLink); i {
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
		file_proto_shortener_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*UserLinks); i {
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
		file_proto_shortener_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*ListShortenLinks); i {
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
		file_proto_shortener_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*ListShortenLinksToDelete); i {
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
		file_proto_shortener_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*FindShortenLinkResponse); i {
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
		file_proto_shortener_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*ShortenLinkResponse); i {
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
		file_proto_shortener_proto_msgTypes[7].Exporter = func(v any, i int) any {
			switch v := v.(*LongLinkResponse); i {
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
		file_proto_shortener_proto_msgTypes[8].Exporter = func(v any, i int) any {
			switch v := v.(*Empty); i {
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
			RawDescriptor: file_proto_shortener_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_shortener_proto_goTypes,
		DependencyIndexes: file_proto_shortener_proto_depIdxs,
		MessageInfos:      file_proto_shortener_proto_msgTypes,
	}.Build()
	File_proto_shortener_proto = out.File
	file_proto_shortener_proto_rawDesc = nil
	file_proto_shortener_proto_goTypes = nil
	file_proto_shortener_proto_depIdxs = nil
}
