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

var msgHandler = map[uint16]func(ISession, []byte) (uint16, interface{}){}

func registerMsgHandler(msgid uint16, handler func(ISession, []byte) (uint16, interface{})) {
	_, ok := msgHandler[msgid]

	if ok {
		fmt.Println("Error: dumplicate msgid ", msgid)
		return
	}
	msgHandler[msgid] = handler
}

func HandleMsg(msgid uint16, ISession, buff []byte) (uint16, interface{}) {
	handler, ok := msgHandler[msgid]

	if ok {
		return handler(buff)
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
	Flag      bool
	ErrorCode uint16
	Token     []byte
}

const MsgId_ReqEnterChat uint16 = 1001

type MsgReqEnterChat struct {
	AppdataId uint64
}

type MsgRetEnterChat struct {
	Flag      bool
	ErrorCode uint16
}

const MsgId_ReqQuitChat uint16 = 1002

type MsgReqQuitChat struct {
	AppdataId uint64
}

type MsgRetQuitChat struct {
	Flag      bool
	ErrorCode uint16
}

const MsgId_ReqAppDataIdList uint16 = 1003

type MsgReqAppDataIdList struct {
}

type MsgRetAppDataIdList struct {
	ErrorCode     uint16
	AppDataIdList []uint64
}
