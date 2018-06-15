package main

import (
	"fmt"
)

//处理函数中需要根据发送给谁的id进行发送，这就需要用到session管理器，根据id查找到对于的session
//并且需要发送返回的消息， 所以可能需要传session进来
//所以消息处理模块和session和db模块有交互
//const define need uppercase for first word or all uppercase with "_" connected
const (
	ReqFrame byte = iota
	RetFrame
	//RpcFrame
	//IqFrame
	//PresenceFrame
	//JsonFrame
	//BinaryFrame
	//PingFrame
	//PongFrame
	TickFrame
	//CloseFrame

	ErrorFrame
	EchoFrame
)

//friend, presence,room, black, message, roommessage
const (
	DataType_Friend uint8 = iota
	DataType_Presence
	DataType_Group
	DataType_Room
	DataType_Black
	DataType_Message
	DataType_RoomMessage
)

//available,subscribe,subscribed,unsubscribe,unsubscribed,unavailable,invisible
const (
	PresenceType_Subscribe uint8 = iota
	PresenceType_Subscribed
	PresenceType_Unsubscribe
	PresenceType_Unsubscribed
	PresenceType_Available
	PresenceType_Unavailable
	PresenceType_Invisible
)

var msgHandler = map[uint16]func(ISession, []byte) (uint16, interface{}){}

func registerMsgHandler(msgid uint16, handler func(ISession, []byte) (uint16, interface{})) {
	_, ok := msgHandler[msgid]

	if ok {
		fmt.Println("Error: dumplicate msgid ", msgid)
		return
	}
	msgHandler[msgid] = handler
}

func HandleMsg(msgid uint16, sess ISession, buff []byte) (uint16, interface{}) {
	handler, ok := msgHandler[msgid]

	if ok {
		return handler(sess, buff)
	}
	return ERR_MSG_INVALID, nil
	//return nil, errors.New("msgid handler not exists")
}

type myint int

func (i myint) Marshal() []byte {
	return nil
}

func (i myint) UnMarshal(buff []byte) int {
	return 0
}

const MsgId_ReqLogin uint16 = 1000

type MsgReqLogin struct {
	Account  string
	Password string
}

type MsgRetLogin struct {
	//Flag      bool
	ErrorCode uint16
	Token     []byte
}

const MsgId_ReqChatLogin uint16 = 1001

type MsgReqChatLogin struct {
	Account  string
	Password string
	AppName  string
	ZoneName string
}

type MsgRetChatLogin struct {
	ErrorCode     uint16
	AppDataIdList []uint64
}

const MsgId_ReqEnterChat uint16 = 1002

type MsgReqEnterChat struct {
	AppdataId uint64
}

type MsgRetEnterChat struct {
	//Flag      bool
	ErrorCode uint16
}

const MsgId_ReqQuitChat uint16 = 1003

type MsgReqQuitChat struct {
}

type MsgRetQuitChat struct {
	//Flag      bool
	ErrorCode uint16
}

const MsgId_ReqCreateAppdata uint16 = 1004

type MsgReqCreateAppdata struct {
	NickName []byte
}

type MsgRetCreateAppdata struct {
	ErrorCode uint16
	AppdataId uint64
}

const MsgId_ReqUserData uint16 = 1005

type MsgReqUserData struct {
	AppdataId uint64
}

type MsgRetUserData struct {
	ErrorCode uint16
	Json      []byte
}

const MsgId_ReqFriendList uint16 = 1006

type MsgReqFriendList struct {
}

type MsgRetFriendList struct {
	ErrorCode uint16
	Json      []byte
}

// const MsgId_ReqUserSubscribe uint16 = 1007

// type MsgReqUserSubscribe struct {
// }

// type MsgRetUserSubscribe struct {
// 	ErrorCode uint16
// 	Json      []byte
// }

const MsgId_Presence uint16 = 1007

type MsgPresence struct {
	PresenceType uint8  `json:"presencetype"` //available,subscribe,subscribed,unsubscribe,unsubscribed,unavailable,invisible
	Who          uint64 `json:"who"`
	TimeStamp    int64  `json:"timestamp"`
	Message      []byte `json:"message"`
}

type MsgPresenceReceipt struct {
	ErrorCode uint16
}

const MsgId_Message uint16 = 1008

type MsgMessage struct {
	//MessageType uint8 //chat, friends, multi
	Who       uint64 //使用who，表示客户端填充的接收者，服务器转发时会修改为发送者
	TimeStamp int64
	Message   []byte
}

type MsgMessageReceipt struct {
	ErrorCode uint16
}

//其它类型的单人消息，服务器收到以后，转发其它人时，都是使用1008的消息格式，但是消息id使用自己的。
const MsgId_AllFriendsMessage uint16 = 1009

type MsgMsgId_AllFriendsMessage struct {
	Message []byte
}

const MsgId_GroupMessage uint16 = 1010

type MsgMsgId_GroupMessage struct {
	Count     uint8
	GroupName []byte
	Message   []byte
}

const MsgId_MultiUsersMessage uint16 = 1011

type MsgMsgId_MultiUsersMessage struct {
	Count   uint8
	To      []uint64
	Message []byte
}

const MsgId_RoomMessage uint16 = 1012

type MsgRoomMessage struct {
	Room    uint64
	From    uint64
	Message []byte
}

const MsgId_RoomUserMessage uint16 = 1013

type MsgRoomUserMessage struct {
	Room    uint64
	Who     uint64
	Message []byte
}

const MsgId_ReqDataList uint16 = 1014

type MsgReqDataList struct {
	DataType uint8 //friend, presence,room, black, message, roommessage
}

type MsgRetDataList struct {
	ErrorCode uint16
	Json      []byte
}

//create/delete user group
const MsgId_ReqGroupCreate uint16 = 1015

type MsgReqGroupCreate struct {
	GroupName []byte
}

type MsgRetGroupCreate struct {
	ErrorCode uint16
}

const MsgId_ReqGroupDelete uint16 = 1016

type MsgReqGroupDelete struct {
	GroupName []byte
}

type MsgRetGroupDelete struct {
	ErrorCode uint16
}

//add/remove black user
const MsgId_ReqAddBlack uint16 = 1017

type MsgReqAddBlack struct {
	AppdataId uint64
}

type MsgRetAddBlack struct {
	ErrorCode uint16
}

const MsgId_ReqRemoveBlack uint16 = 1018

type MsgReqRemoveBlack struct {
	AppdataId uint64
}

type MsgRetRemoveBlack struct {
	ErrorCode uint16
}

//包括从一个组移动到另一个组
const MsgId_ReqAddGroupItem uint16 = 1019

type MsgReqAddGroupItem struct {
	AppdataId uint64
	GroupName []byte
}

type MsgRetAddGroupItem struct {
	ErrorCode uint16
}

//modify user setting
const MsgId_ReqUpdateAppdata uint16 = 1020

type MsgReqUpdateAppdata struct {
	Json []byte
}

type MsgRetUpdateAppdata struct {
	ErrorCode uint16
}

//search user/room
const MsgId_ReqSearch uint16 = 1021

type MsgReqSearch struct {
	SearchType uint8
	Json       []byte
}

type MsgRetSearch struct {
	ErrorCode uint16
	Json      []byte
}

//history message ?

//modify room setting
//create/delete room
//join/quit room
//create/delete room group
//ban/unban room user
//add/remove room role
//invite/kickout
//message broadcast

//define RPC
