package gtmsg

import (
	"errors"

	. "github.com/gtechx/base/common"
)

func RegisterUserMsg() {
	registerMsgHandler(MsgId_ReqLogin, HandlerReqLogin)
	registerMsgHandler(MsgId_EnterChat, HandlerEnterChat)
}

func HandlerReqLogin(buff []byte) (interface{}, error) {
	slen := int(buff[0])
	account := String(buff[1 : 1+alen])
	buff = buff[1+alen:]
	slen = int(buff[0])
	password := String(buff[1 : 1+alen])

	dbmgr := gtdb.Manager()

	if !dbmgr.IsAccountExists(account) {
		return "", errors.New("account not exists")
	}

	tbl_account, err := dbmgr.GetAccount(account)

	if err != nil {
		return "", err
	}

	md5password := GetSaltedPassword(password, tbl_account.Salt)
	if md5password != tbl_account.Password {
		return "", errors.New("password wrong")
	}

	return account, nil
}

func HandlerEnterChat(buff []byte) (interface{}, error) {
	appdataid := Uint64(buff)
	dbmgr := gtdb.Manager()

	if !dbmgr.IsAppDataExists(appdataid) {
		return false, errors.New("id not exists")
	}

	return appdataid, nil
}
