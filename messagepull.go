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
		datamap, err := gtdb.Manager().PullOnlineMessage(config.ServerAddr, 0)

		if err != nil {
			//fmt.Println(err.Error())
			time.Sleep(time.Duration(5) * time.Second)
			continue
		}

		for _, datastr := range datamap {
			data := []byte(datastr)
			id := Uint64(data[0:8])
			fmt.Println("transfer msg to ", id, " data ", string(data[8:]))
			SessMgr().SendMsgToId(id, data[8:])
		}
	}
}
