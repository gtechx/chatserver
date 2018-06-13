package main

import (
	"encoding/json"

	. "github.com/gtechx/base/common"
	"github.com/gtechx/chatserver/db"
)

func RegisterUserMsg() {
	registerMsgHandler(MsgId_ReqUserData, HandlerReqUserData)
	//registerMsgHandler(MsgId_EnterChat, HandlerEnterChat)
}

func HandlerReqUserData(sess ISession, data []byte) (uint16, interface{}) {
	appdata, err := gtdb.Manager().GetAppData(sess.ID())
	errcode := ERR_NONE
	var jsonbytes []byte

	if err != nil {
		errcode = ERR_DB
	} else {
		jsonbytes, err := json.Marshal(pageapp)
		if err != nil {
			errcode = ERR_UNKNOWN
		}
	}

	ret := &MsgRetUserData{errcode, jsonbytes}
	return errcode, ret
}

// type MsgPresence struct {
// 	PresenceType uint8 //available,subscribe,subscribed,unsubscribe,unsubscribed,unavailable,invisible
// 	Who          uint64
// 	Message      string
// }
func HandlerPresence(sess ISession, data []byte) (uint16, interface{}) {
	presencetype := uint8(data[0])
	who := Uint64(data[1:])
	message := String(data[9:])

	errcode := ERR_NONE
	flag, err := gtdb.Manager().IsAppDataExists(who)

	if err != nil {
		errcode = ERR_DB
	} else {
		if !flag {
			errcode = ERR_APPDATAID_NOT_EXISTS
		} else {
			//
		}
	}

	ret := &MsgRetUserData{errcode, jsonbytes}
	return errcode, ret
}
