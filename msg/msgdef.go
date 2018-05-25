package gtmsg

import (
	"errors"
	"fmt"
)

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

var msgHandler = map[uint16]func([]byte) (interface{}, error){}

func registerMsgHandler(msgid uint16, handler func([]byte) (interface{}, error)) {
	_, ok := msgHandler[msgid]

	if ok {
		fmt.Println("Error: dumplicate msgid ", msgid)
		return
	}
	msgHandler[msgid] = handler
}

func Handle(msgid uint16, buff []byte) (interface{}, error) {
	handler, ok := msgHandler[msgid]

	if ok {
		return handler(buff)
	}

	return nil, errors.New("msgid handler not exists")
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

const MsgId_EnterChat uint16 = 1001

type MsgEnterChat struct {
	AppdataId uint64
}
