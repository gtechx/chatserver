package main

import (
	"time"

	"github.com/gtechx/chatserver/db"
)

func serverMonitorInit() {
	go startServerMonitor()
}

func startServerMonitor() {
	timer := time.NewTimer(time.Second * 30)

	select {
	case <-timer.C:
		gtdb.Manager().CheckServerTTL()
	}
}
