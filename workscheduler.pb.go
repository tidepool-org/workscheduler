// Code generated by protoc-gen-go. DO NOT EDIT.
// source: workscheduler.proto

package workscheduler

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Work - a description of the work item to process
type Work struct {
	// describes the source of the work.
	Source *WorkSource `protobuf:"bytes,1,opt,name=source,proto3" json:"source,omitempty"`
	// describes the work to perform.
	Data []byte `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	// describes directions on how to perform the work
	Directions []byte `protobuf:"bytes,3,opt,name=directions,proto3" json:"directions,omitempty"`
	// describes if the work *may* have been attempted before
	Reattempt            bool     `protobuf:"varint,4,opt,name=reattempt,proto3" json:"reattempt,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Work) Reset()         { *m = Work{} }
func (m *Work) String() string { return proto.CompactTextString(m) }
func (*Work) ProtoMessage()    {}
func (*Work) Descriptor() ([]byte, []int) {
	return fileDescriptor_ff451baf5daa9bda, []int{0}
}

func (m *Work) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Work.Unmarshal(m, b)
}
func (m *Work) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Work.Marshal(b, m, deterministic)
}
func (m *Work) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Work.Merge(m, src)
}
func (m *Work) XXX_Size() int {
	return xxx_messageInfo_Work.Size(m)
}
func (m *Work) XXX_DiscardUnknown() {
	xxx_messageInfo_Work.DiscardUnknown(m)
}

var xxx_messageInfo_Work proto.InternalMessageInfo

func (m *Work) GetSource() *WorkSource {
	if m != nil {
		return m.Source
	}
	return nil
}

func (m *Work) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *Work) GetDirections() []byte {
	if m != nil {
		return m.Directions
	}
	return nil
}

func (m *Work) GetReattempt() bool {
	if m != nil {
		return m.Reattempt
	}
	return false
}

type WorkSource struct {
	// describes the source of the work.  May encode the kafka topic and partition and offset.
	Source               string   `protobuf:"bytes,1,opt,name=source,proto3" json:"source,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *WorkSource) Reset()         { *m = WorkSource{} }
func (m *WorkSource) String() string { return proto.CompactTextString(m) }
func (*WorkSource) ProtoMessage()    {}
func (*WorkSource) Descriptor() ([]byte, []int) {
	return fileDescriptor_ff451baf5daa9bda, []int{1}
}

func (m *WorkSource) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WorkSource.Unmarshal(m, b)
}
func (m *WorkSource) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WorkSource.Marshal(b, m, deterministic)
}
func (m *WorkSource) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WorkSource.Merge(m, src)
}
func (m *WorkSource) XXX_Size() int {
	return xxx_messageInfo_WorkSource.Size(m)
}
func (m *WorkSource) XXX_DiscardUnknown() {
	xxx_messageInfo_WorkSource.DiscardUnknown(m)
}

var xxx_messageInfo_WorkSource proto.InternalMessageInfo

func (m *WorkSource) GetSource() string {
	if m != nil {
		return m.Source
	}
	return ""
}

// LagResponse - how much the work scheduler is lagging
type LagResponse struct {
	// the timestamp of the oldest message not yet processed
	MyField              *timestamp.Timestamp `protobuf:"bytes,1,opt,name=my_field,json=myField,proto3" json:"my_field,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *LagResponse) Reset()         { *m = LagResponse{} }
func (m *LagResponse) String() string { return proto.CompactTextString(m) }
func (*LagResponse) ProtoMessage()    {}
func (*LagResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_ff451baf5daa9bda, []int{2}
}

func (m *LagResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LagResponse.Unmarshal(m, b)
}
func (m *LagResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LagResponse.Marshal(b, m, deterministic)
}
func (m *LagResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LagResponse.Merge(m, src)
}
func (m *LagResponse) XXX_Size() int {
	return xxx_messageInfo_LagResponse.Size(m)
}
func (m *LagResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_LagResponse.DiscardUnknown(m)
}

var xxx_messageInfo_LagResponse proto.InternalMessageInfo

func (m *LagResponse) GetMyField() *timestamp.Timestamp {
	if m != nil {
		return m.MyField
	}
	return nil
}

type WorkOutput struct {
	Id                   *KafkaOutput `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *WorkOutput) Reset()         { *m = WorkOutput{} }
func (m *WorkOutput) String() string { return proto.CompactTextString(m) }
func (*WorkOutput) ProtoMessage()    {}
func (*WorkOutput) Descriptor() ([]byte, []int) {
	return fileDescriptor_ff451baf5daa9bda, []int{3}
}

func (m *WorkOutput) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WorkOutput.Unmarshal(m, b)
}
func (m *WorkOutput) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WorkOutput.Marshal(b, m, deterministic)
}
func (m *WorkOutput) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WorkOutput.Merge(m, src)
}
func (m *WorkOutput) XXX_Size() int {
	return xxx_messageInfo_WorkOutput.Size(m)
}
func (m *WorkOutput) XXX_DiscardUnknown() {
	xxx_messageInfo_WorkOutput.DiscardUnknown(m)
}

var xxx_messageInfo_WorkOutput proto.InternalMessageInfo

func (m *WorkOutput) GetId() *KafkaOutput {
	if m != nil {
		return m.Id
	}
	return nil
}

type KafkaOutput struct {
	KafkaTopicMessages   []*KafkaTopicMessages `protobuf:"bytes,1,rep,name=kafkaTopicMessages,proto3" json:"kafkaTopicMessages,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *KafkaOutput) Reset()         { *m = KafkaOutput{} }
func (m *KafkaOutput) String() string { return proto.CompactTextString(m) }
func (*KafkaOutput) ProtoMessage()    {}
func (*KafkaOutput) Descriptor() ([]byte, []int) {
	return fileDescriptor_ff451baf5daa9bda, []int{4}
}

func (m *KafkaOutput) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_KafkaOutput.Unmarshal(m, b)
}
func (m *KafkaOutput) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_KafkaOutput.Marshal(b, m, deterministic)
}
func (m *KafkaOutput) XXX_Merge(src proto.Message) {
	xxx_messageInfo_KafkaOutput.Merge(m, src)
}
func (m *KafkaOutput) XXX_Size() int {
	return xxx_messageInfo_KafkaOutput.Size(m)
}
func (m *KafkaOutput) XXX_DiscardUnknown() {
	xxx_messageInfo_KafkaOutput.DiscardUnknown(m)
}

