package main

import (
	//. "github.com/gtechx/base/common"
	"github.com/gtechx/chatserver/db"
)

func isAppDataExists(id uint64, perrcode *uint16) bool {
	flag, err := gtdb.Manager().IsAppDataExists(id)

	if err != nil {
		*perrcode = ERR_DB
	} else {
		if !flag {
			*perrcode = ERR_APPDATAID_NOT_EXISTS
		} else {
			return true
		}
	}

	return false
}

// func isAppDataNotExists(id uint64, perrcode *uint16) bool {
// 	flag, err := gtdb.Manager().IsAppDataExists(id)

// 	if err != nil {
// 		*perrcode = ERR_DB
// 	} else {
// 		if flag {
// 			*perrcode = ERR_ROOM_EXISTS
// 		} else {
// 			return true
// 		}
// 	}

// 	return false
// }
