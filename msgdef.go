package main

import (
	"fmt"
)

//客户端在发送一个请求以后，需要启动一个定时器来检测该请求是否得到反馈：
//1.如果在规定时间没有反馈，则说明网络超时。
//2.如果收到反馈，则可以根据反馈进行弹框告知用户。
//在收到反馈或者超时以后需要将该检测移除检测队列。

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
	//DataType_Group
	DataType_Presence
	DataType_Room
	DataType_Black
	DataType_Offline_Message
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
	Platform string
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

type MsgReqGroupFriendList struct {
	GroupName []byte
}

type MsgRetGroupFriendList struct {
	ErrorCode uint16
	Count     uint16
	Step      uint16
	Data      []byte
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
	Who          uint64 `json:"who,string"`
	Nickname     string `json:"nickname"`
	TimeStamp    int64  `json:"timestamp,string"`
	Message      string `json:"message"`
}

type MsgPresenceReceipt struct {
	ErrorCode uint16
}

const MsgId_Message uint16 = 1008

type MsgMessage struct {
	//MessageType uint8 //chat, friends, multi
	Who       uint64 `json:"who"` //使用who，表示客户端填充的接收者，服务器转发时会修改为发送者
	TimeStamp int64  `json:"timestamp"`
	Message   []byte `json:"message"`
}

type MsgMessageJson struct {
	//MessageType uint8 //chat, friends, multi
	Who       uint64 `json:"who,string"` //使用who，表示客户端填充的接收者，服务器转发时会修改为发送者
	TimeStamp int64  `json:"timestamp,string"`
	Nickname  string `json:"nickname"`
	Message   string `json:"message"`
}

type MsgMessageReceipt struct {
	ErrorCode uint16
}

//其它类型的单人消息，服务器收到以后，转发其它人时，都是使用1008的消息格式，但是消息id使用自己的。
const MsgId_AllFriendsMessage uint16 = 1009

type MsgAllFriendsMessage struct {
	Message []byte
}

const MsgId_GroupMessage uint16 = 1010

type MsgGroupMessage struct {
	Count     uint8
	GroupName []byte
	Message   []byte
}

const MsgId_MultiUsersMessage uint16 = 1011

type MsgMultiUsersMessage struct {
	Count   uint8
	To      []uint64
	Message []byte
}

// const MsgId_RoomMessage uint16 = 1012

// type MsgRoomMessage struct {
// 	Room    uint64
// 	From    uint64
// 	Message []byte
// }

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
	DataType  uint8
	Json      []byte
}

//create/delete user group
const MsgId_Group uint16 = 1015

type MsgReqGroupJson struct {
	Cmd     string `json:"cmd"`
	Name    string `json:"name"`
	OldName string `json:"oldname"`
	NewName string `json:"newname"`
}

type MsgRetGroupJson struct {
	// Cmd       string `json:"cmd"`
	// Name      string `json:"name,omitempty"`
	// OldName   string `json:"oldname,omitempty"`
	// NewName   string `json:"newname,omitempty"`
	ErrorCode uint16 `json:"errorcode"`
}

const MsgId_ReqGroupRefresh uint16 = 1016

type MsgReqGroupRefresh struct {
	GroupName string
}

