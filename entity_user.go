package main

// import (
// 	"fmt"
// 	"strings"
// 	"time"

// 	. "github.com/gtechx/base/common"
// 	"github.com/gtechx/base/gtnet"
// 	"github.com/gtechx/chatserver/data"
// 	"github.com/gtechx/chatserver/entity"
// )

// type UserEntity struct {
// 	id    uint64
// 	uid   uint64
// 	appid uint64
// 	zone  uint32
// 	conn  gtnet.IConn

// 	recvChan chan []byte
// 	quitChan chan int

// 	*gtdata.DataKey
// }

// func keyJoin(params ...interface{}) string {
// 	var builder strings.Builder
// 	count := len(params)
// 	for i := 0; i < count; i++ {
// 		param := params[i]
// 		builder.WriteString(String(param))
// 		if i != (count - 1) {
// 			builder.WriteString(":")
// 		}
// 	}
// 	return builder.String()
// }

// func newUserEntity(entity gtentity.IEntity) *UserEntity {
// 	newentity := &UserEntity{id: entity.ID(), uid: entity.UID(), appid: entity.APPID(), zone: entity.ZONE(), conn: entity.Conn(), DataKey: &gtdata.DataKey{}}
// 	newentity.init()
// 	newentity.start()
// 	return newentity
// }

// func (this *UserEntity) init() {
// 	this.KeyUID = keyJoin("uid", this.uid)
// 	this.KeyAppData = keyJoin("appdata", this.appid, this.zone, this.uid)
// 	this.KeyGroup = keyJoin("group", this.appid, this.zone, this.uid)
// 	this.KeyFriend = keyJoin("friend", this.appid, this.zone, this.uid)
// 	this.KeyFriendRequest = keyJoin("friend", "request", this.appid, this.zone, this.uid)
// 	this.KeyBlack = keyJoin("black", this.appid, this.zone, this.uid)
// 	this.KeyMessageOffline = keyJoin("message", "offline", this.appid, this.zone, this.uid)
// }

// func (this *UserEntity) ID() uint64 {
// 	return this.id
// }

// func (this *UserEntity) UID() uint64 {
// 	return this.uid
// }

// func (this *UserEntity) APPID() uint64 {
// 	return this.appid
// }

// func (this *UserEntity) ZONE() uint32 {
// 	return this.zone
// }

// func (this *UserEntity) Conn() gtnet.IConn {
// 	return this.conn
// }

// func (this *UserEntity) ForceOffline() {
// }

// func (this *UserEntity) RPC(firstmsgid uint8, secondmsgid uint8, params ...interface{}) {
// 	buff := []byte{}
// 	buff = append(buff, Bytes(firstmsgid)...)
// 	buff = append(buff, Bytes(secondmsgid)...)

// 	for _, param := range params {
// 		data := Bytes(param)
// 		buff = append(buff, Bytes(uint8(len(data)))...) //param len
// 		buff = append(buff, Bytes(param)...)            //param data
// 	}

// 	this.conn.Send(append(Bytes(int16(len(buff))), buff...))
// }

// func (this *UserEntity) stop() {
// 	fmt.Println("UserEntity:" + String(this.uid) + " closed")
// 	if this.conn != nil {
// 		this.conn.SetMsgParser(nil)
// 		this.conn.SetListener(nil)
// 		gtdata.Manager().SetUserOffline(this.DataKey)

// 		this.conn.Close()
// 		this.conn = nil

// 		EntityManager().RemoveEntity(this)
// 	}
// }

// func (this *UserEntity) start() {
// 	this.recvChan = make(chan []byte, 2)
// 	this.quitChan = make(chan int, 1)
// 	this.conn.SetMsgParser(this)
// 	this.conn.SetListener(this)
// 	this.broadcastOnlineMsg()
// 	go this.startProcess()
// }

// func (this *UserEntity) broadcastOnlineMsg() {
// 	err := gtdata.Manager().SetUserOffline(this.DataKey)

// 	if err != nil {
// 		this.RPC(BIG_MSG_ID_ERR, SMALL_MSG_ID_ERR_REDIS)
// 		return
// 	}

// 	grouplist, err := gtdata.Manager().GetGroupList(this.DataKey)

// 	if err != nil {
// 		this.RPC(BIG_MSG_ID_ERR, SMALL_MSG_ID_ERR_REDIS)
// 		return
// 	}

// 	friendlist := []uint64{}

// 	for _, group := range grouplist {
// 		gfriendlist, err := gtdata.Manager().GetFriendList(this.DataKey, group)

// 		if err != nil {
// 			this.RPC(BIG_MSG_ID_ERR, SMALL_MSG_ID_ERR_REDIS)
// 			return
// 		}

