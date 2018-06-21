package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"

	. "github.com/gtechx/base/common"
	//"github.com/gtechx/base/gtnet"

	"github.com/gtechx/base/gtnet"
	"github.com/gtechx/chatserver/config"
	"github.com/gtechx/chatserver/db"
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

	defer gtdb.Manager().UnInitialize()
	err := gtdb.Manager().InitializeRedis(config.RedisAddr, config.RedisPassword, config.RedisDefaultDB)
	if err != nil {
		println("InitializeRedis err:", err.Error())
		return
	}

	err = gtdb.Manager().InitializeMysql(config.MysqlAddr, config.MysqlUserPassword, config.MysqlDefaultDB, config.MysqlTablePrefix)
	if err != nil {
		println("InitializeMysql err:", err.Error())
		return
	}

	//EntityManager().Initialize()

	//register server
	err = gtdb.Manager().RegisterServer(config.ServerAddr)

	if err != nil {
		fmt.Println("register server to gtdata.Manager err:", err)
		return
	}
	defer gtdb.Manager().UnRegisterServer(config.ServerAddr)

	//init loadbalance
	loadBanlanceInit()

	server := gtnet.NewServer()
	err = server.Start(config.ServerNet, config.ServerAddr, onNewConn)
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

	fmt.Println(config.ServerNet + " server start on addr " + config.ServerAddr + " ok...")

	<-quit

	//chatServerStop()
	//gtdb.Manager().UnRegisterServer(config.ServerAddr)
	//gtdata.Manager().UnInitialize()
	//EntityManager().CleanOnlineUsers()
}

func onNewConn(conn net.Conn) {
	//EntityManager().CreateNullEntity(conn)
	fmt.Println("new conn:", conn.RemoteAddr().String())
	isok := false
	defer conn.Close()
	time.AfterFunc(15*time.Second, func() {
		if !isok {
			conn.Close()
		}
	})

	msgtype, id, size, msgid, databuff, err := readMsgHeader(conn)
	isok = true
	fmt.Println(msgtype, id, size, msgid)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if msgid == MsgId_ReqLogin {
		//account login
		_, ret := HandlerReqLogin(databuff)

		senddata := packageMsg(RetFrame, id, MsgId_ReqLogin, ret)
		_, err = conn.Write(senddata)

		if err != nil {
			return
		}
	} else if msgid == MsgId_ReqChatLogin {
		fmt.Println(len(databuff))
		//chat login
		buff := databuff
		slen := int(buff[0])
		account := String(buff[1 : 1+slen])
		buff = buff[1+slen:]
		slen = int(buff[0])
		password := String(buff[1 : 1+slen])
		buff = buff[1+slen:]
		slen = int(buff[0])
		appname := String(buff[1 : 1+slen])
		buff = buff[1+slen:]
		slen = int(buff[0])
		zonename := String(buff[1 : 1+slen])

		fmt.Println(account, password, appname, zonename)
		_, ret := HandlerReqChatLogin(account, password, appname, zonename)

		senddata := packageMsg(RetFrame, id, MsgId_ReqChatLogin, ret)
		_, err = conn.Write(senddata)

		if err != nil {
			return
		}

	waitenterchat:
		msgtype, id, size, msgid, databuff, err = readMsgHeader(conn)
		//errcode := processEnterChat(conn)
		fmt.Println(msgtype, id, size, msgid)
		if err == nil && msgid == MsgId_ReqEnterChat {
			appdataid := Uint64(databuff)
			defer HandlerReqQuitChat(appdataid)
			errcode, ret := HandlerReqEnterChat(appdataid)
			senddata := packageMsg(RetFrame, id, MsgId_ReqEnterChat, ret)
			_, err = conn.Write(senddata)

			if err != nil {
				return
			}

			if errcode == ERR_NONE {
				fmt.Println("sess start:", appdataid)
				lastremoteaddr := conn.RemoteAddr().String()
				lasttime := time.Now()
				sess := SessMgr().CreateSess(conn, appname, zonename, account, appdataid)
				sess.Start()
				gtdb.Manager().UpdateLastLoginInfo(appdataid, lastremoteaddr, lasttime)
			}
		} else if err == nil && msgid == MsgId_ReqCreateAppdata {
			nickname := String(databuff)

			_, ret := HandlerReqCreateAppdata(appname, zonename, account, nickname, conn.RemoteAddr().String())
			senddata := packageMsg(RetFrame, id, MsgId_ReqCreateAppdata, ret)
			_, err = conn.Write(senddata)

			if err != nil {
				return
			}
			goto waitenterchat
		} else if err == nil && msgtype == TickFrame {
			goto waitenterchat
		}
	}
	fmt.Println("conn end")
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
	return ret
}

func readMsgHeader(conn net.Conn) (byte, uint16, uint16, uint16, []byte, error) {
	typebuff := make([]byte, 1)
	idbuff := make([]byte, 2)
	sizebuff := make([]byte, 2)
	msgidbuff := make([]byte, 2)
	var id uint16
	var size uint16
	var msgid uint16
	var databuff []byte

	_, err := conn.Read(typebuff)
	if err != nil {
		fmt.Println(err.Error())
		goto end
	}

	fmt.Println("data type:", typebuff[0])

	if typebuff[0] == TickFrame {
		goto end
	}

	_, err = conn.Read(idbuff)
	if err != nil {
		fmt.Println(err.Error())
		goto end
	}
	id = Uint16(idbuff)

	fmt.Println("id:", id)

	_, err = conn.Read(sizebuff)
	if err != nil {
		fmt.Println(err.Error())
		goto end
	}
	size = Uint16(sizebuff)

	fmt.Println("data size:", size)

	_, err = conn.Read(msgidbuff)
	if err != nil {
		fmt.Println(err.Error())
		goto end
	}
	msgid = Uint16(msgidbuff)

	fmt.Println("msgid:", msgid)

	if size == 0 {
		goto end
	}

	databuff = make([]byte, size)

	_, err = conn.Read(databuff)
	if err != nil {
		fmt.Println(err.Error())
		goto end
	}
end:
	return typebuff[0], id, size, msgid, databuff, err
}

// func processEnterChat(conn net.Conn) uint16 {
// 	isok := false
// 	time.AfterFunc(15*time.Second, func() {
// 		if !isok {
// 			conn.Close()
// 		}
// 	})

// 	msgtype, id, size, msgid, databuff, err := readMsgHeader(conn)

// 	if err != nil {
// 		return ERR_UNKNOWN
// 	}

// 	isok = true

// 	if msgid == MsgId_ReqEnterChat {
// 		appdataid := Uint64(databuff)
// 		errcode, ret := HandlerReqEnterChat(appdataid)
// 		senddata := packageMsg(RetFrame, id, MsgId_ReqEnterChat, ret)
// 		_, err = conn.Write(senddata)

// 		if err != nil {
// 			conn.Close()
// 			return ERR_UNKNOWN
// 		}
// 		return errcode
// 	} else {
// 		conn.Close()
// 		return ERR_MSG_INVALID
// 	}
// }

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
