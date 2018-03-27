package main

// import (
// 	. "github.com/nature19862001/Chat/protocol"
// 	. "github.com/nature19862001/base/common"
// )

// func OnReqFriendList(client *ChatClient, data []byte) {
// 	flist, ret := gDataManager.getFriendList(client.uid)

// 	// if err != nil {
// 	// 	reterr := NewErrorMsg(ERR_REDIS, MsgId_ReqFriendList)
// 	// 	client.send(Bytes(reterr))
// 	// 	return
// 	// }

// 	retmsg := new(MsgRetFriendList)
// 	retmsg.MsgId = MsgId_RetFriendList
// 	retmsg.Result = uint16(ret)
// 	retmsg.GroupCount = byte(len(flist))

// 	data = []byte{}

// 	for groupname, uidarr := range flist {
// 		data = append(data, byte(len(groupname)))
// 		data = append(data, []byte(groupname)...)
// 		data = append(data, Bytes(uint16(len(uidarr)))...)

// 		for _, uid := range uidarr {
// 			data = append(data, Bytes(uid)...)
// 		}
// 	}

// 	retmsg.Data = data

// 	client.send(Bytes(retmsg))
// }

// func OnReqFriendAdd(client *ChatClient, data []byte) {
// 	fuid := Uint64(data[2:])
// 	group := data[10:]
// 	ret := gDataManager.addFriend(client.uid, fuid, string(group))

// 	switch ret {
// 	// case ERR_FRIEND_IN_BLACKLIST:
// 	// 	//do nothing
// 	// case ERR_REDIS:
// 	// 	reterr := NewErrorMsg(ERR_REDIS, MsgId_ReqFriendAdd)
// 	// 	client.send(Bytes(reterr))
// 	// case ERR_FRIEND_EXIST:
// 	// 	reterr := NewErrorMsg(ERR_FRIEND_EXIST, MsgId_ReqFriendAdd)
// 	// 	client.send(Bytes(reterr))
// 	// case ERR_USER_NOT_EXIST:
// 	// 	reterr := NewErrorMsg(ERR_USER_NOT_EXIST, MsgId_ReqFriendAdd)
// 	// 	client.send(Bytes(reterr))
// 	// case ERR_FRIEND_MAX:
// 	// 	reterr := NewErrorMsg(ERR_FRIEND_MAX, MsgId_ReqFriendAdd)
// 	// 	client.send(Bytes(reterr))
// 	case ERR_FRIEND_ADD_NEED_REQ:
// 		code := gDataManager.addFriendReq(client.uid, fuid, string(group))
// 		if code == ERR_NONE {
// 			req := new(MsgFriendReq)
// 			req.MsgId = MsgId_FriendReq
// 			req.Fuid = client.uid
// 			gDataManager.sendMsgToUser(fuid, Bytes(req))
// 		}

// 		// if code != ERR_NONE {
// 		// 	reterr := NewErrorMsg(ERR_REDIS, MsgId_ReqFriendAdd)
// 		// 	client.send(Bytes(reterr))
// 		// 	return
// 		// }

// 		//send req success msg to client
// 		retmsg := new(MsgFriendReqResult)
// 		retmsg.MsgId = MsgId_FriendReqResult
// 		retmsg.Result = uint16(code)
// 		retmsg.Fuid = fuid
// 		client.send(Bytes(retmsg))
// 	// case ERR_NONE:
// 	// 	//send add success msg to client
// 	// 	retmsg := new(MsgReqFriendAddSuccess)
// 	// 	retmsg.MsgId = MsgId_ReqFriendAddSuccess
// 	// 	client.send(Bytes(retmsg))
// 	default:
// 		retmsg := new(MsgRetFriendAdd)
// 		retmsg.Result = uint16(ret)
// 		retmsg.MsgId = MsgId_RetFriendAdd
// 		retmsg.Fuid = fuid
// 		retmsg.Group = group
// 		client.send(Bytes(retmsg))

// 		req := new(MsgFriendReqAgree)
// 		req.MsgId = MsgId_FriendReqAgree
// 		req.Fuid = client.uid
// 		group, code := gDataManager.getGroupOfFriend(fuid, client.uid)
// 		if code == ERR_NONE {
// 			req.Group = []byte(group)
// 			gDataManager.sendMsgToUser(fuid, Bytes(req))
// 		}
// 		// code := gDataManager.sendMsgToUser(fuid, Bytes(req))

// 		// if code != ERR_NONE {
// 		// 	reterr := NewErrorMsg(ERR_REDIS, MsgId_ReqFriendAdd)
// 		// 	client.send(Bytes(reterr))
// 		// 	return
// 		// }
// 	}
// }

// func OnReqFriendDel(client *ChatClient, data []byte) {
// 	fuid := Uint64(data[2:])
// 	ret := gDataManager.deleteFriend(client.uid, fuid)

// 	retmsg := new(MsgRetFriendDel)
// 	retmsg.MsgId = MsgId_RetFriendDel
// 	retmsg.Result = uint16(ret)
// 	retmsg.Fuid = fuid
// 	client.send(Bytes(retmsg))
// }

// func OnReqUserToBlack(client *ChatClient, data []byte) {
// 	fuid := Uint64(data)
// 	ret := gDataManager.addUserToBlacklist(client.uid, fuid)

// 	retmsg := new(MsgRetUserToBlack)
// 	retmsg.MsgId = MsgId_RetUserToBlack
// 	retmsg.Result = uint16(ret)
// 	retmsg.Fuid = fuid
// 	client.send(Bytes(retmsg))
// }

// func OnReqRemoveUserInBlack(client *ChatClient, data []byte) {
// 	fuid := Uint64(data)
// 	ret := gDataManager.removeUserInBlacklist(client.uid, fuid)

// 	retmsg := new(MsgRetRemoveUserInBlack)
// 	retmsg.MsgId = MsgId_RetRemoveUserInBlack
// 	retmsg.Result = uint16(ret)
// 	retmsg.Fuid = fuid
// 	client.send(Bytes(retmsg))
// }

// func OnReqMoveFriendToGroup(client *ChatClient, data []byte) {
// 	fuid := Uint64(data[2:])
// 	group := data[10:]
// 	ret := gDataManager.moveFriendToGroup(client.uid, fuid, string(group))

// 	retmsg := new(MsgRetUserToBlack)
// 	retmsg.MsgId = MsgId_RetMoveFriendToGroup
// 	retmsg.Result = uint16(ret)
// 	retmsg.Fuid = fuid
// 	client.send(Bytes(retmsg))
// }

// func OnReqSetFriendVerifyType(client *ChatClient, data []byte) {
// 	typ := data[2]
// 	ret := gDataManager.setFriendVerifyType(client.uid, typ)

// 	retmsg := new(MsgRetSetFriendVerifyType)
// 	retmsg.MsgId = MsgId_RetSetFriendVerifyType
// 	retmsg.Result = uint16(ret)
// 	client.send(Bytes(retmsg))
// }

// func OnMessage(client *ChatClient, data []byte) {
// 	fuid := Uint64(data[2:])
// 	//msg := data[10:]

// 	ret := gDataManager.sendMsgToUser(fuid, data)

// 	retmsg := new(MsgRetMessage)
// 	retmsg.MsgId = MsgId_RetMessage
// 	retmsg.Result = uint16(ret)
// 	retmsg.Fuid = fuid
// 	client.send(Bytes(retmsg))
// }
