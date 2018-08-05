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
	RegisterRoomMsg()
}

func RegisterRoomMsg() {
	registerMsgHandler(MsgId_ReqCreateRoom, HandlerReqCreateRoom)
	registerMsgHandler(MsgId_ReqDeleteRoom, HandlerReqDeleteRoom)
	registerMsgHandler(MsgId_RoomPresence, HandlerReqRoomPresence)
	registerMsgHandler(MsgId_ReqUpdateRoomSetting, HandlerReqUpdateRoomSetting)
	registerMsgHandler(MsgId_ReqBanRoomUser, HandlerReqBanRoomUser)
	registerMsgHandler(MsgId_ReqJinyanRoomUser, HandlerReqJinyanRoomUser)
	registerMsgHandler(MsgId_ReqUnJinyanRoomUser, HandlerReqUnJinyanRoomUser)
	registerMsgHandler(MsgId_ReqAddRoomAdmin, HandlerReqAddRoomAdmin)
	registerMsgHandler(MsgId_ReqRemoveRoomAdmin, HandlerReqRemoveRoomAdmin)
	registerMsgHandler(MsgId_RoomMessage, HandlerRoomMessage)
}

func HandlerReqCreateRoom(sess ISession, data []byte) (uint16, interface{}) {
	var roommsg *MsgReqCreateRoom = &MsgReqCreateRoom{}
	err := json.Unmarshal(data, roommsg)

	fmt.Println(string(data))
	fmt.Println(roommsg)
	if err != nil {
		fmt.Println(err.Error())
		return ERR_INVALID_JSON, ERR_INVALID_JSON
	}

	errcode := ERR_NONE
	dbMgr := gtdb.Manager()

	tbl_room := &gtdb.Room{Ownerid: sess.ID(), Roomname: roommsg.Name, Roomtype: roommsg.RoomType, Jieshao: roommsg.Jieshao, Notice: roommsg.Notice, Password: roommsg.Password}
	err = dbMgr.CreateRoom(tbl_room)

	if err != nil {
		errcode = ERR_DB
	}

	return errcode, errcode
}

func HandlerReqDeleteRoom(sess ISession, data []byte) (uint16, interface{}) {
	rid := Uint64(data)

	dbMgr := gtdb.Manager()
	flag, err := dbMgr.IsRoomExists(rid)
	errcode := ERR_NONE

	if err != nil {
		errcode = ERR_DB
	} else {
		if !flag {
			errcode = ERR_ROOM_NOT_EXISTS
		} else {
			err = dbMgr.DeleteRoom(rid)

			if err != nil {
				errcode = ERR_DB
			}
		}
	}

	return errcode, errcode
}

