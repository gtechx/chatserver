package main

// import (
// 	. "github.com/gtechx/base/common"
// 	//"github.com/gtechx/base/gtnet"
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
