package main

import (
	"encoding/json"
	"time"

	. "github.com/gtechx/base/common"
	"github.com/gtechx/chatserver/config"
	"github.com/gtechx/chatserver/db"
)

func RegisterUserMsg() {
	registerMsgHandler(MsgId_ReqUserData, HandlerReqUserData)
	registerMsgHandler(MsgId_Presence, HandlerPresence)
	//registerMsgHandler(MsgId_EnterChat, HandlerEnterChat)
}

func HandlerReqUserData(sess ISession, data []byte) (uint16, interface{}) {
	appdata, err := gtdb.Manager().GetAppData(sess.ID())
	errcode := ERR_NONE
	var jsonbytes []byte

	if err != nil {
		errcode = ERR_DB
	} else {
		jsonbytes, err = json.Marshal(appdata)
		if err != nil {
			errcode = ERR_UNKNOWN
		}
	}

	ret := &MsgRetUserData{errcode, jsonbytes}
	return errcode, ret
}

func SendMessageToUserOnline(to uint64, data []byte) uint16 {
	dbMgr := gtdb.Manager()
	online, err := dbMgr.GetUserOnlineInfo(to)
	if err != nil {
		return ERR_DB
	}

	err = gtdb.Manager().SendMsgToUserOnline(append(Bytes(to), data...), online.Serveraddr)
	if err != nil {
		return ERR_DB
	}
	return ERR_NONE
}

func SendMessageToUserOffline(to uint64, data []byte) uint16 {
	err := gtdb.Manager().SendMsgToUserOffline(to, data)
	if err != nil {
		return ERR_DB
	}
	return ERR_NONE
}

func SendMessageToUser(to uint64, data []byte) uint16 {
	dbMgr := gtdb.Manager()
	flag, err := dbMgr.IsUserOnline(to)
	if err != nil {
		return ERR_DB
	}

	if flag {
		return SendMessageToUserOnline(to, data)
	} else {
		return SendMessageToUserOffline(to, data)
	}
}

func SendMessageToFriendsOnline(id uint64, data []byte) uint16 {
	dbMgr := gtdb.Manager()
	friendinfolist, err := dbMgr.GetOnlineFriendInfoList(id)
	if err != nil {
		return ERR_DB
	}

	for _, online := range friendinfolist {
		err = dbMgr.SendMsgToUserOnline(append(Bytes(online.Dataid), data...), online.Serveraddr)
		if err != nil {
			return ERR_DB
		}
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
	message := data[17:]

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
						tbl_from := &gtdb.Friend{Dataid: who, Otherdataid: sess.ID(), Group: config.DefaultGroupName}
						tbl_to := &gtdb.Friend{Dataid: sess.ID(), Otherdataid: who, Group: config.DefaultGroupName}
						err = dbMgr.AddFriend(tbl_from, tbl_to)

						if err != nil {
							errcode = ERR_DB
						} else {
							presencebytes := Bytes(presence)
							errcode = SendMessageToUser(who, presencebytes)
							dbMgr.RemovePresence(sess.ID(), who)
						}
					}
				}
			case PresenceType_Unsubscribe:
				//check if the two are friend, if not omit thie message, else delete friend and send to who.
				flag, err = dbMgr.IsFriend(sess.ID(), who)
				if err != nil {
					errcode = ERR_DB
				} else {
					if !flag {
						errcode = ERR_FRIEND_NOT_EXISTS
					} else {
						err = dbMgr.RemoveFriend(sess.ID(), who)
						if err != nil {
							errcode = ERR_DB
						} else {
							presencebytes := Bytes(presence)
							errcode = SendMessageToUser(who, presencebytes)
						}
					}
				}
			case PresenceType_Unsubscribed:
				//check if server has record, if not omit this message, else send to record sender
				flag, err = dbMgr.IsPresenceExists(sess.ID(), who)
				if err != nil {
					errcode = ERR_DB
				} else {
					if !flag {
						errcode = ERR_UNKNOWN
					} else {
						presencebytes := Bytes(presence)
						errcode = SendMessageToUser(who, presencebytes)
					}
				}
			case PresenceType_Available, PresenceType_Unavailable, PresenceType_Invisible:
				//send to my friend online
				presencebytes := Bytes(presence)
				SendMessageToFriendsOnline(sess.ID(), presencebytes)
			}
		}
	}

	//ret := &MsgRetUserData{errcode, jsonbytes}
	return errcode, errcode
}

// type MsgReqDataList struct {
// 	DataType uint8 //friend, presence,room, black, message, roommessage
// }

// type MsgRetDataList struct {
// 	ErrorCode uint16
// 	Json      []byte
// }
func HandlerReqDataList(sess ISession, data []byte) (uint16, interface{}) {
	datatype := uint8(data[0])

	errcode := ERR_NONE
	dbMgr := gtdb.Manager()
	flag, err := dbMgr.IsAppDataExists(who)
	switch datatype {
	case DataType_Group:
	case DataType_Friend:
	case DataType_Presence:
	case DataType_Black:
	case DataType_Message:
	case DataType_Room:
	case DataType_RoomMessage:
	}
}
