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
