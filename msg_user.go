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
	registerMsgHandler(MsgId_ReqUpdateAppdata, HandlerReqUpdateAppdata)
}

func HandlerReqUserData(sess ISession, data []byte) (uint16, interface{}) {
	id := Uint64(data)
	if id == 0 {
		id = sess.ID()
	}
	appdata, err := gtdb.Manager().GetAppData(id)
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
		fmt.Println("SendMessageToUserOnline to ", to, " serveraddr ", online.Serveraddr)

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

	presencetype := presence.PresenceType
	who := presence.Who
	//timestamp := Int64(data[9:])
	//message := data[17:]

	if who == sess.ID() {
		return ERR_FRIEND_SELF, ERR_FRIEND_SELF
	}

	presence.Nickname = sess.NickName()
	presence.TimeStamp = time.Now().Unix()
	presence.Who = sess.ID()

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
						tbl_from := &gtdb.Friend{Dataid: who, Otherdataid: sess.ID(), Groupname: config.DefaultGroupName}
						tbl_to := &gtdb.Friend{Dataid: sess.ID(), Otherdataid: who, Groupname: config.DefaultGroupName}
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
		grouplist, err := dbMgr.GetGroupList(sess.ID())
		if err != nil {
			errcode = ERR_DB
		} else {
			friendlist := map[string][]*gtdb.FriendJson{}
			for _, group := range grouplist {
				list, err := dbMgr.GetFriendInfoList(sess.ID(), group)
				if err != nil {
					errcode = ERR_DB
				} else {
					friendlist[group] = list
				}
			}
			ret.Json, err = json.Marshal(friendlist)
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
			fmt.Println(list)
			presencelist := []*MsgPresence{}
			for _, presstr := range list {
				var pres *MsgPresence = &MsgPresence{}
				presdata := []byte(presstr)
				err = json.Unmarshal(presdata, pres)
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
				fmt.Println(string(ret.Json))
				if err != nil {
					errcode = ERR_UNKNOWN
					ret.Json = nil
				}
			}
		}
	case DataType_Black:
		blacklist, err := dbMgr.GetBlackInfoList(sess.ID())
		if err != nil {
			errcode = ERR_DB
		} else {
			ret.Json, err = json.Marshal(blacklist)
			fmt.Println(string(ret.Json))
			if err != nil {
				errcode = ERR_JSON_SERIALIZE
				ret.Json = nil
			}
		}
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

func HandlerMessage(sess ISession, data []byte) (uint16, interface{}) {
	var msg *MsgMessageJson = &MsgMessageJson{}
	err := json.Unmarshal(data, msg)

	if err != nil {
		fmt.Println(err.Error())
		return ERR_INVALID_JSON, ERR_INVALID_JSON
	}

	who := msg.Who

	if who == sess.ID() {
		return ERR_MESSAGE_SELF, ERR_MESSAGE_SELF
	}

	msg.TimeStamp = time.Now().Unix()
	msg.Who = sess.ID()
	msg.Nickname = sess.NickName()

	//msg := &MsgMessage{Who: sess.ID(), TimeStamp: timestamp, Message: message}

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
				if flag {
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

func HandlerReqUpdateAppdata(sess ISession, data []byte) (uint16, interface{}) {
	msgmap := make(map[string]interface{})
	err := json.Unmarshal(data, msgmap)
	errcode := ERR_NONE

	if err != nil {
		fmt.Println(err.Error())
		return ERR_INVALID_JSON, ERR_INVALID_JSON
	}

	//TODO: 字段校验暂时未做
	updatemap := make(map[string]interface{})

	value, ok := msgmap["nickname"]

	if ok {
		nick, ok := value.(string)
		if ok && nick != "" {
			updatemap["nickname"] = nick
		}
	}

	value, ok = msgmap["desc"]

	if ok {
		desc, ok := value.(string)
		if ok {
			updatemap["desc"] = desc
		}
	}

	value, ok = msgmap["sex"]

	if ok {
		sex, ok := value.(string)
		if ok && sex != "" {
			updatemap["sex"] = sex
		}
	}

	value, ok = msgmap["birthday"]

	if ok {
		strbirthday, ok := value.(string)
		if ok && strbirthday != "" {
			birthday, err := time.Parse("01/02/2006", strbirthday)
			if err == nil {
				updatemap["birthday"] = birthday
			}
		}
	}

	value, ok = msgmap["country"]

	if ok {
		country, ok := value.(string)
		if ok && country != "" {
			updatemap["country"] = country
		}
	}

	err = gtdb.Manager().UpdateAppDataByMap(updatemap)

	if err != nil {
		errcode = ERR_DB
	}

	return errcode, errcode
}
