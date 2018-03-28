package main

import (
	//. "github.com/gtechx/Chat/common"
	//"github.com/gtechx/base/gtnet"
	"time"
)

func keepLiveInit() {
	go startServerTTLKeep()
}

func startServerTTLKeep() {
	timer := time.NewTimer(time.Second * 30)

	select {
	case <-timer.C:
		DataManager().SetServerTTL(serverAddr, 60)
	}
}
