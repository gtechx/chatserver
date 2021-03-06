package main

import (
	"fmt"
	"time"

	. "github.com/gtechx/base/common"
	"github.com/gtechx/chatserver/config"
	"github.com/gtechx/chatserver/db"
)

func messagePullInit() {
	go startMessagePull()
}

func startMessagePull() {
	for {
		data, err := gtdb.Manager().PullOnlineMessage(config.ServerAddr)

		if err != nil {
			//fmt.Println(err.Error())
			time.Sleep(1 * time.Second)
			continue
		}

		id := Uint64(data[0:8])
		fmt.Println("transfer msg to ", id, " data ", string(data[8:]))
		if !SessMgr().SendMsgToId(id, data[8:]) {
			gtdb.Manager().SendMsgToUserOffline(id, data[8:])
		}
	}
}