var xxx_messageInfo_KafkaOutput proto.InternalMessageInfo

func (m *KafkaOutput) GetKafkaTopicMessages() []*KafkaTopicMessages {
	if m != nil {
		return m.KafkaTopicMessages
	}
	return nil
}

type KafkaTopicMessages struct {
	TopicName            string                    `protobuf:"bytes,1,opt,name=topicName,proto3" json:"topicName,omitempty"`
	PartitionMessages    []*KafkaPartitionMessages `protobuf:"bytes,2,rep,name=partitionMessages,proto3" json:"partitionMessages,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                  `json:"-"`
	XXX_unrecognized     []byte                    `json:"-"`
	XXX_sizecache        int32                     `json:"-"`
}

func (m *KafkaTopicMessages) Reset()         { *m = KafkaTopicMessages{} }
func (m *KafkaTopicMessages) String() string { return proto.CompactTextString(m) }
func (*KafkaTopicMessages) ProtoMessage()    {}
func (*KafkaTopicMessages) Descriptor() ([]byte, []int) {
	return fileDescriptor_ff451baf5daa9bda, []int{5}
}

func (m *KafkaTopicMessages) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_KafkaTopicMessages.Unmarshal(m, b)
}
func (m *KafkaTopicMessages) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_KafkaTopicMessages.Marshal(b, m, deterministic)
}
func (m *KafkaTopicMessages) XXX_Merge(src proto.Message) {
	xxx_messageInfo_KafkaTopicMessages.Merge(m, src)
}
func (m *KafkaTopicMessages) XXX_Size() int {
	return xxx_messageInfo_KafkaTopicMessages.Size(m)
}
func (m *KafkaTopicMessages) XXX_DiscardUnknown() {
	xxx_messageInfo_KafkaTopicMessages.DiscardUnknown(m)
}

var xxx_messageInfo_KafkaTopicMessages proto.InternalMessageInfo

func (m *KafkaTopicMessages) GetTopicName() string {
	if m != nil {
		return m.TopicName
	}
	return ""
}

func (m *KafkaTopicMessages) GetPartitionMessages() []*KafkaPartitionMessages {
	if m != nil {
		return m.PartitionMessages
	}
	return nil
}

type KafkaPartitionMessages struct {
	Partition            int32           `protobuf:"varint,1,opt,name=partition,proto3" json:"partition,omitempty"`
	Messages             []*KafkaMessage `protobuf:"bytes,2,rep,name=messages,proto3" json:"messages,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *KafkaPartitionMessages) Reset()         { *m = KafkaPartitionMessages{} }
func (m *KafkaPartitionMessages) String() string { return proto.CompactTextString(m) }
func (*KafkaPartitionMessages) ProtoMessage()    {}
func (*KafkaPartitionMessages) Descriptor() ([]byte, []int) {
	return fileDescriptor_ff451baf5daa9bda, []int{6}
}

func (m *KafkaPartitionMessages) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_KafkaPartitionMessages.Unmarshal(m, b)
}
func (m *KafkaPartitionMessages) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_KafkaPartitionMessages.Marshal(b, m, deterministic)
}
func (m *KafkaPartitionMessages) XXX_Merge(src proto.Message) {
	xxx_messageInfo_KafkaPartitionMessages.Merge(m, src)
}
func (m *KafkaPartitionMessages) XXX_Size() int {
	return xxx_messageInfo_KafkaPartitionMessages.Size(m)
}
func (m *KafkaPartitionMessages) XXX_DiscardUnknown() {
	xxx_messageInfo_KafkaPartitionMessages.DiscardUnknown(m)
}