func HandlerReqRoomPresence(sess ISession, data []byte) (uint16, interface{}) {
	var presence *MsgRoomPresence = &MsgRoomPresence{}
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
	flag, err := dbMgr.IsRoomExists(who)

	if err != nil {
		errcode = ERR_DB
	} else {
		if !flag {
			errcode = ERR_ROOM_NOT_EXISTS
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

func HandlerReqUpdateRoomSetting(sess ISession, data []byte) (uint16, interface{}) {
	var roomsetting *MsgReqUpdateRoomSetting = &MsgReqUpdateRoomSetting{}
	err := json.Unmarshal(data, roomsetting)

	fmt.Println(string(data))
	fmt.Println(roomsetting)
	if err != nil {
		fmt.Println(err.Error())
		return ERR_INVALID_JSON, ERR_INVALID_JSON
	}

	errcode := ERR_NONE
	dbMgr := gtdb.Manager()

	if roomsetting.Bit&RoomSetting_RoomName != 0 {
		dbMgr.SetRoomName(roomsetting.Rid, roomsetting.RoomName)
		if err != nil {
			return ERR_DB, ERR_DB
		}
	}

	if roomsetting.Bit&RoomSetting_RoomType != 0 {
		dbMgr.SetRoomType(roomsetting.Rid, roomsetting.RoomType)
		if err != nil {
			return ERR_DB, ERR_DB
		}
	}

	if roomsetting.Bit&RoomSetting_Jieshao != 0 {
		dbMgr.SetRoomJieshao(roomsetting.Rid, roomsetting.Jieshao)
		if err != nil {
			return ERR_DB, ERR_DB
		}
	}

	if roomsetting.Bit&RoomSetting_Notice != 0 {
		dbMgr.SetRoomNotice(roomsetting.Rid, roomsetting.Notice)
		if err != nil {
			return ERR_DB, ERR_DB
		}
	}

	if roomsetting.Bit&RoomSetting_Password != 0 {
		dbMgr.SetRoomPassword(roomsetting.Rid, roomsetting.Password)
		if err != nil {
			return ERR_DB, ERR_DB
		}
	}

	return errcode, errcode
}

func HandlerReqBanRoomUser(sess ISession, data []byte) (uint16, interface{}) {
	rid := Uint64(data)
	appdataid := Uint64(data[8:])

	dbMgr := gtdb.Manager()
	flag, err := dbMgr.IsRoomExists(rid)
	errcode := ERR_NONE

	if err != nil {
		errcode = ERR_DB
	} else {
		if !flag {
			errcode = ERR_ROOM_NOT_EXISTS
		} else {
			flag, err = dbMgr.IsAppDataExists(appdataid)

			if err != nil {
				errcode = ERR_DB
			} else {
				if !flag {
					errcode = ERR_ROOM_NOT_EXISTS
				} else {
					err = dbMgr.RemoveRoomUser(rid, appdataid)

					if err != nil {
						errcode = ERR_DB
					}

					//TODO:通知管理员
				}
			}
		}
	}

	return errcode, errcode
}

func HandlerReqJinyanRoomUser(sess ISession, data []byte) (uint16, interface{}) {
	rid := Uint64(data)
	appdataid := Uint64(data[8:])

	dbMgr := gtdb.Manager()
	flag, err := dbMgr.IsRoomExists(rid)
	errcode := ERR_NONE

	if err != nil {
		errcode = ERR_DB
	} else {
		if !flag {
			errcode = ERR_ROOM_NOT_EXISTS
		} else {
			flag, err = dbMgr.IsAppDataExists(appdataid)

			if err != nil {
				errcode = ERR_DB
			} else {
				if !flag {
					errcode = ERR_ROOM_NOT_EXISTS
				} else {
					err = dbMgr.JinyanUserInRoom(rid, appdataid)

					if err != nil {
						errcode = ERR_DB
					}

					//TODO:通知在线成员
				}
			}
		}
	}

	return errcode, errcode
}

func HandlerReqUnJinyanRoomUser(sess ISession, data []byte) (uint16, interface{}) {
	rid := Uint64(data)
	appdataid := Uint64(data[8:])

	dbMgr := gtdb.Manager()
	flag, err := dbMgr.IsRoomExists(rid)
	errcode := ERR_NONE

	if err != nil {
		errcode = ERR_DB
	} else {
		if !flag {
			errcode = ERR_ROOM_NOT_EXISTS
		} else {
			flag, err = dbMgr.IsAppDataExists(appdataid)

			if err != nil {
				errcode = ERR_DB
			} else {
				if !flag {
					errcode = ERR_ROOM_NOT_EXISTS
				} else {
					err = dbMgr.UnJinyanUserInRoom(rid, appdataid)

					if err != nil {
						errcode = ERR_DB
					}
				}
			}
		}
	}

	return errcode, errcode
}

func HandlerReqAddRoomAdmin(sess ISession, data []byte) (uint16, interface{}) {
	rid := Uint64(data)
	appdataid := Uint64(data[8:])

	dbMgr := gtdb.Manager()
	flag, err := dbMgr.IsRoomExists(rid)
	errcode := ERR_NONE

	if err != nil {
		errcode = ERR_DB
	} else {
		if !flag {
			errcode = ERR_ROOM_NOT_EXISTS
		} else {
			flag, err = dbMgr.IsAppDataExists(appdataid)

			if err != nil {
				errcode = ERR_DB
			} else {
				if !flag {
					errcode = ERR_ROOM_NOT_EXISTS
				} else {
					err = dbMgr.AddRoomAdmin(rid, appdataid)

					if err != nil {
						errcode = ERR_DB
					}
				}
			}
		}
	}

	return errcode, errcode
}

func HandlerReqRemoveRoomAdmin(sess ISession, data []byte) (uint16, interface{}) {
	rid := Uint64(data)
	appdataid := Uint64(data[8:])

	dbMgr := gtdb.Manager()
	flag, err := dbMgr.IsRoomExists(rid)
	errcode := ERR_NONE

	if err != nil {
		errcode = ERR_DB
	} else {
		if !flag {
			errcode = ERR_ROOM_NOT_EXISTS
		} else {
			flag, err = dbMgr.IsAppDataExists(appdataid)

			if err != nil {
				errcode = ERR_DB
			} else {
				if !flag {
					errcode = ERR_ROOM_NOT_EXISTS
				} else {
					err = dbMgr.RemoveRoomAdmin(rid, appdataid)

					if err != nil {
						errcode = ERR_DB
					}
				}
			}
		}
	}

	return errcode, errcode
}

func HandlerRoomMessage(sess ISession, data []byte) (uint16, interface{}) {
	var roommsg *MsgRoomMessage = &MsgRoomMessage{}
	err := json.Unmarshal(data, roommsg)

	fmt.Println(string(data))
	fmt.Println(roommsg)
	if err != nil {
		fmt.Println(err.Error())
		return ERR_INVALID_JSON, ERR_INVALID_JSON
	}

	dbMgr := gtdb.Manager()
	flag, err := dbMgr.IsRoomExists(roommsg.Rid)
	errcode := ERR_NONE

	if err != nil {
		errcode = ERR_DB
	} else {
		if !flag {
			errcode = ERR_ROOM_NOT_EXISTS
		} else {
			flag, err = dbMgr.IsUserInRoom(roommsg.Rid, sess.ID())

			if err != nil {
				errcode = ERR_DB
			} else {
				if !flag {
					errcode = ERR_ROOM_USER_INVALID
				} else {
					userlist, err := dbMgr.GetRoomUserList(roommsg.Rid)

					if err != nil {
						errcode = ERR_DB
					} else {
						for _, user := range userlist {
							//broadcast to user in room
						}
					}
				}
			}
		}
	}

	return errcode, errcode
}
