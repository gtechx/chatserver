package main

import (
	"encoding/json"
	"fmt"
	"time"

	. "github.com/gtechx/base/common"
	"github.com/gtechx/chatserver/config"
	"github.com/gtechx/chatserver/db"
)

func init() {
	RegisterUserMsg()
}

func RegisterUserMsg() {
	registerMsgHandler(MsgId_ReqUserData, HandlerReqUserData)
	registerMsgHandler(MsgId_Presence, HandlerPresence)
	registerMsgHandler(MsgId_Message, HandlerMessage)
	registerMsgHandler(MsgId_ReqDataList, HandlerReqDataList)
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
	onlinelist, err := dbMgr.GetUserOnlineInfo(to)
	if err != nil {
		return ERR_DB
	}

	for _, online := range onlinelist {
		err = gtdb.Manager().SendMsgToUserOnline(append(Bytes(to), data...), online.Serveraddr)
		if err != nil {
			return ERR_DB
		}
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
	friendinfolist, err := dbMgr.GetFriendOnlineList(id)
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
	var presence *MsgPresence = &MsgPresence{}
	err := json.Unmarshal(data, presence)

	fmt.Println(string(data))
	fmt.Println(presence)
	if err != nil {
		fmt.Println(err.Error())
		return ERR_INVALID_JSON, ERR_INVALID_JSON
	}

	appdata, err := gtdb.Manager().GetAppData(sess.ID())

	if err != nil {
		return ERR_DB, ERR_DB
	}

	presence.Nickname = appdata.Nickname

	presencetype := presence.PresenceType
	who := presence.Who
	//timestamp := Int64(data[9:])
	//message := data[17:]

	if who == sess.ID() {
		return ERR_FRIEND_SELF, ERR_FRIEND_SELF
	}

	presence.TimeStamp = time.Now().Unix()

	//presence := &MsgPresence{PresenceType: presencetype, Who: sess.ID(), TimeStamp: timestamp, Message: message}

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
						presencebytes, err := json.Marshal(presence)
						if err != nil {
							errcode = ERR_INVALID_JSON
						} else {
							senddata := packageMsg(RetFrame, 0, MsgId_Presence, presencebytes)
							err = dbMgr.AddPresence(sess.ID(), who, presencebytes)
							if err != nil {
								errcode = ERR_DB
							} else {
								//send to who
								errcode = SendMessageToUser(who, senddata)
								// if errcode != ERR_NONE {

								// }
							}
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
						errcode = ERR_PRESENCE_NOT_EXISTS
					} else {
						tbl_from := &gtdb.Friend{Dataid: who, Otherdataid: sess.ID(), Group: config.DefaultGroupName}
						tbl_to := &gtdb.Friend{Dataid: sess.ID(), Otherdataid: who, Group: config.DefaultGroupName}
						err = dbMgr.AddFriend(tbl_from, tbl_to)

						if err != nil {
							errcode = ERR_DB
						} else {
							presencebytes, err := json.Marshal(presence)
							if err != nil {
								errcode = ERR_INVALID_JSON
							} else {
								senddata := packageMsg(RetFrame, 0, MsgId_Presence, presencebytes)
								errcode = SendMessageToUser(who, senddata)
								dbMgr.RemovePresence(sess.ID(), who)
							}
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
							presencebytes, err := json.Marshal(presence)
							if err != nil {
								errcode = ERR_INVALID_JSON
							} else {
								senddata := packageMsg(RetFrame, 0, MsgId_Presence, presencebytes)
								errcode = SendMessageToUser(who, senddata)
							}
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
						errcode = ERR_PRESENCE_NOT_EXISTS
					} else {
						presencebytes, err := json.Marshal(presence)
						if err != nil {
							errcode = ERR_INVALID_JSON
						} else {
							senddata := packageMsg(RetFrame, 0, MsgId_Presence, presencebytes)
							errcode = SendMessageToUser(who, senddata)
							dbMgr.RemovePresence(sess.ID(), who)
						}
					}
				}
			case PresenceType_Available, PresenceType_Unavailable, PresenceType_Invisible:
				//send to my friend online
				presencebytes, err := json.Marshal(presence)
				if err != nil {
					errcode = ERR_INVALID_JSON
				} else {
					senddata := packageMsg(RetFrame, 0, MsgId_Presence, presencebytes)
					SendMessageToFriendsOnline(sess.ID(), senddata)
				}
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
// 	DataType  uint8
// 	Json      []byte
// }
func HandlerReqDataList(sess ISession, data []byte) (uint16, interface{}) {
	datatype := uint8(data[0])

	errcode := ERR_NONE
	dbMgr := gtdb.Manager()

	ret := &MsgRetDataList{}
	ret.DataType = datatype

	switch datatype {
	case DataType_Friend:
		list, err := dbMgr.GetFriendInfoList(sess.ID())
		if err != nil {
			errcode = ERR_DB
		} else {
			ret.Json, err = json.Marshal(list)
			if err != nil {
				errcode = ERR_UNKNOWN
				ret.Json = nil
			}
		}
	case DataType_Presence:
		list, err := dbMgr.GetAllPresence(sess.ID())
		if err != nil {
			errcode = ERR_DB
		} else {
			presencelist := []*MsgPresence{}
			for _, presdata := range list {
				var pres *MsgPresence
				err = json.Unmarshal(presdata[7:], &pres)
				if err != nil {
					errcode = ERR_DB
					break
				}
				presencelist = append(presencelist, pres)
			}

			if err != nil {
				errcode = ERR_DB
			} else {
				ret.Json, err = json.Marshal(presencelist)
				if err != nil {
					errcode = ERR_UNKNOWN
					ret.Json = nil
				}
			}
		}
	case DataType_Black:
	case DataType_Offline_Message:
		list, err := dbMgr.GetOfflineMessage(sess.ID())
		if err != nil {
			errcode = ERR_DB
		} else {
			//msglist := []*MsgMessage{}
			for _, msgdata := range list {
				// var pres *MsgMessage
				// err = json.Unmarshal(msgdata[7:], &pres)
				// if err != nil {
				// 	errcode = ERR_DB
				// 	break
				// }
				// msglist = append(msglist, pres)
				sess.Send(msgdata)
			}

			// if err != nil {
			// 	errcode = ERR_DB
			// } else {
			// 	ret.Json, err = json.Marshal(msglist)
			// 	if err != nil {
			// 		errcode = ERR_UNKNOWN
			// 		ret.Json = nil
			// 	}
			// }
			return errcode, nil
		}
	case DataType_Room:
	case DataType_RoomMessage:
	}

	ret.ErrorCode = errcode
	return errcode, ret
}

func isInBlack(id, otherid uint64) (bool, error) {
	flag, err := gtdb.Manager().IsInBlack(id, otherid)
	if err != nil {
		return false, err
	}
	return flag, nil
}

func HandlerMessage(sess ISession, data []byte) (uint16, interface{}) {
	who := Uint64(data)
	timestamp := Int64(data[8:])
	message := data[16:]

	if who == sess.ID() {
		return ERR_MESSAGE_SELF, ERR_MESSAGE_SELF
	}

	timestamp = time.Now().Unix()

	msg := &MsgMessage{Who: sess.ID(), TimeStamp: timestamp, Message: message}

	errcode := ERR_NONE
	dbMgr := gtdb.Manager()
	flag, err := dbMgr.IsAppDataExists(who)

	if err != nil {
		errcode = ERR_DB
	} else {
		if !flag {
			errcode = ERR_APPDATAID_NOT_EXISTS
		} else {
			flag, err = isInBlack(sess.ID(), who)
			if err != nil {
				errcode = ERR_DB
			} else {
				if !flag {
					errcode = ERR_IN_BLACKLIST
				} else {
					msgbytes, err := json.Marshal(msg)
					if err != nil {
						errcode = ERR_UNKNOWN
					} else {
						senddata := packageMsg(RetFrame, 0, MsgId_Message, msgbytes)
						errcode = SendMessageToUser(who, senddata)
					}
				}
			}
		}
	}

	return errcode, errcode
}