var xxx_messageInfo_KafkaPartitionMessages proto.InternalMessageInfo

func (m *KafkaPartitionMessages) GetPartition() int32 {
	if m != nil {
		return m.Partition
	}
	return 0
}

func (m *KafkaPartitionMessages) GetMessages() []*KafkaMessage {
	if m != nil {
		return m.Messages
	}
	return nil
}

type KafkaMessage struct {
	Data                 []byte   `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *KafkaMessage) Reset()         { *m = KafkaMessage{} }
func (m *KafkaMessage) String() string { return proto.CompactTextString(m) }
func (*KafkaMessage) ProtoMessage()    {}
func (*KafkaMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_ff451baf5daa9bda, []int{7}
}

func (m *KafkaMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_KafkaMessage.Unmarshal(m, b)
}
func (m *KafkaMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_KafkaMessage.Marshal(b, m, deterministic)
}
func (m *KafkaMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_KafkaMessage.Merge(m, src)
}
func (m *KafkaMessage) XXX_Size() int {
	return xxx_messageInfo_KafkaMessage.Size(m)
}
func (m *KafkaMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_KafkaMessage.DiscardUnknown(m)
}

var xxx_messageInfo_KafkaMessage proto.InternalMessageInfo

func (m *KafkaMessage) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterType((*Work)(nil), "workscheduler.Work")
	proto.RegisterType((*WorkSource)(nil), "workscheduler.WorkSource")
	proto.RegisterType((*LagResponse)(nil), "workscheduler.LagResponse")
	proto.RegisterType((*WorkOutput)(nil), "workscheduler.WorkOutput")
	proto.RegisterType((*KafkaOutput)(nil), "workscheduler.KafkaOutput")
	proto.RegisterType((*KafkaTopicMessages)(nil), "workscheduler.KafkaTopicMessages")
	proto.RegisterType((*KafkaPartitionMessages)(nil), "workscheduler.KafkaPartitionMessages")
	proto.RegisterType((*KafkaMessage)(nil), "workscheduler.KafkaMessage")
}

func init() {
	proto.RegisterFile("workscheduler.proto", fileDescriptor_ff451baf5daa9bda)
}

var fileDescriptor_ff451baf5daa9bda = []byte{
	// 509 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x93, 0xc1, 0x6e, 0xd3, 0x40,
	0x10, 0x86, 0xeb, 0xc4, 0x84, 0x74, 0xd2, 0x1e, 0x98, 0x4a, 0x91, 0x49, 0x11, 0x84, 0x15, 0x48,
	0x56, 0x0f, 0x8e, 0x08, 0xaa, 0x88, 0x84, 0xb8, 0x14, 0xe8, 0x85, 0x02, 0xe9, 0xa6, 0x12, 0x12,
	0x17, 0xd8, 0xc6, 0x1b, 0xb3, 0x8a, 0x9d, 0x5d, 0xd9, 0x6b, 0xa1, 0x3c, 0x01, 0x07, 0x5e, 0x87,
	0x07, 0x44, 0x6b, 0xc7, 0x71, 0x9c, 0x38, 0x3d, 0xf4, 0xe6, 0xfd, 0xe7, 0x9b, 0x99, 0x7f, 0xc7,
	0xb3, 0x70, 0xf2, 0x5b, 0xc6, 0xf3, 0x64, 0xfa, 0x8b, 0xfb, 0x69, 0xc8, 0x63, 0x4f, 0xc5, 0x52,
	0x4b, 0x3c, 0xae, 0x88, 0xbd, 0x67, 0x81, 0x94, 0x41, 0xc8, 0x07, 0x59, 0xf0, 0x36, 0x9d, 0x0d,
	0xb4, 0x88, 0x78, 0xa2, 0x59, 0xa4, 0x72, 0xbe, 0x77, 0xba, 0x0d, 0xf0, 0x48, 0xe9, 0x65, 0x1e,
	0x24, 0x7f, 0x2d, 0xb0, 0xbf, 0xc9, 0x78, 0x8e, 0xaf, 0xa0, 0x95, 0xc8, 0x34, 0x9e, 0x72, 0xc7,
	0xea, 0x5b, 0x6e, 0x67, 0xf8, 0xd8, 0xab, 0xf6, 0x36, 0xd0, 0x24, 0x03, 0xe8, 0x0a, 0x44, 0x04,
	0xdb, 0x67, 0x9a, 0x39, 0x8d, 0xbe, 0xe5, 0x1e, 0xd1, 0xec, 0x1b, 0x9f, 0x02, 0xf8, 0x22, 0xe6,
	0x53, 0x2d, 0xe4, 0x22, 0x71, 0x9a, 0x59, 0x64, 0x43, 0xc1, 0x27, 0x70, 0x18, 0x73, 0xa6, 0xb5,
	0xf1, 0xe0, 0xd8, 0x7d, 0xcb, 0x6d, 0xd3, 0x52, 0x20, 0x2f, 0x00, 0xca, 0x3e, 0xd8, 0xad, 0x58,
	0x3a, 0x2c, 0xfa, 0x92, 0x0f, 0xd0, 0xb9, 0x62, 0x01, 0xe5, 0x89, 0x92, 0x8b, 0x84, 0xe3, 0x39,
	0xb4, 0xa3, 0xe5, 0x8f, 0x99, 0xe0, 0xa1, 0xbf, 0xf2, 0xde, 0xf3, 0xf2, 0x2b, 0x7b, 0xc5, 0x95,
	0xbd, 0x9b, 0x62, 0x26, 0xf4, 0x61, 0xb4, 0xbc, 0x34, 0x28, 0x19, 0xe5, 0xbd, 0xbe, 0xa6, 0x5a,
	0xa5, 0x1a, 0xcf, 0xa0, 0x21, 0xca, 0xf4, 0xea, 0xd5, 0x3f, 0xb1, 0xd9, 0x9c, 0xe5, 0x1c, 0x6d,
	0x08, 0x9f, 0xfc, 0x84, 0xce, 0x86, 0x84, 0xd7, 0x80, 0x73, 0x73, 0xbc, 0x91, 0x4a, 0x4c, 0x3f,
	0xf3, 0x24, 0x61, 0x01, 0x4f, 0x1c, 0xab, 0xdf, 0x74, 0x3b, 0xc3, 0xe7, 0x75, 0xa5, 0x2a, 0x20,
	0xad, 0x49, 0x26, 0x7f, 0x2c, 0xc0, 0x5d, 0xd4, 0x0c, 0x4f, 0x1b, 0xe1, 0x0b, 0x8b, 0x8a, 0x99,
	0x94, 0x02, 0x4e, 0xe0, 0x91, 0x62, 0xb1, 0x16, 0x66, 0xd0, 0x6b, 0x1b, 0x8d, 0xcc, 0xc6, 0xcb,
	0x3a, 0x1b, 0xe3, 0x6d, 0x98, 0xee, 0xe6, 0x13, 0x09, 0xdd, 0x7a, 0xd8, 0x98, 0x59, 0xe3, 0x99,
	0x99, 0x07, 0xb4, 0x14, 0xf0, 0x0d, 0xb4, 0xa3, 0xaa, 0x87, 0xd3, 0x3a, 0x0f, 0xab, 0x6a, 0x74,
	0x0d, 0x13, 0x02, 0x47, 0x9b, 0x91, 0xf5, 0x92, 0x59, 0xe5, 0x92, 0x0d, 0xff, 0x35, 0xe1, 0x38,
	0xdb, 0x93, 0xa2, 0x18, 0x9e, 0x83, 0x3d, 0x96, 0x61, 0x88, 0xdd, 0x9d, 0x3f, 0xff, 0xd1, 0x2c,
	0x7b, 0xef, 0xa4, 0x66, 0x9b, 0xc9, 0x01, 0x8e, 0xc0, 0x1e, 0x8b, 0x45, 0xb0, 0x37, 0x6d, 0x8f,
	0x4e, 0x0e, 0xf0, 0x1d, 0xb4, 0x2e, 0x99, 0x08, 0xb9, 0x8f, 0xfb, 0x1f, 0xca, 0x1d, 0xe9, 0x17,
	0xd0, 0x7e, 0x2f, 0x23, 0x15, 0x72, 0xcd, 0xef, 0x2a, 0x50, 0x17, 0xca, 0xb7, 0x2e, 0x37, 0x7f,
	0x9d, 0x0a, 0x7d, 0x0f, 0xf3, 0x23, 0xb0, 0x27, 0x5a, 0xaa, 0x7b, 0x64, 0xbe, 0x85, 0xe6, 0x15,
	0xdb, 0x3f, 0xaf, 0xed, 0x97, 0xb3, 0xf1, 0x4c, 0xc9, 0xc1, 0xc5, 0xd9, 0x77, 0x57, 0x0b, 0x9f,
	0x2b, 0x29, 0x43, 0x4f, 0xc6, 0xc1, 0xa0, 0xc2, 0x56, 0x4f, 0xb7, 0xad, 0xac, 0xf2, 0xeb, 0xff,
	0x01, 0x00, 0x00, 0xff, 0xff, 0x38, 0xb4, 0x6f, 0x3b, 0x02, 0x05, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// WorkSchedulerClient is the client API for WorkScheduler service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type WorkSchedulerClient interface {
	// Retrieve work to be performed
	Poll(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*Work, error)
	// confirm that work scheduler is alive
	Ping(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*empty.Empty, error)
	// report to work scheduler that a work item failed
	Failed(ctx context.Context, in *WorkSource, opts ...grpc.CallOption) (*empty.Empty, error)
	// report to work scheduler that a work item is completed
	Complete(ctx context.Context, in *WorkSource, opts ...grpc.CallOption) (*WorkOutput, error)
	// cause work scheduler to exit
	Quit(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*empty.Empty, error)
	// cause work scheduler to stop providing work
	Stop(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*empty.Empty, error)
	Lag(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*LagResponse, error)
}

type workSchedulerClient struct {
	cc grpc.ClientConnInterface
}

func NewWorkSchedulerClient(cc grpc.ClientConnInterface) WorkSchedulerClient {
	return &workSchedulerClient{cc}
}

func (c *workSchedulerClient) Poll(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*Work, error) {
	out := new(Work)
	err := c.cc.Invoke(ctx, "/workscheduler.WorkScheduler/Poll", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *workSchedulerClient) Ping(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/workscheduler.WorkScheduler/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *workSchedulerClient) Failed(ctx context.Context, in *WorkSource, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/workscheduler.WorkScheduler/Failed", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *workSchedulerClient) Complete(ctx context.Context, in *WorkSource, opts ...grpc.CallOption) (*WorkOutput, error) {
	out := new(WorkOutput)
	err := c.cc.Invoke(ctx, "/workscheduler.WorkScheduler/Complete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *workSchedulerClient) Quit(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/workscheduler.WorkScheduler/Quit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *workSchedulerClient) Stop(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/workscheduler.WorkScheduler/Stop", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *workSchedulerClient) Lag(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*LagResponse, error) {
	out := new(LagResponse)
	err := c.cc.Invoke(ctx, "/workscheduler.WorkScheduler/Lag", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// WorkSchedulerServer is the server API for WorkScheduler service.
type WorkSchedulerServer interface {
	// Retrieve work to be performed
	Poll(context.Context, *empty.Empty) (*Work, error)
	// confirm that work scheduler is alive
	Ping(context.Context, *empty.Empty) (*empty.Empty, error)
	// report to work scheduler that a work item failed
	Failed(context.Context, *WorkSource) (*empty.Empty, error)
	// report to work scheduler that a work item is completed
	Complete(context.Context, *WorkSource) (*WorkOutput, error)
	// cause work scheduler to exit
	Quit(context.Context, *empty.Empty) (*empty.Empty, error)
	// cause work scheduler to stop providing work
	Stop(context.Context, *empty.Empty) (*empty.Empty, error)
	Lag(context.Context, *empty.Empty) (*LagResponse, error)
}

// UnimplementedWorkSchedulerServer can be embedded to have forward compatible implementations.
type UnimplementedWorkSchedulerServer struct {
}

func (*UnimplementedWorkSchedulerServer) Poll(ctx context.Context, req *empty.Empty) (*Work, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Poll not implemented")
}
func (*UnimplementedWorkSchedulerServer) Ping(ctx context.Context, req *empty.Empty) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (*UnimplementedWorkSchedulerServer) Failed(ctx context.Context, req *WorkSource) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Failed not implemented")
}
func (*UnimplementedWorkSchedulerServer) Complete(ctx context.Context, req *WorkSource) (*WorkOutput, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Complete not implemented")
}
func (*UnimplementedWorkSchedulerServer) Quit(ctx context.Context, req *empty.Empty) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Quit not implemented")
}
func (*UnimplementedWorkSchedulerServer) Stop(ctx context.Context, req *empty.Empty) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Stop not implemented")
}
func (*UnimplementedWorkSchedulerServer) Lag(ctx context.Context, req *empty.Empty) (*LagResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Lag not implemented")
}

func RegisterWorkSchedulerServer(s *grpc.Server, srv WorkSchedulerServer) {
	s.RegisterService(&_WorkScheduler_serviceDesc, srv)
}

func _WorkScheduler_Poll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorkSchedulerServer).Poll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/workscheduler.WorkScheduler/Poll",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorkSchedulerServer).Poll(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _WorkScheduler_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorkSchedulerServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/workscheduler.WorkScheduler/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorkSchedulerServer).Ping(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _WorkScheduler_Failed_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WorkSource)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorkSchedulerServer).Failed(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/workscheduler.WorkScheduler/Failed",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorkSchedulerServer).Failed(ctx, req.(*WorkSource))
	}
	return interceptor(ctx, in, info, handler)
}

func _WorkScheduler_Complete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WorkSource)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorkSchedulerServer).Complete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/workscheduler.WorkScheduler/Complete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorkSchedulerServer).Complete(ctx, req.(*WorkSource))
	}
	return interceptor(ctx, in, info, handler)
}

func _WorkScheduler_Quit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorkSchedulerServer).Quit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/workscheduler.WorkScheduler/Quit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorkSchedulerServer).Quit(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _WorkScheduler_Stop_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorkSchedulerServer).Stop(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/workscheduler.WorkScheduler/Stop",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorkSchedulerServer).Stop(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _WorkScheduler_Lag_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorkSchedulerServer).Lag(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/workscheduler.WorkScheduler/Lag",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorkSchedulerServer).Lag(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _WorkScheduler_serviceDesc = grpc.ServiceDesc{
	ServiceName: "workscheduler.WorkScheduler",
	HandlerType: (*WorkSchedulerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Poll",
			Handler:    _WorkScheduler_Poll_Handler,
		},
		{
			MethodName: "Ping",
			Handler:    _WorkScheduler_Ping_Handler,
		},
		{
			MethodName: "Failed",
			Handler:    _WorkScheduler_Failed_Handler,
		},
		{
			MethodName: "Complete",
			Handler:    _WorkScheduler_Complete_Handler,
		},
		{
			MethodName: "Quit",
			Handler:    _WorkScheduler_Quit_Handler,
		},
		{
			MethodName: "Stop",
			Handler:    _WorkScheduler_Stop_Handler,
		},
		{
			MethodName: "Lag",
			Handler:    _WorkScheduler_Lag_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "workscheduler.proto",
}
