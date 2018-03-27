package main

import (
	"fmt"
	"time"

	. "github.com/nature19862001/base/common"
	"github.com/nature19862001/base/gtnet"
	"github.com/nature19862001/base/php"
)

const (
	SMALL_MSG_ID_LOGIN uint8 = iota
	//SMALL_MSG_ID_APP_LOGIN
	SMALL_MSG_ID_TOKEN_LOGIN
	SMALL_MSG_ID_TOKEN_REQ
)

type NullEntity struct {
	id    uint64
	uid   uint64
	appid uint64
	zone  uint32
	conn  gtnet.IConn

	recvChan chan []byte
}

func newNullEntity(id uint64, conn gtnet.IConn) *NullEntity {
	return &NullEntity{id: id, conn: conn}
}

func (this *NullEntity) ID() uint64 {
	return this.id
}

func (this *NullEntity) UID() uint64 {
	return this.uid
}

func (this *NullEntity) APPID() uint64 {
	return this.appid
}

func (this *NullEntity) ZONE() uint32 {
	return this.zone
}

func (this *NullEntity) Conn() gtnet.IConn {
	return this.conn
}

func (this *NullEntity) ForceOffline() {
}

func (this *NullEntity) RPC(firstmsgid uint8, secondmsgid uint8, params ...interface{}) {
	buff := []byte{}
	buff = append(buff, Bytes(firstmsgid)...)
	buff = append(buff, Bytes(secondmsgid)...)

	for _, param := range params {
		data := Bytes(param)
		buff = append(buff, Bytes(uint8(len(data)))...)
		buff = append(buff, Bytes(param)...)
	}

	this.conn.Send(append(Bytes(int16(len(buff))), buff...))
}

func (this *NullEntity) start() {
	this.recvChan = make(chan []byte)
	this.conn.SetMsgParser(this)
	this.conn.SetListener(this)

	go this.startProcess()
}

func (this *NullEntity) startProcess() {
	timer := time.NewTimer(time.Second * 30)

	select {
	case <-timer.C:
		this.conn.Close()
	case data := <-this.recvChan:
		this.process(data)
	}

	this.conn = nil
}

func (this *NullEntity) process(data []byte) bool {
	//bigmsgid := Uint8(data)
	smallmsgid := Uint8(data[1:])
	//msgProcesser[bigmsgid][smallmsgid](this, data[2:])
	result := false

	switch smallmsgid {
	case SMALL_MSG_ID_LOGIN:
		uid := Uint64(data[2:10])
		appid := Uint64(data[10:18])
		zone := Uint32(data[18:22])
		password := string(data[22:])

		upass, err := DataManager().GetPassword(uid)

		if err != nil {
			return false
		}

		if upass != password {
			return false
		}

		flag, err := DataManager().IsAppExists(appid)

		if err != nil {
			return false
		}

		if !flag {
			return false
		}

		apptype, err := DataManager().GetAppType(appid)

		if err != nil {
			return false
		}

		if apptype == "game" {
			flag, err = DataManager().IsAppZone(appid, zone)

			if err != nil {
				return false
			}

			if !flag {
				return false
			}
			this.zone = zone
		}

		this.uid = uid
		this.appid = appid
		EntityManager().CreateEntity(TYPE_USER, this)
		fmt.Println("addr:" + this.conn.ConnAddr() + " logined success")

		this.RPC(BIG_MSG_ID_USER, SMALL_MSG_ID_LOGIN_SUCCESS)
		return true
		// case SMALL_MSG_ID_APP_LOGIN:
		// 	appname := string(data[2:34])
		// 	password := string(data[34:])
		// 	//check app login
		// 	//if login ok, then wait for app server verify
		// 	code := gDataManager.checkAppLogin(appname, password)

		// 	if code == ERR_NONE {
		// 		code = gDataManager.setAppOnline(appname)
		// 		if code == ERR_NONE {
		// 			//newAppClient(appname, this.conn)
		// 			fmt.Println("addr:" + this.conn.ConnAddr() + " app logined success")
		// 		}
		// 	}
		// 	result = code != ERR_NONE
		// 	this.RPC(BIG_MSG_ID_LOGIN, SMALL_MSG_ID_APP_LOGIN, uint16(code))
		// case SMALL_MSG_ID_TOKEN_LOGIN:
		// 	token := data[2:]
		// 	str := Authcode(string(token))
		// 	pos := strings.Index(str, ":")

		// 	code := ERR_NONE
		// 	timestamp := Int64(str[:pos])

		// 	if time.Now().Unix()-timestamp > 3600 {
		// 		code = ERR_TIME_OUT
		// 		result = true
		// 	} else {
		// 		uid := Uint64(str[pos:])
		// 		//newChatClient(uid, this.conn)
		// 		fmt.Println("addr:" + this.conn.ConnAddr() + " logined with token success")
		// 	}
		// 	this.RPC(BIG_MSG_ID_LOGIN, SMALL_MSG_ID_TOKEN_LOGIN, uint16(code))
		// case SMALL_MSG_ID_TOKEN_REQ:
		// 	uid := Uint64(data[2:10])
		// 	password := string(data[10:])

		// 	code := gDataManager.checkLogin(uid, password)

		// 	token := ""
		// 	if code == ERR_NONE {
		// 		token = Authcode(String(time.Now().Unix())+":"+String(uid), "ENCODE")
		// 		fmt.Println("uid:" + String(uid) + " get token success")
		// 	}
		// 	result = true
		// 	this.RPC(BIG_MSG_ID_LOGIN, SMALL_MSG_ID_TOKEN_REQ, uint16(code), uid, token)
	}

	return result
}

