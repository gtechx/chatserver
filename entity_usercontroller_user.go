package main

import (
	. "github.com/gtechx/base/common"
	"github.com/gtechx/chatserver/data"
)

const (
	SMALL_MSG_ID_CREATE_ACCOUNT uint8 = iota
	SMALL_MSG_ID_REGISTE_SUCCESS
	SMALL_MSG_ID_LOGIN_SUCCESS
	SMALL_MSG_ID_TICK
	SMALL_MSG_ID_ECHO
	SMALL_MSG_ID_LOGOUT
	SMALL_MSG_ID_ONLINE
	SMALL_MSG_ID_OFFLINE
	SMALL_MSG_ID_BUSY
	SMALL_MSG_ID_HIDE

	SMALL_MSG_ID_COUNT
)

func init() {
	if msgProcesser == nil {
		msgProcesser = make([][]func(*UserEntity, []byte), BIG_MSG_ID_COUNT)
	}
	msgProcesser[BIG_MSG_ID_USER] = make([]func(*UserEntity, []byte), SMALL_MSG_ID_COUNT)
	msgProcesser[BIG_MSG_ID_USER][SMALL_MSG_ID_CREATE_ACCOUNT] = onCreateAccount
	msgProcesser[BIG_MSG_ID_USER][SMALL_MSG_ID_TICK] = onTick
	msgProcesser[BIG_MSG_ID_USER][SMALL_MSG_ID_ECHO] = onEcho
	msgProcesser[BIG_MSG_ID_USER][SMALL_MSG_ID_LOGOUT] = onLogout
}

func onCreateAccount(entity *UserEntity, data []byte) {
	account := String(data[:32])
	password := string(data[8:])

	flag, err := gtdata.Manager().IsAccountExists(account)

	if err != nil {
		entity.RPC(BIG_MSG_ID_ERR, SMALL_MSG_ID_ERR_REDIS)
		return
	}

	if flag {
		entity.RPC(BIG_MSG_ID_ERR, SMALL_MSG_ID_ERR_CODE, ERR_ACCOUNT_EXISTS)
		return
	}

	err = gtdata.Manager().CreateAccount(account, password, "")

	if err != nil {
		entity.RPC(BIG_MSG_ID_ERR, SMALL_MSG_ID_ERR_REDIS)
		return
	}

	entity.RPC(BIG_MSG_ID_USER, SMALL_MSG_ID_REGISTE_SUCCESS)
}

func onTick(entity *UserEntity, data []byte) {
	entity.RPC(BIG_MSG_ID_USER, SMALL_MSG_ID_TICK)
}

func onEcho(entity *UserEntity, data []byte) {
	entity.RPC(BIG_MSG_ID_USER, SMALL_MSG_ID_TICK, data)
}

func onLogout(entity *UserEntity, data []byte) {
	err := gtdata.Manager().SetUserOffline(entity.EntityKey)

	if err != nil {
		entity.RPC(BIG_MSG_ID_ERR, SMALL_MSG_ID_ERR_REDIS)
		return
	}

	grouplist, err := gtdata.Manager().GetGroupList(entity.EntityKey)

	if err != nil {
		entity.RPC(BIG_MSG_ID_ERR, SMALL_MSG_ID_ERR_REDIS)
		return
	}

	friendlist := []uint64{}

	for _, group := range grouplist {
		gfriendlist, err := gtdata.Manager().GetFriendList(entity.EntityKey, group)

		if err != nil {
			entity.RPC(BIG_MSG_ID_ERR, SMALL_MSG_ID_ERR_REDIS)
			return
		}

		friendlist = append(friendlist, gfriendlist...)
	}

	for _, fuid := range friendlist {
		flag, err := gtdata.Manager().IsUserOnline(fuid)
		if err != nil {
			entity.RPC(BIG_MSG_ID_ERR, SMALL_MSG_ID_ERR_REDIS)
			return
		}

		if flag {
			//send offline message to online friend
			entity.RPC(BIG_MSG_ID_USER, SMALL_MSG_ID_OFFLINE, entity.UID())
		}
	}
}
