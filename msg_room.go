package main

import (
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

	registerMsgHandler(MsgId_ReqRoomList, HandlerReqRoomList)
	registerMsgHandler(MsgId_ReqRoomPresenceList, HandlerReqRoomPresenceList)
	registerMsgHandler(MsgId_ReqRoomUserList, HandlerReqRoomUserList)
}

func HandlerReqCreateRoom(sess ISession, data []byte) (uint16, interface{}) {
	errcode := ERR_NONE
	var roommsg *MsgReqCreateRoom = &MsgReqCreateRoom{}
	if !jsonUnMarshal(data, roommsg, &errcode) {
		return errcode, errcode
	}

	createRoom(sess.ID(), roommsg, &errcode)
	return errcode, errcode
}

func HandlerReqDeleteRoom(sess ISession, data []byte) (uint16, interface{}) {
	rid := Uint64(data)
	errcode := ERR_NONE

	if !isRoomExists(rid, &errcode) {
		return errcode, errcode
	}

	var uselist []*gtdb.RoomUser
	if !getRoomUserIds(rid, &uselist, &errcode) {
		return errcode, errcode
	}

	if !deleteRoom(rid, &errcode) {
		return errcode, errcode
	}

	//通知房间其他人，房间解散
	var presence *MsgRoomPresence = &MsgRoomPresence{}
	presence.PresenceType = PresenceType_Dismiss
	presence.Rid = rid
	presence.Who = sess.ID()
	presence.Nickname = sess.NickName()
	presence.TimeStamp = time.Now().Unix()

	var presencebytes []byte
	jsonMarshal(presence, &presencebytes, &errcode)
	senddata := packageMsg(RetFrame, 0, MsgId_RoomPresence, presencebytes)

	myid := sess.ID()
	for _, user := range uselist {
		if user.Dataid != myid {
			SendMessageToUser(user.Dataid, senddata)
		}
	}
	return errcode, errcode
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

	switch presencetype {
	case PresenceType_Subscribe:
		presence.Nickname = sess.NickName()
		presence.Who = sess.ID()

		if !isNotRoomUser(rid, sess.ID(), &errcode) {
			return errcode, errcode
		}

		if !isRoomNotFull(rid, &errcode) {
			return errcode, errcode
		}

		var roomtype byte
		if !getRoomType(rid, &roomtype, &errcode) {
			return errcode, errcode
		}

		if roomtype == RoomType_Apply {
			var presencebytes []byte
			if !jsonMarshal(presence, &presencebytes, &errcode) {
				return errcode, errcode
			}

			if !addRoomPresence(rid, sess.ID(), presencebytes, &errcode) {
				return errcode, errcode
			}

			if !sendPresenceToRoomAdmin(rid, presencebytes, &errcode) {
				return errcode, errcode
			}
		} else if roomtype == RoomType_Everyone {
			if !addRoomUser(rid, sess.ID(), presence, &errcode) {
				return errcode, errcode
			}
		} else if roomtype == RoomType_Password {
			if !isRoomPassword(rid, presence.Password, &errcode) {
				return errcode, errcode
			}

			if !addRoomUser(rid, sess.ID(), presence, &errcode) {
				return errcode, errcode
			}
		}
	case PresenceType_Subscribed:
		if !isRoomAdmin(rid, sess.ID(), &errcode) {
			return errcode, errcode
		}

		if !isAppDataExists(who, &errcode) {
			return errcode, errcode
		}

		if !isRoomPresenceExists(rid, who, &errcode) {
			return errcode, errcode
		}

		if !addRoomUser(rid, sess.ID(), presence, &errcode) {
			return errcode, errcode
		}

		if !removeRoomPresence(rid, sess.ID(), &errcode) {
			return errcode, errcode
		}
	case PresenceType_UnSubscribe:
		if !isRoomUser(rid, sess.ID(), &errcode) {
			return errcode, errcode
		}

		if !isNotRoomOwner(rid, sess.ID(), &errcode) {
			return errcode, errcode
		}

		if !removeRoomUser(rid, sess.ID(), &errcode) {
			return errcode, errcode
		}

		var presencebytes []byte
		if !jsonMarshal(presence, &presencebytes, &errcode) {
			return errcode, errcode
		}

		if !sendPresenceToRoomUser(rid, presencebytes, &errcode) {
			return errcode, errcode
		}

		if !removeRoomPresence(rid, sess.ID(), &errcode) {
			return errcode, errcode
		}
	case PresenceType_UnSubscribed:
		if !isRoomAdmin(rid, sess.ID(), &errcode) {
			return errcode, errcode
		}

		if !isAppDataExists(who, &errcode) {
			return errcode, errcode
		}

		if !isRoomPresenceExists(rid, who, &errcode) {
			return errcode, errcode
		}

		var presencebytes []byte
		if !jsonMarshal(presence, &presencebytes, &errcode) {
			return errcode, errcode
		}

		senddata := packageMsg(RetFrame, 0, MsgId_RoomPresence, presencebytes)
		errcode = SendMessageToUser(who, senddata)

		if !removeRoomPresence(rid, sess.ID(), &errcode) {
			return errcode, errcode
		}
	case PresenceType_Available, PresenceType_UnAvailable, PresenceType_Invisible:
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
	errcode := ERR_NONE
	var roomsetting *MsgReqUpdateRoomSetting = &MsgReqUpdateRoomSetting{}
	if !jsonUnMarshal(data, roomsetting, &errcode) {
		return errcode, errcode
	}

	dbMgr := gtdb.Manager()

	if roomsetting.Bit&RoomSetting_RoomName != 0 {
		err := dbMgr.SetRoomName(roomsetting.Rid, roomsetting.RoomName)
		if err != nil {
			return ERR_DB, ERR_DB
		}
	}

	if roomsetting.Bit&RoomSetting_RoomType != 0 {
		err := dbMgr.SetRoomType(roomsetting.Rid, roomsetting.RoomType)
		if err != nil {
			return ERR_DB, ERR_DB
		}
	}

	if roomsetting.Bit&RoomSetting_Jieshao != 0 {
		err := dbMgr.SetRoomJieshao(roomsetting.Rid, roomsetting.Jieshao)
		if err != nil {
			return ERR_DB, ERR_DB
		}
	}

	if roomsetting.Bit&RoomSetting_Notice != 0 {
		err := dbMgr.SetRoomNotice(roomsetting.Rid, roomsetting.Notice)
		if err != nil {
			return ERR_DB, ERR_DB
		}
	}

	if roomsetting.Bit&RoomSetting_Password != 0 {
		err := dbMgr.SetRoomPassword(roomsetting.Rid, roomsetting.Password)
		if err != nil {
			return ERR_DB, ERR_DB
		}
	}

	return errcode, errcode
}

func HandlerReqBanRoomUser(sess ISession, data []byte) (uint16, interface{}) {
	rid := Uint64(data)
	appdataid := Uint64(data[8:])
	errcode := ERR_NONE

	if !isRoomExists(rid, &errcode) {
		return errcode, errcode
	}

	if !isAppDataExists(appdataid, &errcode) {
		return errcode, errcode
	}

	removeRoomUser(rid, appdataid, &errcode)

	return errcode, errcode
}

func HandlerReqJinyanRoomUser(sess ISession, data []byte) (uint16, interface{}) {
	rid := Uint64(data)
	appdataid := Uint64(data[8:])
	errcode := ERR_NONE

	if !isRoomExists(rid, &errcode) {
		return errcode, errcode
	}

	if !isAppDataExists(appdataid, &errcode) {
		return errcode, errcode
	}

	jinyanRoomUser(rid, appdataid, &errcode)

	return errcode, errcode
}

func HandlerReqUnJinyanRoomUser(sess ISession, data []byte) (uint16, interface{}) {
	rid := Uint64(data)
	appdataid := Uint64(data[8:])
	errcode := ERR_NONE

	if !isRoomExists(rid, &errcode) {
		return errcode, errcode
	}

	if !isAppDataExists(appdataid, &errcode) {
		return errcode, errcode
	}

	unjinyanRoomUser(rid, appdataid, &errcode)

	return errcode, errcode
}

func HandlerReqAddRoomAdmin(sess ISession, data []byte) (uint16, interface{}) {
	rid := Uint64(data)
	appdataid := Uint64(data[8:])
	errcode := ERR_NONE

	if !isRoomExists(rid, &errcode) {
		return errcode, errcode
	}

	if !isAppDataExists(appdataid, &errcode) {
		return errcode, errcode
	}

	addRoomAdmin(rid, appdataid, &errcode)

	return errcode, errcode
}

func HandlerReqRemoveRoomAdmin(sess ISession, data []byte) (uint16, interface{}) {
	rid := Uint64(data)
	appdataid := Uint64(data[8:])
	errcode := ERR_NONE

	if !isRoomExists(rid, &errcode) {
		return errcode, errcode
	}

	if !isAppDataExists(appdataid, &errcode) {
		return errcode, errcode
	}

	removeRoomAdmin(rid, appdataid, &errcode)

	return errcode, errcode
}

func HandlerRoomMessage(sess ISession, data []byte) (uint16, interface{}) {
	errcode := ERR_NONE
	var roommsg *MsgRoomMessage = &MsgRoomMessage{}
	if !jsonUnMarshal(data, roommsg, &errcode) {
		return errcode, errcode
	}

	if !isRoomExists(roommsg.Rid, &errcode) {
		return errcode, errcode
	}

	if !isRoomUser(roommsg.Rid, sess.ID(), &errcode) {
		return errcode, errcode
	}

	roommsg.TimeStamp = time.Now().Unix()
	roommsg.Who = sess.ID()
	roommsg.Nickname = sess.NickName()

	var msgbytes []byte
	if !jsonMarshal(roommsg, &msgbytes, &errcode) {
		return errcode, errcode
	}

	//最后一个if判断不需要
	// if !sendMessageToRoomUser(roommsg.Rid, msgbytes, &errcode) {
	// 	return errcode, errcode
	// }
	sendMessageToRoomUser(roommsg.Rid, msgbytes, &errcode)

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

func HandlerReqRoomList(sess ISession, data []byte) (uint16, interface{}) {
	errcode := ERR_NONE

	var roomlist []*gtdb.Room
	if !getRoomList(sess.ID(), &roomlist, &errcode) {
		return errcode, errcode
	}

	var msgbytes []byte
	if !jsonMarshal(roomlist, &msgbytes, &errcode) {
		return errcode, errcode
	}

	ret := &MsgRetRoomList{}
	ret.Json = msgbytes
	ret.ErrorCode = errcode
	return errcode, ret
}

func HandlerReqRoomPresenceList(sess ISession, data []byte) (uint16, interface{}) {
	rid := Uint64(data)
	errcode := ERR_NONE

	if !isRoomExists(rid, &errcode) {
		return errcode, errcode
	}
	if !isRoomAdmin(rid, sess.ID(), &errcode) {
		return errcode, errcode
	}
	var datalist map[string]string
	if !getRoomPresenceList(rid, &datalist, &errcode) {
		return errcode, errcode
	}

	var msgbytes []byte
	if !jsonMarshal(datalist, &msgbytes, &errcode) {
		return errcode, errcode
	}

	ret := &MsgRetRoomPresenceList{}
	ret.Json = msgbytes
	ret.ErrorCode = errcode
	return errcode, ret
}

func HandlerReqRoomUserList(sess ISession, data []byte) (uint16, interface{}) {
	rid := Uint64(data)
	errcode := ERR_NONE

	if !isRoomUser(rid, sess.ID(), &errcode) {
		return errcode, errcode
	}

	var datalist []*gtdb.RoomUser
	if !getRoomUserList(rid, &datalist, &errcode) {
		return errcode, errcode
	}

	var msgbytes []byte
	if !jsonMarshal(datalist, &msgbytes, &errcode) {
		return errcode, errcode
	}

	ret := &MsgRetRoomUserList{}
	ret.Json = msgbytes
	ret.ErrorCode = errcode
	return errcode, ret
}