// 		friendlist = append(friendlist, gfriendlist...)
// 	}

// 	for _, fuid := range friendlist {
// 		flag, err := gtdata.Manager().IsUserOnline(fuid)
// 		if err != nil {
// 			this.RPC(BIG_MSG_ID_ERR, SMALL_MSG_ID_ERR_REDIS)
// 			return
// 		}

// 		if flag {
// 			//send online message to online friend
// 			this.RPC(BIG_MSG_ID_USER, SMALL_MSG_ID_ONLINE, this.uid)
// 		}
// 	}
// }

// func (this *UserEntity) startProcess() {
// 	timer := time.NewTimer(time.Second * 40)
// 	countTimeOut := 0

// 	for {
// 		select {
// 		case <-this.quitChan:
// 			goto end
// 		case <-timer.C:
// 			fmt.Println("countTimeOut++")
// 			countTimeOut++
// 			if countTimeOut >= 2 {
// 				goto end
// 			}
// 		case data := <-this.recvChan:
// 			result := this.process(data)

// 			if result {
// 				goto end
// 			}

// 			countTimeOut = 0
// 		}
// 		timer.Reset(time.Second * 40)
// 	}
// end:
// 	timer.Stop()
// 	this.stop()
// 	fmt.Println("chat process end")
// }

// func (this *UserEntity) process(data []byte) bool {
// 	bigmsgid := Uint8(data)
// 	smallmsgid := Uint8(data[1:])

// 	if bigmsgid >= uint8(BIG_MSG_ID_COUNT) {
// 		fmt.Println("unknown bigmsgid:", bigmsgid)
// 		return false
// 	}

// 	if smallmsgid >= uint8(len(msgProcesser[bigmsgid])) {
// 		fmt.Println("unknown smallmsgid:", smallmsgid)
// 		return false
// 	}

// 	fn := msgProcesser[bigmsgid][smallmsgid]
// 	fn(this, data[2:])

// 	// switch msgid {
// 	// case MsgId_Tick:
// 	// 	ret := new(MsgTick)
// 	// 	ret.MsgId = MsgId_Tick
// 	// 	this.send(Bytes(ret))
// 	// case MsgId_ReqLoginOut:
// 	// 	ret := new(MsgRetLoginOut)
// 	// 	ret.Result = 1
// 	// 	ret.MsgId = MsgId_ReqRetLoginOut
// 	// 	this.send(Bytes(ret))
// 	// 	return true
// 	// case MsgId_Echo:
// 	// 	// ret := new(Echo)
// 	// 	// ret.MsgId = MsgId_Echo
// 	// 	// ret.Data = data[2:]
// 	// 	this.send(data)
// 	// default:
// 	// 	fn, ok := msgProcesser[msgid]
// 	// 	if ok {
// 	// 		fn(this, data)
// 	// 	} else {
// 	// 		fmt.Println("unknown msgid:", msgid)
// 	// 	}
// 	// }

// 	return true
// }

// // IMsgParser start
// func (this *UserEntity) ParseHeader(data []byte) int {
// 	size := Int(data)
// 	//fmt.Println("header size :", size)
// 	//p.conn.Send(data)
// 	return size
// }

// func (this *UserEntity) ParseMsg(data []byte) {
// 	//fmt.Println("UserEntity:", this.conn.ConnAddr(), "say:", String(data))
// 	newdata := make([]byte, len(data))
// 	copy(newdata, data)
// 	this.recvChan <- newdata
// }

// // IMsgParser end

// // IConnListener start
// func (this *UserEntity) OnError(errorcode int, msg string) {
// 	//fmt.Println("tcpserver error, errorcode:", errorcode, "msg:", msg)
// }

// func (this *UserEntity) OnPreSend([]byte) {

// }

// func (this *UserEntity) OnPostSend([]byte, int) {
// 	// if this.state == state_logouted {
// 	// 	this.Close()
// 	// }
// }

// func (this *UserEntity) OnClose() {
// 	//fmt.Println("tcpserver closed:", this.UserEntityAddr)
// 	this.quitChan <- 1
// }

// func (this *UserEntity) OnRecvBusy([]byte) {
// 	//str := "server is busy"
// 	//p.conn.Send(Bytes(int16(len(str))))
// 	//this.conn.Send(append(Bytes(int16(len(str))), []byte(str)...))
// }

// func (this *UserEntity) OnSendBusy([]byte) {
// 	// str := "server is busy"
// 	// p.conn.Send(Bytes(int16(len(str))))
// 	// p.conn.Send([]byte(str))
// }

// // IConnListener end
