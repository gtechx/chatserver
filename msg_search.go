package main

import (
	"encoding/json"

	. "github.com/gtechx/base/common"
	"github.com/gtechx/chatserver/db"
)

func init() {
	RegisterSearchMsg()
}

func RegisterSearchMsg() {
	registerMsgHandler(MsgId_ReqIdSearch, HandlerReqIdSearch)
	registerMsgHandler(MsgId_ReqNicknameSearch, HandlerReqNicknameSearch)
	registerMsgHandler(MsgId_ReqRoomSearch, HandlerReqRoomSearch)
}

func HandlerReqIdSearch(sess ISession, data []byte) (uint16, interface{}) {
	appdataid := Uint64(data)
	errcode := ERR_NONE
	dbMgr := gtdb.Manager()
	ret := &MsgRetIdSearch{}

	searchret, err := dbMgr.SearchUserById(appdataid)
	if err != nil {
		errcode = ERR_DB
	} else {
		if searchret != nil {
			ret.Json, err = json.Marshal(searchret)
			if err != nil {
				errcode = ERR_JSON_SERIALIZE
				ret.Json = nil
			}
		}
	}
	ret.ErrorCode = errcode

	return errcode, ret
}

func HandlerReqNicknameSearch(sess ISession, data []byte) (uint16, interface{}) {
	nickname := String(data)
	errcode := ERR_NONE
	dbMgr := gtdb.Manager()
	ret := &MsgRetNicknameSearch{}

	searchret, err := dbMgr.SearchUserByNickname(nickname)
	if err != nil {
		errcode = ERR_DB
	} else {
		if searchret != nil {
			ret.Json, err = json.Marshal(searchret)
			if err != nil {
				errcode = ERR_JSON_SERIALIZE
				ret.Json = nil
			}
		}
	}
	ret.ErrorCode = errcode

	return errcode, ret
}

func HandlerReqRoomSearch(sess ISession, data []byte) (uint16, interface{}) {
	roomname := String(data)
	errcode := ERR_NONE
	dbMgr := gtdb.Manager()
	ret := &MsgRetRoomSearch{}

	searchret, err := dbMgr.SearchRoom(roomname)
	if err != nil {
		errcode = ERR_DB
	} else {
		if searchret != nil {
			ret.Json, err = json.Marshal(searchret)
			if err != nil {
				errcode = ERR_JSON_SERIALIZE
				ret.Json = nil
			}
		}
	}
	ret.ErrorCode = errcode

	return errcode, ret
}
