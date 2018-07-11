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
	registerMsgHandler(MsgId_Group, HandlerGroup)
	registerMsgHandler(MsgId_ReqGroupRefresh, HandlerGroupRefresh)
	//registerMsgHandler(MsgId_EnterChat, HandlerEnterChat)
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
	case "refresh":
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
	}
	ret.ErrorCode = errcode

	return errcode, ret
}
