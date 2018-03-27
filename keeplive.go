package main

import (
	//. "github.com/nature19862001/Chat/common"
	//"github.com/nature19862001/base/gtnet"
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