type MsgRetGroupRefresh struct {
	ErrorCode uint16
	Json      []byte
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
const MsgId_ReqAddUserToGroup uint16 = 1019

type MsgReqAddUserToGroup struct {
	AppdataId uint64
	GroupName []byte
}

type MsgRetAddUserToGroup struct {
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

const MsgId_ReqModifyFriendComment uint16 = 1021

type MsgReqModifyFriendComment struct {
	Id      uint64
	Comment string
}

type MsgRetModifyFriendComment struct {
	ErrorCode uint16
}

const MsgId_KickOut uint16 = 1022

//search user/room
const MsgId_ReqIdSearch uint16 = 1023

type MsgReqIdSearch struct {
	Id uint64
}

type MsgRetIdSearch struct {
	ErrorCode uint16
	Json      []byte
}

const MsgId_ReqNicknameSearch uint16 = 1024

type MsgReqNicknameSearch struct {
	Nickname string
}

type MsgRetNicknameSearch struct {
	ErrorCode uint16
	Json      []byte
}

const MsgId_ReqRoomSearch uint16 = 1025

type MsgReqRoomSearch struct {
	Roomname string
}

type MsgRetRoomSearch struct {
	ErrorCode uint16
	Json      []byte
}

//history message ?

//create/delete room
const MsgId_ReqCreateRoom uint16 = 1100

type MsgReqCreateRoom struct {
	Name     string `json:"name"`
	RoomType byte   `json:"roomtype"`
	Password string `json:"password"`
	Jieshao  string `json:"jieshao"`
	Notice   string `json:"notice"` //公告
}

type MsgRetCreateRoom struct {
	ErrorCode uint16
}

const MsgId_ReqDeleteRoom uint16 = 1101

type MsgReqDeleteRoom struct {
	Rid uint64
}

type MsgRetDeleteRoom struct {
	ErrorCode uint16
}

//modify room setting
const MsgId_ReqUpdateRoomSetting uint16 = 1102

const (
	RoomSetting_None     byte = 0
	RoomSetting_RoomName byte = 0x1
	RoomSetting_RoomType byte = 0x2
	RoomSetting_Jieshao  byte = 0x4
	RoomSetting_Notice   byte = 0x8
	RoomSetting_Password byte = 0x10
)

const (
	RoomType_None byte = iota
	RoomType_Everyone
	RoomType_Apply
	RoomType_Password
	RoomType_Temp
)

type MsgReqUpdateRoomSetting struct {
	Rid      uint64 `json:"rid,string"`
	Bit      byte   `json:"bit"`
	RoomName string `json:"roomname"`
	RoomType byte   `json:"roomtype"` //1.everyone 2.need apply 3.password 4.temp
	Jieshao  string `json:"jieshao"`
	Notice   string `json:"notice"` //公告
	Password string `json:"password"`
}

type MsgRetUpdateRoomSetting struct {
	ErrorCode uint16
}

//join/quit room
// const MsgId_ReqJoinRoom uint16 = 1103

// type MsgReqJoinRoom struct {
// 	Rid uint64
// }

// type MsgRetJoinRoom struct {
// 	ErrorCode uint16
// }

// const MsgId_ReqQuitRoom uint16 = 1104

// type MsgReqQuitRoom struct {
// 	Rid uint64
// }

// type MsgRetQuitRoom struct {
// 	ErrorCode uint16
// }

const MsgId_RoomPresence uint16 = 1103

type MsgRoomPresence struct {
	PresenceType uint8  `json:"presencetype"` //available,subscribe,subscribed,unsubscribe,unsubscribed,unavailable,invisible
	Rid          uint64 `json:"rid,string"`
	Who          uint64 `json:"who,string"`
	Nickname     string `json:"nickname"`
	TimeStamp    int64  `json:"timestamp,string"`
	Password     string `json:"password"`
	Message      string `json:"message"`
}

type MsgRoomPresenceReceipt struct {
	ErrorCode uint16
}

//ban room user
const MsgId_ReqBanRoomUser uint16 = 1105

type MsgReqBanRoomUser struct {
	Rid       uint64
	AppdataId uint64
}

type MsgRetBanRoomUser struct {
	ErrorCode uint16
}

//jinyan/unjinyan room user
const MsgId_ReqJinyanRoomUser uint16 = 1106

type MsgReqJinyanRoomUser struct {
	Rid       uint64
	AppdataId uint64
}

type MsgRetJinyanRoomUser struct {
	ErrorCode uint16
}

const MsgId_ReqUnJinyanRoomUser uint16 = 1107

type MsgReqUnJinyanRoomUser struct {
	Rid       uint64
	AppdataId uint64
}

type MsgRetUnJinyanRoomUser struct {
	ErrorCode uint16
}

//add/remove room admin
const MsgId_ReqAddRoomAdmin uint16 = 1108

type MsgReqAddRoomAdmin struct {
	Rid       uint64
	AppdataId uint64
}

type MsgRetAddRoomAdmin struct {
	ErrorCode uint16
}

const MsgId_ReqRemoveRoomAdmin uint16 = 1109

type MsgReqRemoveRoomAdmin struct {
	Rid       uint64
	AppdataId uint64
}

type MsgRetRemoveRoomAdmin struct {
	ErrorCode uint16
}

//room message
const MsgId_RoomMessage uint16 = 1110

type MsgRoomMessage struct {
	Rid       uint64 `json:"rid,string"`
	Who       uint64 `json:"who,string"` //使用who，表示客户端填充的接收者，服务器转发时会修改为发送者
	TimeStamp int64  `json:"timestamp,string"`
	Nickname  string `json:"nickname"`
	Message   string `json:"message"`
}

type MsgRoomReceipt struct {
	ErrorCode uint16
}

//create/delete room group
//invite user
//message broadcast

//公共频道消息传输，如世界、国家等
const MsgId_PublicChannelMessage uint16 = 1200

type MsgPublicChannelMessage struct {
	Who       uint64 `json:"who,string"` //使用who，表示客户端填充的接收者，服务器转发时会修改为发送者
	TimeStamp int64  `json:"timestamp,string"`
	Nickname  string `json:"nickname"`
	Message   string `json:"message"`
	Channel   string `json:"channel"`
}

type MsgPublicChannelMessageReceipt struct {
	ErrorCode uint16
}

//define RPC
