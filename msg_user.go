package main

import (
	"errors"

	. "github.com/gtechx/base/common"
	"github.com/gtechx/chatserver/db"
	"github.com/satori/go.uuid"
)

func RegisterUserMsg() {
	registerMsgHandler(MsgId_ReqLogin, HandlerReqLogin)
	registerMsgHandler(MsgId_EnterChat, HandlerEnterChat)
}

func HandlerReqLogin(buff []byte) (interface{}, error) {
	slen := int(buff[0])
	account := String(buff[1 : 1+slen])
	buff = buff[1+slen:]
	slen = int(buff[0])
	password := String(buff[1 : 1+slen])

	dbmgr := gtdb.Manager()

	ok, err := dbmgr.IsAccountExists(account)

	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, errors.New("account not exists")
	}

	tbl_account, err := dbmgr.GetAccount(account)

	if err != nil {
		return nil, err
	}

	md5password := GetSaltedPassword(password, tbl_account.Salt)
	if md5password != tbl_account.Password {
		return nil, errors.New("password wrong")
	}

	token, err := uuid.NewV4()

	if err != nil {
		return nil, err
	}

	return token.Bytes(), nil
}

func HandlerEnterChat(buff []byte) (interface{}, error) {
	appdataid := Uint64(buff)
	dbmgr := gtdb.Manager()

	ok, err := dbmgr.IsAppDataExists(appdataid)

	if err != nil {
		return nil, err
	}

	if !ok {
		return false, errors.New("id not exists")
	}

	return appdataid, nil
}
