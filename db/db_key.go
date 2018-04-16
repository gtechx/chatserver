package gtdb

import (
	"strings"
	"time"

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
	// Uid      uint64
	// Appid    uint64
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

func (datakey *DataKey) Init(appname, zonename, account string) {
	datakey.Appname = appname
	datakey.Zonename = zonename
	datakey.Account = account
	// datakey.Uid = uid
	// datakey.Appid = appid

	datakey.Update()
}

func (datakey *DataKey) SetAccount(appname, zonename, account string) {
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

type Admin struct {
	Account             string `redis:"account" json:"account" gorm:"primary_key"`
	Allowsetting        bool   `redis:"setting" json:"setting" gorm:"tinyint(1)"`
	Allowuseradd        bool   `redis:"useradd" json:"useradd" gorm:"tinyint(1)"`
	Allowuserdel        bool   `redis:"userdel" json:"userdel" gorm:"tinyint(1)"`
	Allowmonitoronline  bool   `redis:"monitoronline" json:"monitoronline" gorm:"tinyint(1)"`
	Allowappadd         bool   `redis:"appadd" json:"appadd" gorm:"tinyint(1)"`
	Allowappdel         bool   `redis:"appdel" json:"appdel" gorm:"tinyint(1)"`
	Allowzoneadd        bool   `redis:"zoneadd" json:"zoneadd" gorm:"tinyint(1)"`
	Allowzonedel        bool   `redis:"zonedel" json:"zonedel" gorm:"tinyint(1)"`
	Allowmonitorappdata bool   `redis:"monitorappdata" json:"monitorappdata" gorm:"tinyint(1)"`
}

type Account struct {
	Account   string    `redis:"account" json:"account" gorm:"primary_key"`
	Password  string    `redis:"password" json:"password" gorm:"not null"`
	Email     string    `redis:"email" json:"email"`
	Salt      string    `redis:"salt" json:"salt" gorm:"type:varchar(6);not null;default:''"`
	Regip     string    `redis:"regip" json:"regip"`
	CreatedAt time.Time `redis:"createdate" json:"createdate"`
}

type App struct {
	Name      string    `redis:"name" json:"name" gorm:"primary_key"`
	Owner     string    `redis:"owner" json:"owner"`
	Desc      string    `redis:"desc" json:"desc"`
	Share     string    `redis:"share" json:"share"`
	CreatedAt time.Time `redis:"createdate" json:"createdate"`
}

type AppZone struct {
	Name  string `redis:"name" json:"name"`
	Owner string `redis:"owner" json:"owner"`
}

type AppShare struct {
	Name      string `redis:"name" json:"name"`
	Othername string `redis:"othername" json:"othername"`
}

type AppData struct {
	ID        uint64    `redis:"id" json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Account   string    `redis:"account" json:"account"`
	Appname   string    `redis:"appname" json:"appname"`
	Zonename  string    `redis:"zonename" json:"zonename"`
	Nickname  string    `redis:"nickname" json:"nickname"`
	Desc      string    `redis:"desc" json:"desc"`
	Sex       string    `redis:"sex" json:"sex"`
	Birthday  string    `redis:"birthday" json:"birthday"`
	Country   string    `redis:"country" json:"country"`
	Lastip    string    `redis:"lastip" json:"lastip"`
	Lastlogin time.Time `redis:"lastlogin" json:"lastlogin"`
	CreatedAt time.Time `redis:"createdate" json:"createdate"`
}

type Online struct {
	Dataid uint64 `redis:"dataid" json:"dataid"`
	// Account    string    `redis:"account" json:"account"`
	// Appname    string    `redis:"appname" json:"appname"`
	// Zonename   string    `redis:"zonename" json:"zonename"`
	Serveraddr string    `redis:"serveraddr" json:"serveraddr"`
	State      string    `redis:"state" json:"state"`
	CreatedAt  time.Time `redis:"createdate" json:"createdate"`
}

type Friend struct {
	Dataid      uint64 `redis:"dataid" json:"dataid"`
	Otherdataid uint64 `redis:"otherdataid" json:"otherdataid"`
	// Account      string    `redis:"account" json:"account"`
	// Otheraccount string    `redis:"otheraccount" json:"otheraccount"`
	// Appname      string    `redis:"appname" json:"appname"`
	// Zonename     string    `redis:"zonename" json:"zonename"`
	Group     string    `redis:"group" json:"group"`
	Comment   string    `redis:"comment" json:"comment"`
	CreatedAt time.Time `redis:"createdate" json:"createdate"`
}

type Black struct {
	Dataid      uint64 `redis:"dataid" json:"dataid"`
	Otherdataid uint64 `redis:"otherdataid" json:"otherdataid"`
	// Account      string    `redis:"account" json:"account"`
	// Otheraccount string    `redis:"otheraccount" json:"otheraccount"`
	// Appname      string    `redis:"appname" json:"appname"`
	// Zonename     string    `redis:"zonename" json:"zonename"`
	CreatedAt time.Time `redis:"createdate" json:"createdate"`
}

type Group struct {
	Name        string `redis:"name" json:"name"`
	Dataid      uint64 `redis:"dataid" json:"dataid"`
	Otherdataid uint64 `redis:"otherdataid" json:"otherdataid"`
	// Account  string `redis:"account" json:"account"`
	// Appname  string `redis:"appname" json:"appname"`
	// Zonename string `redis:"zonename" json:"zonename"`
}

var db_tables []interface{} = []interface{}{
	&Admin{},
	&Account{},
	&App{},
	&AppZone{},
	&AppShare{},
	&AppData{},
	&Online{},
	&Friend{},
	&Black{},
	&Group{},
}
