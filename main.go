package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"
	//. "github.com/gtechx/Chat/common"
	//"github.com/gtechx/base/gtnet"

	"github.com/gtechx/base/gtnet"
	"github.com/gtechx/chatserver/config"
	"github.com/gtechx/chatserver/data"
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
	server := gtnet.NewServer(config.ServerNet, config.ServerAddr, onNewConn)
	err = server.Start()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer server.Stop()

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

func onNewConn(conn net.Conn) {
	//EntityManager().CreateNullEntity(conn)
	fmt.Println("new conn:", conn.RemoteAddr().String())
	isok := false
	time.AfterFunc(5*time.Second, func() {
		if !isok {
			conn.Close()
		}
	})

	typebuff := make([]byte, 1)
	idbuff := make([]byte, 2)
	sizebuff := make([]byte, 2)
	msgidbuff := make([]byte, 2)

	_, err := conn.Read(typebuff)
	if err != nil {
		fmt.Println(err.Error())
		conn.Close()
		return
	}

	fmt.Println("data type:", typebuff[0])

	_, err = conn.Read(idbuff)
	if err != nil {
		fmt.Println(err.Error())
		conn.Close()
		return
	}
	id := Int(idbuff)

	fmt.Println("id:", id)

	_, err = conn.Read(sizebuff)
	if err != nil {
		fmt.Println(err.Error())
		conn.Close()
		return
	}
	size := Int(sizebuff)

	fmt.Println("data size:", size)

	_, err = conn.Read(msgidbuff)
	if err != nil {
		fmt.Println(err.Error())
		conn.Close()
		return
	}
	msgid := Uint16(msgidbuff)

	fmt.Println("msgid:", msgid)

	databuff := make([]byte, size)

	_, err = conn.Read(databuff)
	if err != nil {
		fmt.Println(err.Error())
		conn.Close()
		return
	}

	// fmt.Println("recv data:", String(databuff), " from "+conn.RemoteAddr().String())

	// if String(databuff) != "wyq" {
	// 	conn.Close()
	// 	return
	// }
	isok = true
	errorcode, ret := HandlerReqLogin(databuff)

	senddata := packageMsg(RetFrame, id, MsgRetLogin, ret)
	_, err = conn.Write(senddata)

	if err != nil {
		conn.Close()
		return
	}

	if errorcode == ERR_NONE {
		SessMgr().CreateSess(conn)
	}
}

func packageMsg(msgtype uint8, id uint16, msgid uint16, data interface{}) []byte {
	ret := []byte{}
	databuff := Bytes(data)
	datalen := uint16(len(databuff))
	ret = append(ret, byte(msgtype))
	ret = append(ret, Bytes(id)...)
	ret = append(ret, Bytes(datalen)...)
	ret = append(ret, Bytes(msgid)...)

	if datalen > 0 {
		ret = append(ret, databuff...)
	}
}

//first, login with account,appname and zonename
//server will return all appdataid in the zone of app
//client need to use one of the appdataid to enter chat.

//before receive chat server chat msg, client need send ready msg to server.
//账号登录的时候发送账号、密码,返回登录成功的token
//登录聊天有两种情况
//1.聊天APP应用，没有分区
//2.游戏带分区聊天应用
//登录聊天的时候需要发送账号、密码，返回appdataidlist
//进入聊天发送appdataid, 服务器根据appdataid创建session
//客户端发送可以接受消息命令，服务器设置玩家在线
