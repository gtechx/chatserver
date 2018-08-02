package gtdb

import (
	"fmt"
	"strings"
	"time"

	. "github.com/gtechx/base/common"
	"github.com/jinzhu/gorm"
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

//关于结构体中time类型json化问题，所有终端统一使用RFC3339Nano字符串格式，即go中time默认json的格式。

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
	Account      string    `redis:"account" json:"account" gorm:"primary_key"`
	Adminadmin   bool      `redis:"adminadmin" json:"adminadmin" gorm:"tinyint(1);default:0"`
	Adminuser    bool      `redis:"adminuser" json:"adminuser" gorm:"tinyint(1);default:0"`
	Adminapp     bool      `redis:"adminapp" json:"adminapp" gorm:"tinyint(1);default:0"`
	Adminappdata bool      `redis:"adminappdata" json:"adminappdata" gorm:"tinyint(1);default:0"`
	Adminonline  bool      `redis:"adminonline" json:"adminonline" gorm:"tinyint(1);default:0"`
	Adminmessage bool      `redis:"adminmessage" json:"adminmessage" gorm:"tinyint(1);default:0"`
	Appcount     uint32    `redis:"appcount" json:"appcount" gorm:"default:0"`
	Expire       time.Time `redis:"expire" json:"expire" gorm:"type:datetime"`

	AdminApps []AdminApp `json:"-" gorm:"foreignkey:Adminaccount;association_foreignkey:Account"`
}

// func (admin *Admin) MarshalJSON() ([]byte, error) {
// 	// 定义一个该结构体的别名
// 	type Alias Admin
// 	// 定义一个新的结构体
// 	tmpSt := struct {
// 		Alias
// 		Expire string `json:"expire"`
// 	}{
// 		Alias:  (Alias)(*admin),
// 		Expire: admin.Expire.Format("01/02/2006"),
// 	}
// 	return json.Marshal(tmpSt)
// }

type AdminApp struct {
	Adminaccount string `redis:"adminaccount" json:"adminaccount"`
	Appname      string `redis:"appname" json:"appname"`
}

type Account struct {
	Account   string    `redis:"account" json:"account" gorm:"primary_key"`
	Password  string    `redis:"password" json:"-" gorm:"not null"`
	Email     string    `redis:"email" json:"email"`
	Mobile    string    `redis:"mobile" json:"mobile"`
	Salt      string    `redis:"salt" json:"-" gorm:"type:varchar(6);not null;default:''"`
	Regip     string    `redis:"regip" json:"regip"`
	Isbaned   bool      `redis:"isbaned" json:"isbaned" gorm:"tinyint(1);default:0"`
	CreatedAt time.Time `redis:"createdate" json:"createdate"`

	Apps []App `json:"-" gorm:"foreignkey:Owner;association_foreignkey:Account"`
}

// func (acc *Account) MarshalJSON() ([]byte, error) {
// 	// 定义一个该结构体的别名
// 	type Alias Account
// 	// 定义一个新的结构体
// 	tmpSt := struct {
// 		Alias
// 		CreateDate string `json:"createdate"`
// 	}{
// 		Alias:      (Alias)(*acc),
// 		CreateDate: acc.CreatedAt.Format("01/02/2006 15:04:05"),
// 	}
// 	return json.Marshal(tmpSt)
// }

