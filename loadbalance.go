package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	. "github.com/gtechx/base/common"
)

func loadBanlanceInit() {
	go startHTTPServer()
	go starUserRegister()
}

func startHTTPServer() {
	http.HandleFunc("/serverlist", getServerList)
	http.ListenAndServe(":9001", nil)
}

func getServerList(rw http.ResponseWriter, req *http.Request) {
	serverlist, _ := DataManager().GetServerList()

	ret := "{\r\n\t\"serverlist\":\r\n\t[\r\n"
	for i := 0; i < len(serverlist); i++ {
		ret += "\t\t{ \"addr\":\"" + serverlist[i] + "\" }"
		if i != len(serverlist)-1 {
			ret += ",\r\n"
		}
	}
	ret += "\r\n\t]\r\n"
	ret += "}\r\n"

	io.WriteString(rw, ret)
}

func starUserRegister() {
	http.HandleFunc("/register", register)
	http.HandleFunc("/create", create)
	http.HandleFunc("/tokenverify", tokenVerify)
	http.ListenAndServe(":8081", nil)
}

type Error struct {
	ErrorMsg  string `json:"error"`
	ErrorCode uint16 `json:"errorcode"`
}

func tokenVerify(rw http.ResponseWriter, req *http.Request) {
	token := req.PostFormValue("token")
	str := Authcode(token)
	pos := strings.Index(str, ":")
	timestamp := Int64(str[:pos])
	uid := Uint64(str[pos:])

	errmsg := new(Error)
	if time.Now().Unix()-timestamp > 3600 {
		errmsg.ErrorCode = ERR_TIME_OUT
		errmsg.ErrorMsg = "ERR_TIME_OUT"
	} else {
		flag, err := DataManager().IsUIDExists(uid)
		if err != nil {
			errmsg.ErrorCode = ERR_REDIS
			errmsg.ErrorMsg = "ERR_REDIS"
		} else if !flag {
			errmsg.ErrorCode = ERR_ACCOUNT_NOT_EXISTS
			errmsg.ErrorMsg = "ERR_ACCOUNT_NOT_EXISTS"
		}
	}

	// if gDataManager.verifyAppLoginData(token, uid) {
	// 	if verifyAppLogin(uid) {
	// 		errmsg.ErrorCode = 0
	// 	} else {
	// 		errmsg.ErrorCode = 2
	// 		errmsg.ErrorMsg = "uid not logined"
	// 	}
	// } else {
	// 	errmsg.ErrorCode = 1
	// 	errmsg.ErrorMsg = "uuid is not exist or uid is not right"
	// }

	data, _ := json.Marshal(&errmsg)
	io.WriteString(rw, string(data))
}

func writeHeader(rw http.ResponseWriter) {
	ret := "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Transitional//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd\">"
	ret += "<html xmlns=\"http://www.w3.org/1999/xhtml\">"
	ret += "<head>"
	ret += "<meta http-equiv=\"Content-Type\" content=\"text/html; charset=utf-8\" />"
	ret += "<title>register</title>"
	ret += "<meta content=\"GTech Inc.\" name=\"Copyright\" />"
	ret += "<script src=\"http://cdn.bootcss.com/blueimp-md5/1.1.0/js/md5.min.js\"></script>"
	ret += "</head>"
	ret += "<body>"

	io.WriteString(rw, ret)
}

func writeFooter(rw http.ResponseWriter) {
	ret := "</body>"
	ret += "</html>"

	io.WriteString(rw, ret)
}

func register(rw http.ResponseWriter, req *http.Request) {
	//req.ParseForm()
	fmt.Println(req.Method)
	fmt.Println(req.RemoteAddr)
	// fmt.Println(req.PostForm)
	// fmt.Println(req.Form["username"])
	// fmt.Println(req.PostForm["username"])
	// fmt.Println(req.PostFormValue("username"))
	// nickname := req.PostFormValue("nickname")
	// password := req.PostFormValue("password")
	// regip := req.RemoteAddr
	// method := req.Method

	writeHeader(rw)

	ret := "<form method=\"post\" action=\"/create\" onsubmit=\"return true;\">"
	ret += "账号：<input type=\"text\" name=\"account\" />"
	ret += "<br/>"
	ret += "密码：<input type=\"password\" name=\"password1\" oninput=\"document.getElementById('password').value = md5(this.value);\" onpropertychange=\"document.getElementById('password').value = md5(this.value);\" />"
	ret += "<input type=\"hidden\" name=\"password\" id=\"password\" />"
	ret += "<br/>"
	ret += "<input type=\"submit\" name=\"login_button\" value=\"提交\">"
	ret += "</form>"

	io.WriteString(rw, ret)

	writeFooter(rw)
}

func create(rw http.ResponseWriter, req *http.Request) {
	var err error
	account := req.PostFormValue("account")
	password := req.PostFormValue("password")
	regip := req.RemoteAddr
	method := req.Method

	writeHeader(rw)

	ret := ""

	if method != "POST" {
		// ret := "{\r\n\terrorcode:1,\r\n"
		// ret = "\r\n\terror:\"need post\",\r\n"
		// ret += "}"
		ret += "<span>请使用post方法!</span><br/>"
		//io.WriteString(rw, ret)
		goto errend
	}

	if account == "" {
		// ret := "{\r\n\terrorcode:1,\r\n"
		// ret = "\r\n\terror:\"need nickname\",\r\n"
		// ret += "}"
		ret += "<span>请输入昵称!</span><br/>"
		//io.WriteString(rw, ret)
		goto errend
	}

	if password == "" {
		// ret := "{\r\n\terrorcode:2,\r\n"
		// ret = "\r\n\terror:\"need password\",\r\n"
		// ret += "}"
		ret += "<span>请输入密码!</span><br/>"
		//io.WriteString(rw, ret)
		goto errend
	}

	err = DataManager().CreateAccount(account, password, regip)

	if err != nil {
		// ret := "{\r\n\terrorcode:3,\r\n"
		// ret = "\r\n\terror:\"server error\",\r\n"
		// ret += "}"
		ret += "<span>注册失败，服务器内部错误!</span><br/>"
		//io.WriteString(rw, ret)
		goto errend
	}

	// ret := "{\r\n\terrorcode:0,\r\n"
	// ret = "\r\n\terror:\"\",\r\n"
	// ret = "\r\n\tuid:" + String(uid) + ",\r\n"
	// ret += "}"
	ret += "<span>注册成功，登录账号：" + account + "</span><br/>"
	goto end
errend:
	ret += "<form method=\"post\" action=\"/create\">"
	ret += "账号：<input type=\"text\" name=\"account\" />"
	ret += "<br/>"
	ret += "密码：<input type=\"password\" name=\"password1\" oninput=\"document.getElementById('password').value = md5(this.value);\" onpropertychange=\"document.getElementById('password').value = md5(this.value);\" />"
	ret += "<input type=\"hidden\" name=\"password\" id=\"password\" />"
	ret += "<br/>"
	ret += "<input type=\"submit\" name=\"login_button\" value=\"提交\">"
	ret += "</form>"
end:
	io.WriteString(rw, ret)

	writeFooter(rw)
}
