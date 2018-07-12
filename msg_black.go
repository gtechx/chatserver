package main

import (
	. "github.com/gtechx/base/common"
	"github.com/gtechx/chatserver/db"
)

func init() {
	RegisterBlackMsg()
}

func RegisterBlackMsg() {
	registerMsgHandler(MsgId_ReqAddBlack, HandlerReqAddBlack)
	registerMsgHandler(MsgId_ReqRemoveBlack, HandlerReqRemoveBlack)
}

func isInBlack(id, otherid uint64) (bool, error) {
	flag, err := gtdb.Manager().IsInBlack(id, otherid)
	if err != nil {
		return false, err
	}
	return flag, nil
}

func HandlerReqAddBlack(sess ISession, data []byte) (uint16, interface{}) {
	appdataid := Uint64(data)
	errcode := ERR_NONE
	dbMgr := gtdb.Manager()

	flag, err := dbMgr.IsAppDataExists(appdataid)
	if err != nil {
		errcode = ERR_DB
	} else {
		if !flag {
			errcode = ERR_APPDATAID_NOT_EXISTS
		} else {
			flag, err = dbMgr.IsInBlack(sess.ID(), appdataid)
			if err != nil {
				errcode = ERR_DB
			} else {
				if flag {
					errcode = ERR_IN_BLACKLIST
				} else {
					tbl_black := &gtdb.Black{Dataid: sess.ID(), Otherdataid: appdataid}
					err = dbMgr.AddBlack(tbl_black)
					if err != nil {
						errcode = ERR_DB
					}
				}
			}
		}
	}

	return errcode, errcode
}

func HandlerReqRemoveBlack(sess ISession, data []byte) (uint16, interface{}) {
	appdataid := Uint64(data)
	errcode := ERR_NONE
	dbMgr := gtdb.Manager()

	flag, err := dbMgr.IsInBlack(sess.ID(), appdataid)
	if err != nil {
		errcode = ERR_DB
	} else {
		if !flag {
			errcode = ERR_NOT_IN_BLACKLIST
		} else {
			err = dbMgr.RemoveFromBlack(sess.ID(), appdataid)
			if err != nil {
				errcode = ERR_DB
			}
		}
	}

	return errcode, errcode
}