func (acc *Account) BeforeDelete(tx *gorm.DB) error {
	fmt.Println("BeforeDelete Account", acc)

	var apps []App
	for {
		if err := tx.Model(acc).Limit(100).Related(&apps, "Apps").Error; err != nil {
			return err
		}
		if len(apps) == 0 {
			break
		}
		for _, app := range apps {
			if err := tx.Delete(&app).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

type AccountAdminApp struct {
	Account string `redis:"account" json:"account"`
	Appname string `redis:"appname" json:"appname"`
}

type App struct {
	ID      uint64 `redis:"id" json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Appname string `redis:"appname" json:"appname" gorm:"primary_key"`
	Owner   string `redis:"owner" json:"owner"`
	Desc    string `redis:"desc" json:"desc"`
	Share   string `redis:"share" json:"share"`
	//Channelmax uint8     `redis:"channelmax" json:"channelmax"`
	//Friendmax  uint8     `redis:"friendmax" json:"friendmax"`
	CreatedAt time.Time `redis:"createdate" json:"createdate"`

	AppZones  []AppZone  `json:"-" gorm:"foreignkey:Owner;association_foreignkey:Appname"`
	AppShares []AppShare `json:"-" gorm:"foreignkey:Appname;association_foreignkey:Appname"`
	AppDatas  []AppData  `json:"-" gorm:"foreignkey:Appname;association_foreignkey:Appname"`
}

// func (app *App) MarshalJSON() ([]byte, error) {
// 	// 定义一个该结构体的别名
// 	type Alias App
// 	// 定义一个新的结构体
// 	tmpSt := struct {
// 		Alias
// 		CreateDate string `json:"createdate"`
// 	}{
// 		Alias:      (Alias)(*app),
// 		CreateDate: app.CreatedAt.Format("01/02/2006 15:04:05"),
// 	}
// 	return json.Marshal(tmpSt)
// }

func (app *App) BeforeDelete(tx *gorm.DB) error {
	fmt.Println("BeforeDelete App", app)

	//var zones []AppZone
	//tx.Model(app).Related(&zones, "AppZones")
	//fmt.Println(zones)
	//delete zones of this app
	if err := tx.Delete(&AppZone{}, "owner = ?", app.Appname).Error; err != nil {
		return err
	}

	//delete appshare of this app
	if err := tx.Delete(&AppShare{}, "appname = ? OR othername = ?", app.Appname, app.Appname).Error; err != nil {
		return err
	}

	//delete appdatas of this app
	var appdatas []AppData
	for {
		if err := tx.Model(app).Limit(1000).Related(&appdatas, "AppDatas").Error; err != nil {
			return err
		}
		if len(appdatas) == 0 {
			break
		}
		for _, appdata := range appdatas {
			if err := tx.Delete(&appdata).Error; err != nil {
				return err
			}
		}
	}

	//update share colomn who share with me
	if err := tx.Model(&App{}).Where("share = ?", app.Appname).Update("share", "").Error; err != nil {
		return err
	}

	if err := tx.Delete(&AccountApp{}, "appname = ?", app.Appname).Error; err != nil {
		return err
	}

	if err := tx.Delete(&AccountZone{}, "appname = ?", app.Appname).Error; err != nil {
		return err
	}

	// for _, zone := range zones {
	// 	tx.Delete(&zone, "name = ? AND owner = ?", zone.Name, zone.Owner)
	// }
	return nil
}

func (app *App) AfterDelete(tx *gorm.DB) error {
	fmt.Println("AfterDelete App", app)
	return nil
}

type AppZone struct {
	Zonename string `redis:"zonename" json:"zonename"`
	Owner    string `redis:"owner" json:"owner"`
}

func (appzone *AppZone) BeforeDelete(tx *gorm.DB) error {
	if err := tx.Delete(&AccountZone{}, "zonename = ?", appzone.Zonename).Error; err != nil {
		return err
	}

	return nil
}

type AppShare struct {
	Appname   string `redis:"appname" json:"appname"`
	Othername string `redis:"othername" json:"othername"`
}

type AppChatChannel struct {
	Appname     string `redis:"appname" json:"appname"`
	Channelname string `redis:"channelname" json:"channelname"`
}

type AppData struct {
	ID        uint64    `redis:"id" json:"id,string" gorm:"primary_key;AUTO_INCREMENT"`
	Account   string    `redis:"account" json:"account"`
	Appname   string    `redis:"appname" json:"appname"`
	Zonename  string    `redis:"zonename" json:"zonename"`
	Nickname  string    `redis:"nickname" json:"nickname"`
	Desc      string    `redis:"desc" json:"desc"`
	Sex       string    `redis:"sex" json:"sex"`
	Birthday  time.Time `redis:"birthday" json:"birthday"`
	Country   string    `redis:"country" json:"country"`
	Isbaned   bool      `redis:"isbaned" json:"isbaned" gorm:"tinyint(1);default:0"`
	Isjinyan  bool      `redis:"isjinyan" json:"isjinyan" gorm:"tinyint(1);default:0"`
	Regip     string    `redis:"regip" json:"regip"`
	Lastip    string    `redis:"lastip" json:"lastip"`
	Lastlogin time.Time `redis:"lastlogin" json:"lastlogin"`
	CreatedAt time.Time `redis:"createdate" json:"createdate"`

	Onlines []Online `json:"-" gorm:"foreignkey:Dataid;association_foreignkey:ID"`
	Friends []Friend `json:"-" gorm:"foreignkey:Dataid;association_foreignkey:ID"`
	Blacks  []Black  `json:"-" gorm:"foreignkey:Dataid;association_foreignkey:ID"`
	Groups  []Group  `json:"-" gorm:"foreignkey:Dataid;association_foreignkey:ID"`
}

//对外可见的数据，用于join查询。需要这些数据的结构需要组合该结构
type AppDataPublic struct {
	Account  string    `redis:"account" json:"account" gorm:"-"`
	Appname  string    `redis:"appname" json:"appname" gorm:"-"`
	Zonename string    `redis:"zonename" json:"zonename" gorm:"-"`
	Nickname string    `redis:"nickname" json:"nickname" gorm:"-"`
	Desc     string    `redis:"desc" json:"desc" gorm:"-"`
	Sex      string    `redis:"sex" json:"sex" gorm:"-"`
	Birthday time.Time `redis:"birthday" json:"birthday" gorm:"-"`
	Country  string    `redis:"country" json:"country" gorm:"-"`
	Isbaned  bool      `redis:"isbaned" json:"isbaned" gorm:"-"`
	Isjinyan bool      `redis:"isjinyan" json:"isjinyan" gorm:"-"`
	Isonline bool      `redis:"isonline" json:"isonline" gorm:"-"`
}

func (appdata *AppData) toAccountApp() *AccountApp {
	return &AccountApp{Account: appdata.Account, Appname: appdata.Appname}
}

func (appdata *AppData) toAccountZone() *AccountZone {
	return &AccountZone{Account: appdata.Account, Appname: appdata.Appname, Zonename: appdata.Zonename}
}

//AccountApp 账号所登录过的APP记录
type AccountApp struct {
	Account string `redis:"account" json:"account"`
	Appname string `redis:"appname" json:"appname"`
}

//AccountZone 账号所登录过的app区记录
type AccountZone struct {
	Account  string `redis:"account" json:"account"`
	Appname  string `redis:"appname" json:"appname"`
	Zonename string `redis:"zonename" json:"zonename"`
}

// func (appdata *AppData) MarshalJSON() ([]byte, error) {
// 	// 定义一个该结构体的别名
// 	type Alias AppData
// 	// 定义一个新的结构体
// 	tmpSt := struct {
// 		Alias
// 		Birthday   string `json:"birthday"`
// 		Lastlogin  string `json:"lastlogin"`
// 		CreateDate string `json:"createdate"`
// 	}{
// 		Alias:      (Alias)(*appdata),
// 		Birthday:   appdata.Birthday.Format("01/02/2006"),
// 		Lastlogin:  appdata.Lastlogin.Format("01/02/2006 15:04:05"),
// 		CreateDate: appdata.CreatedAt.Format("01/02/2006 15:04:05"),
// 	}
// 	return json.Marshal(tmpSt)
// }

func (appdata *AppData) BeforeDelete(tx *gorm.DB) error {
	fmt.Println("BeforeDelete AppData", appdata)

	if err := tx.Delete(&Online{}, "dataid = ?", appdata.ID).Error; err != nil {
		return err
	}

	if err := tx.Delete(&Friend{}, "dataid = ?", appdata.ID).Error; err != nil {
		return err
	}

	if err := tx.Delete(&Black{}, "dataid = ?", appdata.ID).Error; err != nil {
		return err
	}

	if err := tx.Delete(&Group{}, "dataid = ?", appdata.ID).Error; err != nil {
		return err
	}

	return nil
}

type Online struct {
	Dataid uint64 `redis:"dataid" json:"dataid,string" gorm:"not null"`
	// Account    string    `redis:"account" json:"account"`
	// Appname    string    `redis:"appname" json:"appname"`
	// Zonename   string    `redis:"zonename" json:"zonename"`
	Serveraddr string `redis:"serveraddr" json:"serveraddr"`
	//State      string `redis:"state" json:"state"`
	Platform  string    `redis:"platform" json:"platform"`
	CreatedAt time.Time `redis:"createdate" json:"createdate"`
}

type Friend struct {
	Dataid      uint64 `redis:"dataid" json:"dataid,string"`
	Otherdataid uint64 `redis:"otherdataid" json:"otherdataid,string"`
	// Account      string    `redis:"account" json:"account"`
	// Otheraccount string    `redis:"otheraccount" json:"otheraccount"`
	// Appname      string    `redis:"appname" json:"appname"`
	// Zonename     string    `redis:"zonename" json:"zonename"`
	Groupname string    `redis:"groupname" json:"groupname"`
	Comment   string    `redis:"comment" json:"comment"`
	CreatedAt time.Time `redis:"createdate" json:"createdate"`
}

type Black struct {
	Dataid      uint64 `redis:"dataid" json:"dataid,string"`
	Otherdataid uint64 `redis:"otherdataid" json:"otherdataid"`
	// Account      string    `redis:"account" json:"account"`
	// Otheraccount string    `redis:"otheraccount" json:"otheraccount"`
	// Appname      string    `redis:"appname" json:"appname"`
	// Zonename     string    `redis:"zonename" json:"zonename"`
	CreatedAt time.Time `redis:"createdate" json:"createdate"`
}

type Group struct {
	Groupname string `redis:"groupname" json:"groupname"`
	Dataid    uint64 `redis:"dataid" json:"dataid"`
	//Otherdataid uint64 `redis:"otherdataid" json:"otherdataid"`
	// Account  string `redis:"account" json:"account"`
	// Appname  string `redis:"appname" json:"appname"`
	// Zonename string `redis:"zonename" json:"zonename"`
}

type AccountBaned struct {
	Account   string    `redis:"account" json:"account" gorm:"unique;not null"`
	Why       string    `redis:"why" json:"why"`
	Dateline  time.Time `redis:"dateline" json:"dateline"`
	CreatedAt time.Time `redis:"createdate" json:"createdate"`
}

type AppBaned struct {
	Appname   string    `redis:"appname" json:"appname" gorm:"unique;not null"`
	Why       string    `redis:"why" json:"why"`
	Dateline  time.Time `redis:"dateline" json:"dateline"`
	CreatedAt time.Time `redis:"createdate" json:"createdate"`
}

type AppDataBaned struct {
	Dataid    uint64    `redis:"dataid" json:"dataid,string" gorm:"unique;not null"`
	Why       string    `redis:"why" json:"why"`
	Dateline  time.Time `redis:"dateline" json:"dateline"`
	CreatedAt time.Time `redis:"createdate" json:"createdate"`
}

type AppDataJinyan struct {
	Dataid    uint64    `redis:"dataid" json:"dataid,string" gorm:"unique;not null"`
	Why       string    `redis:"why" json:"why"`
	Dateline  time.Time `redis:"dateline" json:"dateline"`
	CreatedAt time.Time `redis:"createdate" json:"createdate"`
}

type Room struct {
	Rid       uint64    `redis:"rid" json:"rid,string" gorm:"primary_key;AUTO_INCREMENT"`
	Ownerid   uint64    `redis:"ownerid" json:"ownerid,string" gorm:"not null"`
	Roomtype  byte      `redis:"jointype" json:"jointype" gorm:"default:1"`
	Jointype  byte      `redis:"jointype" json:"jointype" gorm:"default:1"`
	Jieshao   string    `redis:"jieshao" json:"jieshao"`
	Notice    string    `redis:"notice" json:"notice"` //公告
	Password  string    `redis:"password" json:"-"`
	CreatedAt time.Time `redis:"createdate" json:"createdate"`

	//other info for join
	Msgsetting byte `redis:"msgsetting" json:"msgsetting" gorm:"-"` //RoomUser.Msgsetting
}

type RoomUser struct {
	Rid         uint64    `redis:"rid" json:"rid,string" gorm:"not null"`
	Dataid      uint64    `redis:"dataid" json:"dataid,string" gorm:"not null"`
	Isowner     bool      `redis:"isowner" json:"isowner"`
	Isadmin     bool      `redis:"isadmin" json:"isadmin"`
	Isjinyan    bool      `redis:"isjinyan" json:"isjinyan"`
	Displayname string    `redis:"displayname" json:"displayname"`
	Msgsetting  byte      `redis:"msgsetting" json:"msgsetting" gorm:"default:1"`
	CreatedAt   time.Time `redis:"createdate" json:"createdate"`

	//other info, for join
	AppDataPublic
}

var db_tables []interface{} = []interface{}{
	&Admin{},
	&AdminApp{},
	&Account{},
	&AccountAdminApp{},
	&App{},
	&AppZone{},
	&AppShare{},
	&AppData{},
	&AccountApp{},
	&AccountZone{},
	&Online{},
	&Friend{},
	&Black{},
	&Group{},
	&AccountBaned{},
	&AppBaned{},
	&AppDataBaned{},
	&AppDataJinyan{},

	&Room{},
	&RoomUser{},
}

type AppDataJson struct {
	Nickname string    `json:"nickname"`
	Desc     string    `json:"desc"`
	Sex      string    `json:"sex"`
	Birthday time.Time `json:"birthday"`
	Country  string    `json:"country"`
}

type FriendJson struct {
	Dataid    uint64 `json:"who,string"`
	Nickname  string `json:"nickname"`
	Desc      string `json:"desc"`
	Groupname string `json:"groupname"`
	Comment   string `json:"comment"`
}

type SearchUserJson struct {
	Dataid   uint64 `json:"who,string"`
	Nickname string `json:"nickname"`
	Country  string `json:"country"`
}

type SearchRoomJson struct {
	Dataid   uint64 `json:"who,string"`
	Nickname string `json:"nickname"`
	Country  string `json:"country"`
}
