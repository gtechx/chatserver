package main

import (
	"encoding/json"
	"fmt"
	"time"

	. "github.com/gtechx/base/common"
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

//正条件，在条件满足的时候才去做事情
func isRoomFull(rid uint64, perrcode *uint16) bool {
	dbMgr := gtdb.Manager()

	usercount, err := dbMgr.GetRoomUserCount(rid)
	if err != nil {
		*perrcode = ERR_DB
	} else {
		maxusercount, err := dbMgr.GetRoomMaxUser(rid)
		if err != nil {
			*perrcode = ERR_DB
		} else {
			if usercount != maxusercount {
				*perrcode = ERR_ROOM_NOT_FULL
			} else {
				return true
			}
		}
	}

	return false
}

//反条件，在条件满足的时候才去做事情
func isNotRoomFull(rid uint64, perrcode *uint16) bool {
	dbMgr := gtdb.Manager()

	usercount, err := dbMgr.GetRoomUserCount(rid)
	if err != nil {
		*perrcode = ERR_DB
	} else {
		maxusercount, err := dbMgr.GetRoomMaxUser(rid)
		if err != nil {
			*perrcode = ERR_DB
		} else {
			if usercount == maxusercount {
				*perrcode = ERR_ROOM_FULL
			} else {
				return true
			}
		}
	}

	return false
}

func addRoomUser(rid, appdataid uint64, presence *MsgRoomPresence) uint16 {
	errcode := ERR_NONE
	dbMgr := gtdb.Manager()

	tbl_roomuser := &gtdb.RoomUser{Rid: rid, Dataid: appdataid}
	err := dbMgr.AddRoomUser(tbl_roomuser)

	if err != nil {
		errcode = ERR_DB
	} else {
		presencebytes, err := json.Marshal(presence)
		if err != nil {
			errcode = ERR_INVALID_JSON
		} else {
			senddata := packageMsg(RetFrame, 0, MsgId_RoomPresence, presencebytes)
			userlist, err := dbMgr.GetRoomUserIds(rid)

			if err != nil {
				errcode = ERR_DB
			} else {
				for _, user := range userlist {
					//broadcast to user in room
					errcode = SendMessageToUser(user, senddata)
				}
			}
		}
	}

	return errcode
}

func isRoomPassword(rid uint64, password string, perrcode *uint16) bool {
	roompassword, err := gtdb.Manager().GetRoomPassword(rid)

	if err != nil {
		*perrcode = ERR_DB
	} else {
		if password != roompassword {
			*perrcode = ERR_ROOM_PASSWORD_INVALID
		} else {
			return true
		}
	}

	return false
}

func isRoomExists(rid uint64, perrcode *uint16) bool {
	flag, err := gtdb.Manager().IsRoomExists(rid)

	if err != nil {
		*perrcode = ERR_DB
	} else {
		if !flag {
			*perrcode = ERR_ROOM_NOT_EXISTS
		} else {
			return true
		}
	}

	return false
}

func isRoomNotExists(rid uint64, perrcode *uint16) bool {
	flag, err := gtdb.Manager().IsRoomExists(rid)

	if err != nil {
		*perrcode = ERR_DB
	} else {
		if flag {
			*perrcode = ERR_ROOM_EXISTS
		} else {
			return true
		}
	}

	return false
}

func isRoomUser(rid, appdataid uint64, perrcode *uint16) bool {
	flag, err := gtdb.Manager().IsRoomUser(rid, appdataid)

	if err != nil {
		*perrcode = ERR_DB
	} else {
		if !flag {
			*perrcode = ERR_ROOM_USER_INVALID
		} else {
			return true
		}
	}

	return false
}

func isNotRoomUser(rid, appdataid uint64, perrcode *uint16) bool {
	flag, err := gtdb.Manager().IsRoomUser(rid, appdataid)

	if err != nil {
		*perrcode = ERR_DB
	} else {
		if flag {
			*perrcode = ERR_ROOM_USER_EXISTS
		} else {
			return true
		}
	}

	return false
}

func HandlerReqRoomPresence(sess ISession, data []byte) (uint16, interface{}) {
	errcode := ERR_NONE
	var presence *MsgRoomPresence = &MsgRoomPresence{}
	if !jsonUnMarshal(data, presence, &errcode) {
		return errcode, errcode
	}

	if !isRoomExists(presence.Rid, &errcode) {
		return errcode, errcode
	}

	presencetype := presence.PresenceType
	rid := presence.Rid
	who := presence.Who

	presence.TimeStamp = time.Now().Unix()

	dbMgr := gtdb.Manager()

	switch presencetype {
	case PresenceType_Subscribe:
		presence.Nickname = sess.NickName()
		presence.Who = sess.ID()

		if isRoomUser(rid, sess.ID(), &errcode) {
			return errcode, errcode
		}

		if isRoomFull(rid, &errcode) {
			return errcode, errcode
		}

		roomtype, _ := dbMgr.GetRoomType(rid)

		if roomtype == RoomType_Apply {
			var presencebytes []byte
			if !jsonMarshal(presence, &presencebytes, &errcode) {
				return errcode, errcode
			}

			var admins []uint64
			if !getRoomAdminIds(rid, &admins, &errcode) {
				return errcode, errcode
			}
			senddata := packageMsg(RetFrame, 0, MsgId_RoomPresence, presencebytes)

			err = dbMgr.AddRoomPresence(rid, sess.ID(), presencebytes)
			if err != nil {
				errcode = ERR_DB
			} else {
				//send to admin
				for _, id := range admins {
					errcode = SendMessageToUser(id, senddata)
				}
			}
		} else if roomtype == RoomType_Everyone {
			errcode = addRoomUser(rid, sess.ID(), presence)
		} else if roomtype == RoomType_Password {
			if isRoomPassword(rid, presence.Password, &errcode) {
				errcode = addRoomUser(rid, sess.ID(), presence)
			}
		}
	case PresenceType_Subscribed:
		flag, err = dbMgr.IsRoomAdmin(rid, sess.ID())
		if err != nil {
			errcode = ERR_DB
		} else {
			if !flag {
				errcode = ERR_ROOM_ADMIN_INVALID
			} else {
				flag, err = dbMgr.IsAppDataExists(who)

				if err != nil {
					errcode = ERR_DB
				} else {
					if !flag {
						errcode = ERR_APPDATAID_NOT_EXISTS
					} else {
						flag, err = dbMgr.IsRoomPresenceExists(rid, who)

						if err != nil {
							errcode = ERR_DB
						} else {
							if !flag {
								errcode = ERR_ROOM_PRESENCE_NOT_EXISTS
							} else {
								tbl_roomuser := &gtdb.RoomUser{Rid: rid, Dataid: who}
								err = dbMgr.AddRoomUser(tbl_roomuser)

								if err != nil {
									errcode = ERR_DB
								} else {
									presencebytes, err := json.Marshal(presence)
									if err != nil {
										errcode = ERR_INVALID_JSON
									} else {
										senddata := packageMsg(RetFrame, 0, MsgId_RoomPresence, presencebytes)
										userlist, err := dbMgr.GetRoomUserIds(rid)

										if err != nil {
											errcode = ERR_DB
										} else {
											for _, user := range userlist {
												//broadcast to user in room
												errcode = SendMessageToUser(user, senddata)
											}
										}
										dbMgr.RemoveRoomPresence(rid, who)
									}
								}
							}
						}
					}
				}
			}
		}
	case PresenceType_Unsubscribe:
		//check if the two are friend, if not omit thie message, else delete friend and send to who.
		flag, err = dbMgr.IsRoomUser(rid, sess.ID())
		if err != nil {
			errcode = ERR_DB
		} else {
			if !flag {
				errcode = ERR_ROOM_USER_INVALID
			} else {
				err = dbMgr.RemoveRoomUser(rid, sess.ID())
				if err != nil {
					errcode = ERR_DB
				} else {
					presencebytes, err := json.Marshal(presence)
					if err != nil {
						errcode = ERR_INVALID_JSON
					} else {
						senddata := packageMsg(RetFrame, 0, MsgId_RoomPresence, presencebytes)
						userlist, err := dbMgr.GetRoomUserIds(rid)

						if err != nil {
							errcode = ERR_DB
						} else {
							for _, user := range userlist {
								//broadcast to user in room
								errcode = SendMessageToUser(user, senddata)
							}
						}
						dbMgr.RemoveRoomPresence(rid, who)
					}
				}
			}
		}
	case PresenceType_Unsubscribed:
		flag, err = dbMgr.IsRoomAdmin(rid, sess.ID())
		if err != nil {
			errcode = ERR_DB
		} else {
			if !flag {
				errcode = ERR_ROOM_ADMIN_INVALID
			} else {
				flag, err = dbMgr.IsAppDataExists(who)

				if err != nil {
					errcode = ERR_DB
				} else {
					if !flag {
						errcode = ERR_APPDATAID_NOT_EXISTS
					} else {
						flag, err = dbMgr.IsRoomPresenceExists(rid, who)

						if err != nil {
							errcode = ERR_DB
						} else {
							if !flag {
								errcode = ERR_ROOM_PRESENCE_NOT_EXISTS
							} else {
								presencebytes, err := json.Marshal(presence)
								if err != nil {
									errcode = ERR_INVALID_JSON
								} else {
									senddata := packageMsg(RetFrame, 0, MsgId_RoomPresence, presencebytes)
									errcode = SendMessageToUser(who, senddata)
									dbMgr.RemoveRoomPresence(rid, who)
								}
							}
						}
					}
				}
			}
		}
	case PresenceType_Available, PresenceType_Unavailable, PresenceType_Invisible:
		//send to my friend online
		// presencebytes, err := json.Marshal(presence)
		// if err != nil {
		// 	errcode = ERR_INVALID_JSON
		// } else {
		// 	senddata := packageMsg(RetFrame, 0, MsgId_Presence, presencebytes)
		// 	SendMessageToFriendsOnline(sess.ID(), senddata)
		// }
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
					err = dbMgr.JinyanRoomUser(rid, appdataid)

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
					err = dbMgr.UnJinyanRoomUser(rid, appdataid)

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
	errcode := ERR_NONE
	var roommsg *MsgRoomMessage = &MsgRoomMessage{}
	if !jsonUnMarshal(data, roommsg, &errcode) {
		return errcode, errcode
	}
	// err := json.Unmarshal(data, roommsg)

	// fmt.Println(string(data))
	// fmt.Println(roommsg)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return ERR_INVALID_JSON, ERR_INVALID_JSON
	// }

	if isRoomNotExists(roommsg.Rid, &errcode) {
		return errcode, errcode
	}

	if isNotRoomUser(roommsg.Rid, sess.ID(), &errcode) {
		return errcode, errcode
	}

	var userlist []uint64
	if !getRoomUserIds(roommsg.Rid, &userlist, &errcode) {
		return errcode, errcode
	}

	var msgbytes []byte
	if !jsonMarshal(roommsg, &msgbytes, &errcode) {
		return errcode, errcode
	}

	roommsg.TimeStamp = time.Now().Unix()
	roommsg.Who = sess.ID()
	roommsg.Nickname = sess.NickName()

	senddata := packageMsg(RetFrame, 0, MsgId_RoomMessage, msgbytes)
	for _, user := range userlist {
		//broadcast to user in room
		errcode = SendMessageToUser(user, senddata)
	}

	// if getRoomUserIds(roommsg.Rid, &userlist, &errcode) {
	// 	var msgbytes []byte
	// 	if jsonMarshal(roommsg, &msgbytes, &errcode) {
	// 		senddata := packageMsg(RetFrame, 0, MsgId_RoomMessage, msgbytes)
	// 		for _, user := range userlist {
	// 			//broadcast to user in room
	// 			errcode = SendMessageToUser(user, senddata)
	// 		}
	// 	}
	// }

	// userlist, err := gtdb.Manager().GetRoomUserIds(roommsg.Rid)

	// if err != nil {
	// 	errcode = ERR_DB
	// } else {
	// 	msgbytes, err := json.Marshal(roommsg)
	// 	if err != nil {
	// 		errcode = ERR_JSON_SERIALIZE
	// 	} else {
	// 		senddata := packageMsg(RetFrame, 0, MsgId_RoomMessage, msgbytes)
	// 		for _, user := range userlist {
	// 			//broadcast to user in room
	// 			errcode = SendMessageToUser(user, senddata)
	// 		}
	// 	}
	// }

	return errcode, errcode
}

func getRoomUserIds(rid uint64, ids *[]uint64, perrcode *uint16) bool {
	userlist, err := gtdb.Manager().GetRoomUserIds(rid)

	if err != nil {
		*perrcode = ERR_DB
	} else {
		*ids = userlist
		return true
	}

	return false
}

func getRoomAdminIds(rid uint64, ids *[]uint64, perrcode *uint16) bool {
	userlist, err := gtdb.Manager().GetRoomAdminIds(rid)

	if err != nil {
		*perrcode = ERR_DB
	} else {
		*ids = userlist
		return true
	}

	return false
}

func jsonMarshal(data interface{}, out *[]byte, perrcode *uint16) bool {
	databytes, err := json.Marshal(data)
	if err != nil {
		*perrcode = ERR_JSON_SERIALIZE
	} else {
		*out = databytes
		return true
	}

	return false
}

func jsonUnMarshal(data []byte, out interface{}, perrcode *uint16) bool {
	err := json.Unmarshal(data, out)
	if err == nil {
		return true
	}
	*perrcode = ERR_INVALID_JSON
	return false
}
