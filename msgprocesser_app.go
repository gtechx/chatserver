package main

// import (
// 	"fmt"
// 	. "github.com/nature19862001/Chat/protocol"
// 	. "github.com/nature19862001/base/common"
// 	"strings"
// 	"time"
// )

// func AppOnReqToken(client *AppClient, data []byte) {
// 	uid := Uint64(data[2:10])

// 	ok := gDataManager.isAppUser(client.appName, uid)

// 	ret := new(AppMsgRetToken)
// 	ret.MsgId = AppMsgId_RetToken

// 	code := ERR_NONE
// 	if ok {
// 		token := Authcode(String(time.Now().Unix())+":"+String(uid), "ENCODE")
// 		ret.Token = []byte(token)
// 		fmt.Println("app uid:" + String(uid) + " get token success")
// 	} else {
// 		code = ERR_USER_NOT_EXIST
// 	}
// 	ret.Result = uint16(code)
// 	client.send(Bytes(ret))
// }

// func AppOnReqTokenVerify(client *AppClient, data []byte) {
// 	token := data[2:]
// 	str := Authcode(string(token))
// 	pos := strings.Index(str, ":")

// 	ret := new(AppMsgRetTokenVerify)
// 	ret.MsgId = AppMsgId_RetTokenVerify

// 	code := ERR_NONE
// 	timestamp := Int64(str[:pos])

// 	if time.Now().Unix()-timestamp > 3600 {
// 		code = ERR_TIME_OUT
// 	} else {
// 		uid := Uint64(str[pos:])
// 		ok := gDataManager.isUserExist(uid)
// 		if !ok {
// 			code = ERR_USER_NOT_EXIST
// 		} else {
// 			fmt.Println("app uid:" + String(uid) + " verify token success")
// 		}
// 	}
// 	ret.Result = uint16(code)
// 	client.send(Bytes(ret))
// }
