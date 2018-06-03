package main

import (
	"fmt"

	. "github.com/gtechx/base/common"
	"github.com/gtechx/chatserver/config"
	"github.com/gtechx/chatserver/db"
	"github.com/satori/go.uuid"
)

// func RegisterUserMsg() {
// 	//registerMsgHandler(MsgId_ReqLogin, HandlerReqLogin)
// 	//registerMsgHandler(MsgId_EnterChat, HandlerEnterChat)
// }

func checkAccount(account, password string) uint16 {
	dbmgr := gtdb.Manager()
	errcode := ERR_NONE

	ok, err := dbmgr.IsAccountExists(account)

	if err != nil {
		errcode = ERR_DB
	} else {
		if !ok {
			errcode = ERR_ACCOUNT_NOT_EXISTS
		} else {
			tbl_account, err := dbmgr.GetAccount(account)

			if err != nil {
				errcode = ERR_DB
			} else {
				md5password := GetSaltedPassword(password, tbl_account.Salt)
				if md5password != tbl_account.Password {
					errcode = ERR_PASSWORD_INVALID
				}
			}
		}
	}

	return errcode
}

func HandlerReqLogin(buff []byte) (uint16, interface{}) {
	slen := int(buff[0])
	account := String(buff[1 : 1+slen])
	buff = buff[1+slen:]
	slen = int(buff[0])
	password := String(buff[1 : 1+slen])

	var tokenbytes []byte
	//dbmgr := gtdb.Manager()
	errcode := checkAccount(account, password)

	if errcode == ERR_NONE {
		token, err := uuid.NewV4()

		if err != nil {
			errcode = ERR_UNKNOWN
		} else {
			tokenbytes = token.Bytes()
		}
	}

	fmt.Println("tokenbytes len:", len(tokenbytes))
	ret := &MsgRetLogin{errcode, tokenbytes}
	return errcode, ret
	//sess.Send(ret)
}

func HandlerReqChatLogin(account, password, appname, zonename string) (uint16, interface{}) {
	errcode := checkAccount(account, password)
	if errcode == ERR_NONE {
		idlist, err := gtdb.Manager().GetAppDataIdList(appname, zonename, account)
		if err != nil {
			errcode = ERR_DB
		}
		ret := &MsgRetChatLogin{errcode, idlist}
		return errcode, ret
	}
	ret := &MsgRetChatLogin{ErrorCode: errcode}
	return errcode, ret
}

func HandlerReqCreateAppdata(appname, zonename, account, nickname, regip string) (uint16, interface{}) {
	dbMgr := gtdb.Manager()
	errcode := ERR_NONE
	id := uint64(0)

	flag, err := dbMgr.IsNicknameExists(appname, zonename, account, nickname)
	if err != nil {
		errcode = ERR_DB
	} else if flag {
		errcode = ERR_NICKNAME_EXISTS
	} else {
		tbl_appdata := &gtdb.AppData{Appname: appname, Zonename: zonename, Account: account, Nickname: nickname, Regip: regip}
		err = dbMgr.CreateAppData(tbl_appdata)

		if err != nil {
			errcode = ERR_DB
		} else {
			id = tbl_appdata.ID
		}
	}

	ret := &MsgRetCreateAppdata{errcode, id}
	return errcode, ret
}

// func HandlerReqAppDataIdList(appname, zonename, account string) (uint16, interface{}) {
// 	idlist, err := gtdb.Manager().GetAppDataIdList(appname, zonename, account)
// 	errcode := ERR_NONE
// 	if err != nil {
// 		errcode = ERR_DB
// 	}
// 	ret := &MsgRetAppDataIdList{errcode, idlist}
// 	//sess.Send(ret)
// 	return errcode, ret
// }

func HandlerReqEnterChat(appdataid uint64) (uint16, interface{}) {
	dbmgr := gtdb.Manager()
	errcode := ERR_NONE

	ok, err := dbmgr.IsAppDataExists(appdataid)

	if err != nil {
		errcode = ERR_DB
	} else {
		if !ok {
			errcode = ERR_APPDATAID_NOT_EXISTS
		} else {
			tbl_online := &gtdb.Online{Dataid: appdataid, Serveraddr: config.ServerAddr, State: "available"}
			err = dbmgr.SetUserOnline(tbl_online)
			if err != nil {
				errcode = ERR_DB
			}
		}
	}

	ret := &MsgRetEnterChat{errcode}
	return errcode, ret
}

func HandlerReqQuitChat(appdataid uint64) (uint16, interface{}) {
	dbmgr := gtdb.Manager()
	errcode := ERR_NONE

	ok, err := dbmgr.IsAppDataExists(appdataid)

	if err != nil {
		errcode = ERR_DB
	} else {
		if !ok {
			errcode = ERR_APPDATAID_NOT_EXISTS
		} else {
			err = dbmgr.SetUserOffline(appdataid)
			if err != nil {
				errcode = ERR_DB
			}
		}
	}

	ret := &MsgRetQuitChat{errcode}
	//sess.Send(ret)
	return errcode, ret
}
