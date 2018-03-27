package main

// import (
// 	. "github.com/nature19862001/base/common"
// 	//"github.com/nature19862001/base/gtnet"
// )

func messagePullInit() {
	go startMessagePull()
}

func startMessagePull() {
	for {
		data, err := DataManager().PullOnlineMessage(serverAddr, 15)

		if err != nil {
			continue
		}

		if data != nil {
			//fmt.Println(data)
			// uid := Uint64(data[0:8])
			// sendMsgToUid(uid, data[8:])
		}
	}
}
