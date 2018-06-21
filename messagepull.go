package main

import (
	. "github.com/gtechx/base/common"
	"github.com/gtechx/chatserver/db"
)

func messagePullInit() {
	go startMessagePull()
}

func startMessagePull() {
	for {
		data, err := gtdb.Manager().PullOnlineMessage(serverAddr, 15)

		if err != nil {
			continue
		}

		if data != nil {
			//fmt.Println(data)
			id := Uint64(data[0:8])
			ok := SessMgr().SendMsgToId(id, data[8:])
			if !ok {
				gtdb.Manager().SendMsgToUserOffline(id, data[8:])
			}
		}
	}
}
