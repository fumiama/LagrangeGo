// Code generated by protoc-gen-golite. DO NOT EDIT.
// source: pb/service/oidb/OidbSvcTrpcTcp0xFD4_1.proto

package oidb

// Fetch Friends List
type OidbSvcTrpcTcp0XFD4_1 struct {
	Field2     uint32                       `protobuf:"varint,2,opt"` // 300
	Field4     uint32                       `protobuf:"varint,4,opt"` // 0
	Field6     uint32                       `protobuf:"varint,6,opt"` // 1
	Body       []*OidbSvcTrpcTcp0XFD4_1Body `protobuf:"bytes,10001,rep"`
	Field10002 []uint32                     `protobuf:"varint,10002,rep"` // [13578, 13579, 13573, 13572, 13568]
	Field10003 uint32                       `protobuf:"varint,10003,opt"`
}

type OidbSvcTrpcTcp0XFD4_1Body struct {
	Type   uint32      `protobuf:"varint,1,opt"`
	Number *OidbNumber `protobuf:"bytes,2,opt"`
	_      [0]func()
}

type OidbSvcTrpcTcp0XFD4_1Response struct {
	DisplayFriendCount uint32        `protobuf:"varint,3,opt"`
	Timestamp          uint32        `protobuf:"varint,6,opt"`
	SelfUin            uint32        `protobuf:"varint,7,opt"`
	Friends            []*OidbFriend `protobuf:"bytes,101,rep"`
}