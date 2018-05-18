package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	//. "github.com/gtechx/Chat/common"
	//"github.com/gtechx/base/gtnet"

	"github.com/gtechx/base/gtnet"
	"github.com/gtechx/chatserver/config"
	"github.com/gtechx/chatserver/data"
	"github.com/gtechx/chatserver/service"
)

var quit chan os.Signal

var nettype string = "tcp"
var serverAddr string = "127.0.0.1:9090"
var redisNet string = "tcp"
var redisAddr string = "192.168.93.16:6379"

func main() {
	//var err error
	quit = make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)

	pnet := flag.String("net", "", "-net=")
	paddr := flag.String("addr", "", "-addr=")
	//predisnet := flag.String("redisnet", redisNet, "-redisnet=")
	predisaddr := flag.String("redisaddr", "", "-redisaddr=")

	flag.Parse()

	if pnet != nil && *pnet != "" {
		config.ServerNet = *pnet
	}
	if paddr != nil && *paddr != "" {
		config.ServerAddr = *paddr
	}
	if predisaddr != nil && *predisaddr != "" {
		config.RedisAddr = *predisaddr
	}
	// nettype = *pnet
	// serverAddr = *paddr
	// redisNet = *predisnet
	// redisAddr = *predisaddr

	err := gtdata.Manager().Initialize(config.RedisAddr, config.RedisPassword, config.DefaultDB, config.StartUID, config.StartAPPID)

	if err != nil {
		fmt.Println("Initialize gtdata.Manager err:", err)
		return
	}

	EntityManager().Initialize()

	//register server
	err = gtdata.Manager().RegisterServer(config.ServerAddr)

	if err != nil {
		fmt.Println("register server to gtdata.Manager err:", err)
		return
	}

	//init loadbalance
	loadBanlanceInit()

	//init chat server
	// ok = chatServerStart(nettype, serverAddr)

	// if !ok {
	// 	fmt.Println("chat server init failed!!!")
	// 	return
	// }
	service := gtservice.NewService("chatserver", config.ServerNet, config.ServerAddr, onNewConn)
	err = service.Start()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer service.Stop()

	//keep live init
	keepLiveInit()

	//other server live monitor init
	serverMonitorInit()

	//msg from other server monitor
	messagePullInit()

	fmt.Println("server start ok...")

	<-quit

	//chatServerStop()
	gtdata.Manager().UnRegisterServer(config.ServerAddr)
	gtdata.Manager().UnInitialize()
	EntityManager().CleanOnlineUsers()
}

func onNewConn(conn gtnet.IConn) {
	EntityManager().CreateNullEntity(conn)
}

//first, login with account,appname and zonename
//server will return all appdataid in the zone of app
//client need to use one of the appdataid to enter chat.

//before receive chat server chat msg, client need send ready msg to server.
