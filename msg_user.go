package main

import (
	"encoding/json"
	"time"

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
		jsonbytes, err = json.Marshal(pageapp)
		if err != nil {
			errcode = ERR_UNKNOWN
		}
	}

	ret := &MsgRetUserData{errcode, jsonbytes}
	return errcode, ret
}

func SendMessageToUser(to uint64, data []byte) uint16 {
	dbMgr := gtdb.Manager()
	flag, err := dbMgr.IsUserOnline(to)
	if err != nil {
		return ERR_DB
	}

	if flag {
		err = dbMgr.SendMsgToUserOnline(to, append(Bytes(to), data...))
	} else {
		err = dbMgr.SendMsgToUserOffline(to, data)
	}

	if err != nil {
		return ERR_DB
	}

	return ERR_NONE
}

// type MsgPresence struct {
// 	PresenceType uint8 //available,subscribe,subscribed,unsubscribe,unsubscribed,unavailable,invisible
// 	Who          uint64
// 	Message      string
// }
func HandlerPresence(sess ISession, data []byte) (uint16, interface{}) {
	presencetype := uint8(data[0])
	who := Uint64(data[1:])
	timestamp := Int64(data[9:])
	message := String(data[17:])

	timestamp = time.Now().Unix()

	presence := &MsgPresence{PresenceType: presencetype, Who: sess.ID(), TimeStamp: timestamp, Message: message}

	errcode := ERR_NONE
	dbMgr := gtdb.Manager()
	flag, err := dbMgr.IsAppDataExists(who)

	if err != nil {
		errcode = ERR_DB
	} else {
		if !flag {
			errcode = ERR_APPDATAID_NOT_EXISTS
		} else {
			//
			switch presencetype {
			case PresenceType_Subscribe:
				flag, err = dbMgr.IsFriend(sess.ID(), who)
				if err != nil {
					errcode = ERR_DB
				} else {
					if flag {
						errcode = ERR_FRIEND_EXISTS
					} else {
						//send presence to who and record this presence for who's answer
						presencebytes := Bytes(presence)
						err = dbMgr.AddPresence(sess.ID(), who, presencebytes)
						if err != nil {
							errcode = ERR_DB
						} else {
							//send to who
							errcode = SendMessageToUser(who, presencebytes)
							// if errcode != ERR_NONE {

							// }
						}
					}
				}
			case PresenceType_Subscribed:
				//check if server has record, if not omit this message, else send to record sender
				flag, err = dbMgr.IsPresenceExists(sess.ID(), who)
				if err != nil {
					errcode = ERR_DB
				} else {
					if !flag {
						errcode = ERR_UNKNOWN
					} else {
						dbMgr.AddFriend()
						presencebytes := Bytes(presence)
						errcode = SendMessageToUser(who, presencebytes)
						dbMgr.RemovePresence(sess.ID(), who)
					}
				}
			case PresenceType_Unsubscribe:
				//check if the two are friend, if not omit thie message, else delete friend and send to who.
			case PresenceType_Unsubscribed:
				//check if server has record, if not omit this message, else send to record sender
			}
		}
	}

	//ret := &MsgRetUserData{errcode, jsonbytes}
	return errcode, errcode
}
