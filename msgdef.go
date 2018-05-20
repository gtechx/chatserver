package main

import (
	"fmt"
)

//const define need uppercase for first word or all uppercase with "_" connected
const (
	DataFrame byte = iota
	//JsonFrame
	//BinaryFrame
	PingFrame
	PongFrame
	CloseFrame
	AckFrame
	ErrorFrame
	EchoFrame
)

var msgHandler = map[uint16]func(uint64, []byte){}

func registerMsgHandler(msgid uint16, handler func(uint64, []byte)) {
	_, ok := msgHandler[msgid]

	if ok {
		fmt.Println("Error: dumplicate msgid ", msgid)
		return
	}
	msgHandler[msgid] = handler
}
