package gtdata

import (
	"strings"

	. "github.com/gtechx/base/common"
)

func keyJoin(params ...interface{}) string {
	var builder strings.Builder
	count := len(params)
	for i := 0; i < count; i++ {
		param := params[i]
		builder.WriteString(String(param))
		if i != (count - 1) {
			builder.WriteString(":")
		}
	}
	return builder.String()
}

//Key[Hset|Zset|Set][store data][by][field]
//store data:
//1.User 表示存储的是user表的数据
//2.UidAccount表示存储的是uid account的键值对
//3.By表示根据by后面的field不同，有n条独立的这样的数据.没有by的表示这个key只有一个，一般用来存储统计数据
type DataKey struct {
	KeyUserHsetByAccount      string //hset:user:account:xxx
	KeyUserHsetUidAccount     string //hset:user:uid:account
	KeyUserZsetRegdateAccount string //zset:user:regdate
	KeyUserSet                string //set:user
	//KeyUID            string

	KeyAppSet                         string //set:app
	KeyAppHsetAppidAppname            string //hset:app:appid:appname
	KeyAppHsetByAppname               string //hset:app:appname:xxx
	KeyAppSetAppnameByAccount         string //set:app:account:xxx
	KeyAppZsetRegdateAppnameByAccount string //zset:app:regdate:account:xxx
	KeyAppSetShareByAppname           string //set:app:share:xxx
	KeyAppSetZonenameByAppname        string //set:app:zone:xxx

	KeyAppDataHsetByAppidZonenameAccount                   string //hset:app:data:xxx:xxx:xxx
	KeyAppDataSetGroupByAppidZonenameAccount               string //set:app:data:group:xxx:xxx:xxx
	KeyAppDataHsetFriendByAppidZonenameAccount             string //hset:app:data:friend:xxx:xxx:xxx
	KeyAppDataHsetFriendrequestGroupByAppidZonenameAccount string //hset:app:data:friend:request:xxx:xxx:xxx
	KeyAppDataSetBlackByAppidZonenameAccount               string //set:app:data:black:xxx:xxx:xxx
	KeyAppDataListMsgByAppidZonenameAccount                string //list:app:data:msg:offline:xxx:xxx:xxx

	// KeyAppData        string
	// KeyGroup          string
	// KeyFriend         string
	// KeyFriendRequest  string
	// KeyBlack          string
	// KeyMessageOffline string
	Appname  string
	Zonename string
	Account  string
	Uid      uint64
	Appid    uint64
}

func (datakey *DataKey) Update() {
	datakey.KeyUserHsetByAccount = keyJoin("hset:user:account", datakey.Account)
	datakey.KeyUserHsetUidAccount = "hset:user:uid:account"
	datakey.KeyUserZsetRegdateAccount = "zset:user:regdate"
	datakey.KeyUserSet = "set:user"

	datakey.KeyAppSet = "set:app"
	datakey.KeyAppHsetAppidAppname = "hset:app:appid:appname"
	datakey.KeyAppHsetByAppname = keyJoin("hset:app:appname", datakey.Appname)
	datakey.KeyAppSetAppnameByAccount = keyJoin("set:app:account", datakey.Account)
	datakey.KeyAppZsetRegdateAppnameByAccount = keyJoin("zset:app:regdate:account", datakey.Account)
	datakey.KeyAppSetShareByAppname = keyJoin("set:app:share", datakey.Appname)
	datakey.KeyAppSetZonenameByAppname = keyJoin("set:app:zone", datakey.Appname)

	datakey.KeyAppDataHsetByAppidZonenameAccount = keyJoin("hset:app:data", datakey.Appname, datakey.Zonename, datakey.Account)
	datakey.KeyAppDataSetGroupByAppidZonenameAccount = keyJoin("set:app:data:group", datakey.Appname, datakey.Zonename, datakey.Account)
	datakey.KeyAppDataHsetFriendByAppidZonenameAccount = keyJoin("hset:app:data:friend", datakey.Appname, datakey.Zonename, datakey.Account)
	datakey.KeyAppDataHsetFriendrequestGroupByAppidZonenameAccount = keyJoin("hset:app:data:friend:request", datakey.Appname, datakey.Zonename, datakey.Account)
	datakey.KeyAppDataSetBlackByAppidZonenameAccount = keyJoin("set:app:data:black", datakey.Appname, datakey.Zonename, datakey.Account)
	datakey.KeyAppDataListMsgByAppidZonenameAccount = keyJoin("list:app:data:msg:offline", datakey.Appname, datakey.Zonename, datakey.Account)
}

func (datakey *DataKey) Init(appname, zonename, account string, uid, appid uint64) {
	datakey.Appname = appname
	datakey.Zonename = zonename
	datakey.Account = account
	datakey.Uid = uid
	datakey.Appid = appid

	datakey.Update()
}

func (datakey *DataKey) SetAccount(appname, zonename, account string, uid, appid uint64) {
	datakey.Account = account
	datakey.Update()
}

func (datakey *DataKey) SetAppname(appname string) {
	datakey.Appname = appname
	datakey.Update()
}

func (datakey *DataKey) SetZonename(zonename string) {
	datakey.Zonename = zonename
	datakey.Update()
}

type App struct {
	Appid   uint64 `redis:"appid" json:"appid"`
	Name    string `redis:"name" json:"name"`
	Owner   uint64 `redis:"owner" json:"owner"`
	Desc    string `redis:"desc" json:"desc"`
	Regdate int64  `redis:"regdate" json:"regdate"`
	Type    string `redis:"type" json:"type"`
	Share   uint64 `redis:"share" json:"share"`
}

type AppData struct {
	Id       uint64 `redis:"id" json:"id"`
	Name     string `redis:"name" json:"name"`
	Desc     string `redis:"desc" json:"desc"`
	Regdate  int64  `redis:"regdate" json:"regdate"`
	Sex      string `redis:"sex" json:"sex"`
	Birthday string `redis:"birthday" json:"birthday"`
	Country  string `redis:"country" json:"country"`
}