func (this *NullEntity) ParseHeader(data []byte) int {
	size := Int(data)
	//fmt.Println("header size :", size)
	//p.conn.Send(data)
	return size
}

func (this *NullEntity) ParseMsg(data []byte) {
	newdata := make([]byte, len(data))
	copy(newdata, data)
	this.recvChan <- newdata
}

func (this *NullEntity) OnError(errorcode int, msg string) {
	//fmt.Println("tcpserver error, errorcode:", errorcode, "msg:", msg)
}

func (this *NullEntity) OnPreSend([]byte) {

}

func (this *NullEntity) OnPostSend([]byte, int) {
	// if this.state == state_logouted {
	// 	this.Close()
	// }
}

func (this *NullEntity) OnClose() {
	//fmt.Println("tcpserver closed:", this.clientAddr)
	//this.Close()
}

func (this *NullEntity) OnRecvBusy([]byte) {
	//str := "server is busy"
	//p.conn.Send(Bytes(int16(len(str))))
	//this.conn.Send(append(Bytes(int16(len(str))), []byte(str)...))
}

func (this *NullEntity) OnSendBusy([]byte) {
	// str := "server is busy"
	// p.conn.Send(Bytes(int16(len(str))))
	// p.conn.Send([]byte(str))
}

var UC_KEY string = "1111aaaa"

func Authcode(str string, args ...interface{}) string {
	operation := "DECODE"
	key := ""
	var expiry int64 = 0

	argc := len(args)
	if argc >= 3 {
		texpiry, ok := args[2].(int64)

		if ok {
			expiry = texpiry
		}

		ttexpiry, ok := args[2].(int)

		if ok {
			expiry = int64(ttexpiry)
		}
	}

	if argc >= 2 {
		tkey, ok := args[1].(string)

		if ok {
			key = tkey
		}
	}

	if argc >= 1 {
		toperation, ok := args[0].(string)

		if ok {
			operation = toperation
		}
	}

	ckey_length := 4

	if key == "" {
		key = php.Md5(UC_KEY)
	} else {
		key = php.Md5(key)
	}
	//key = php.Md5(key ? key : UC_KEY)
	keya := php.Md5(php.Substr(key, 0, 16))
	keyb := php.Md5(php.Substr(key, 16, 16))

	keyc := ""
	if ckey_length != 0 {
		if operation == "DECODE" {
			keyc = php.Substr(str, 0, ckey_length)
		} else {
			keyc = php.Substr(php.Md5(String(php.Microtime())), -ckey_length)
		}
	}

	cryptkey := keya + php.Md5(keya+keyc)
	key_length := len(cryptkey)

	if operation == "DECODE" {
		str = php.Base64_decode(php.Substr(str, ckey_length))
	} else {
		var rexpiry int64
		if expiry == 0 {
			rexpiry = 0
		} else {
			rexpiry = expiry + php.Time()
		}
		str1 := php.Sprintf("%010d", rexpiry)
		str = str1 + php.Substr(php.Md5(str+keyb), 0, 16) + str
	}
	string_length := len(str)

	result := ""
	box := php.Range(0, 255, 1)

	j := 0
	i := 0
	a := 0
	rndkey := make([]int, 256)
	for i = 0; i <= 255; i++ {
		rndkey[i] = php.Ord(string(cryptkey[i%key_length]))
	}

	j = 0
	i = 0
	for ; i < 256; i++ {
		j = (j + box[i] + rndkey[i]) % 256
		tmp := box[i]
		box[i] = box[j]
		box[j] = tmp
	}

	j = 0
	i = 0
	for ; i < string_length; i++ {
		a = (a + 1) % 256
		j = (j + box[a]) % 256
		tmp := box[a]
		box[a] = box[j]
		box[j] = tmp
		ntmp := box[(box[a]+box[j])%256]
		//nstr := string(rune(str[i])))
		ndata := int(str[i])
		nres := ndata ^ ntmp
		nstr := php.Chr(nres)
		result = result + nstr //php.Chr(nres)
	}

	if operation == "DECODE" {
		num := Int64(php.Substr(result, 0, 10)) //utils.StrToInt64(string(byteresult[:10])) //php.Substr(result, 0, 10))
		if (num == 0 || num-php.Time() > 0) && php.Substr(result, 10, 16) == php.Substr(php.Md5(php.Substr(result, 26)+keyb), 0, 16) {
			return php.Substr(result, 26) //string(byteresult[26:]) //php.Substr(result, 26)
		} else {
			return ""
		}
	} else {
		return keyc + php.Base64_encode(result) //php.Str_replace("=", "", php.Base64_encode(result)) //base64.StdEncoding.EncodeToString(byteresult)) //php.Base64_encode(result))
	}
}
