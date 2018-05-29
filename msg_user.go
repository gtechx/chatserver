package main

import (
	. "github.com/gtechx/base/common"
	"github.com/gtechx/chatserver/config"
	"github.com/gtechx/chatserver/db"
)

func RegisterUserMsg() {
	//registerMsgHandler(MsgId_ReqLogin, HandlerReqLogin)
	registerMsgHandler(MsgId_EnterChat, HandlerEnterChat)
}

func HandlerReqLogin(buff []byte) (uint16, interface{}) {
	slen := int(buff[0])
	account := String(buff[1 : 1+slen])
	buff = buff[1+slen:]
	slen = int(buff[0])
	password := String(buff[1 : 1+slen])

	dbmgr := gtdb.Manager()
	errcode := ERR_NONE
	var tokenbytes []byte

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
				} else {
					token, err := uuid.NewV4()

					if err != nil {
						errcode = ERR_UNKNOWN
					} else {
						tokenbytes = token.Bytes()
					}
				}
			}
		}
	}

	ret := &MsgRetLogin{errcode == ERR_NONE, errcode, tokenbytes}
	return errcode, ret
	//sess.Send(ret)
}

func HandlerReqAppDataIdList(sess ISession, buff []byte) (uint16, interface{}) {
	idlist, err := gtdb.Manager().GetAppDataIdList(sess.AppName(), sess.ZoneName(), sess.Account())
	errcode := ERR_NONE
	if err != nil {
		errcode = ERR_DB
	}
	ret := &MsgRetAppDataIdList{errcode, idlist}
	//sess.Send(ret)
	return errcode, ret
}

func HandlerEnterChat(sess ISession, buff []byte) (uint16, interface{}) {
	appdataid := Uint64(buff)
	dbmgr := gtdb.Manager()
	errcode := ERR_NONE

	ok, err := dbmgr.IsAppDataExists(appdataid)

	if err != nil {
		errcode = ERR_DB
	} else {
		if !ok {
			errcode = ERR_APPDATAID_NOT_EXISTS
		} else {
			tbl_online := &gtdb.Online{appdataid, config.ServerAddr, "available"}
			err = dbmgr.SetUserOnline(tbl_online)
			if err != nil {
				errcode = errcode = ERR_DB
			}
		}
	}

	ret := &MsgRetEnterChat{errcode == ERR_NONE, errcode}
	//sess.Send(ret)
	return errcode, ret
}

func HandlerQuitChat(sess ISession, buff []byte) (uint16, interface{}) {
	appdataid := Uint64(buff)
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
				errcode = errcode = ERR_DB
			}
		}
	}

	ret := &MsgRetQuitChat{errcode == ERR_NONE, errcode}
	//sess.Send(ret)
	return errcode, ret
}
