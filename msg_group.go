package main

import (
	"encoding/json"
	"fmt"

	. "github.com/gtechx/base/common"
	"github.com/gtechx/chatserver/db"
)

func init() {
	RegisterGroupMsg()
}

func RegisterGroupMsg() {
	registerMsgHandler(MsgId_Group, HandlerGroup)
	registerMsgHandler(MsgId_ReqGroupRefresh, HandlerGroupRefresh)
	registerMsgHandler(MsgId_ReqAddUserToGroup, HandlerReqAddUserToGroup)
}

func HandlerReqAddUserToGroup(sess ISession, data []byte) (uint16, interface{}) {
	appdataid := Uint64(data)
	groupname := String(data[8:])
	errcode := ERR_NONE
	dbMgr := gtdb.Manager()

	flag, err := dbMgr.IsAppDataExists(appdataid)
	if err != nil {
		errcode = ERR_DB
	} else {
		if !flag {
			errcode = ERR_APPDATAID_NOT_EXISTS
		} else {
			flag, err = dbMgr.IsGroupExists(sess.ID(), groupname)
			if err != nil {
				errcode = ERR_DB
			} else {
				if !flag {
					errcode = ERR_GROUP_NOT_EXISTS
				} else {
					err = dbMgr.MoveToGroup(sess.ID(), appdataid, groupname)
					if err != nil {
						errcode = ERR_DB
					}
				}
			}
		}
	}

	return errcode, errcode
}

func HandlerGroupRefresh(sess ISession, data []byte) (uint16, interface{}) {
	groupname := String(data)
	errcode := ERR_NONE
	dbMgr := gtdb.Manager()
	ret := &MsgRetGroupRefresh{}

	flag, err := dbMgr.IsGroupExists(sess.ID(), groupname)
	if err != nil {
		errcode = ERR_DB
	} else {
		if !flag {
			errcode = ERR_GROUP_NOT_EXISTS
		} else {
			friendlist := map[string][]*gtdb.FriendJson{}
			list, err := dbMgr.GetFriendInfoList(sess.ID(), groupname)
			if err != nil {
				errcode = ERR_DB
			} else {
				friendlist[groupname] = list
			}
			ret.Json, err = json.Marshal(friendlist)
			if err != nil {
				errcode = ERR_JSON_SERIALIZE
				ret.Json = nil
			}
		}
	}
	ret.ErrorCode = errcode

	return errcode, ret
}

func HandlerGroup(sess ISession, data []byte) (uint16, interface{}) {
	var groupmsg *MsgReqGroupJson = &MsgReqGroupJson{}
	err := json.Unmarshal(data, groupmsg)

	fmt.Println(string(data))
	fmt.Println(groupmsg)
	if err != nil {
		fmt.Println(err.Error())
		return ERR_INVALID_JSON, ERR_INVALID_JSON
	}

	errcode := ERR_NONE
	dbMgr := gtdb.Manager()

	switch groupmsg.Cmd {
	case "create":
		flag, err := dbMgr.IsGroupExists(sess.ID(), groupmsg.Name)
		if err != nil {
			errcode = ERR_DB
		} else {
			if flag {
				errcode = ERR_GROUP_NOT_EXISTS
			} else {
				tbl_group := &gtdb.Group{Groupname: groupmsg.Name, Dataid: sess.ID()}
				err = dbMgr.AddGroup(tbl_group)
				if err != nil {
					errcode = ERR_DB
				}
			}
		}
	case "delete":
		flag, err := dbMgr.IsGroupExists(sess.ID(), groupmsg.Name)
		if err != nil {
			errcode = ERR_DB
		} else {
			if !flag {
				errcode = ERR_GROUP_NOT_EXISTS
			} else {
				//check if group has friend
				count, err := dbMgr.GetFriendCountInGroup(sess.ID(), groupmsg.Name)

				if err != nil {
					errcode = ERR_DB
				} else {
					if count > 0 {
						errcode = ERR_GROUP_NOT_EMPTY
					} else {
						err = dbMgr.RemoveGroup(sess.ID(), groupmsg.Name)
						if err != nil {
							errcode = ERR_DB
						}
					}
				}
			}
		}
	case "rename":
		flag, err := dbMgr.IsGroupExists(sess.ID(), groupmsg.OldName)
		if err != nil {
			errcode = ERR_DB
		} else {
			if !flag {
				errcode = ERR_OLD_GROUP_NOT_EXISTS
			} else {
				flag, err := dbMgr.IsGroupExists(sess.ID(), groupmsg.NewName)
				if err != nil {
					errcode = ERR_DB
				} else {
					if flag {
						errcode = ERR_NEW_GROUP_EXISTS
					} else {
						err := dbMgr.RenameGroup(sess.ID(), groupmsg.OldName, groupmsg.NewName)
						if err != nil {
							errcode = ERR_DB
						}
					}
				}
			}
		}
	}

	return errcode, errcode
}
