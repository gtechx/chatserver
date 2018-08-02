package gtdb

import (
	"time"

	"github.com/jinzhu/gorm"
)

// type OnlineAppData struct {
// 	ID       uint64 `redis:"id" json:"id,string"`
// 	Account  string `redis:"account" json:"account"`
// 	Appname  string `redis:"appname" json:"appname"`
// 	Zonename string `redis:"zonename" json:"zonename"`
// 	Nickname string `redis:"nickname" json:"nickname"`
// 	//Isbaned   bool      `redis:"isbaned" json:"isbaned"`
// 	Isjinyan bool `redis:"isjinyan" json:"isjinyan"`
// 	// Desc     string    `redis:"desc" json:"desc"`
// 	Sex string `redis:"sex" json:"sex"`
// 	// Birthday time.Time `redis:"birthday" json:"-"`
// 	Country string `redis:"country" json:"country"`
// 	// Regip     string    `redis:"regip" json:"regip"`

// 	Serveraddr string    `redis:"serveraddr" json:"serveraddr"`
// 	Platform   string    `redis:"platform" json:"platform"`
// 	Onlinedate time.Time `redis:"onlinedate" json:"-"`
// }

// func (ps *OnlineAppData) MarshalJSON() ([]byte, error) {
// 	// 定义一个该结构体的别名
// 	type Alias OnlineAppData
// 	// 定义一个新的结构体
// 	tmpSt := struct {
// 		Alias
// 		CreateDate string `json:"createdate"`
// 	}{
// 		Alias:      (Alias)(*ps),
// 		CreateDate: ps.Onlinedate.Format("01/02/2006 15:04:05"),
// 	}
// 	return json.Marshal(tmpSt)
// }

type OnlineFilter struct {
	Account string
	// Appname  string
	// Zonename string
	Nickname string
	Sex      string
	Country  string

	Isjinyan int

	Serveraddr string
	Platform   string

	Onlinebegindate *time.Time
	Onlineenddate   *time.Time
}

func (filter *OnlineFilter) apply(db *gorm.DB) *gorm.DB {
	retdb := db
	if filter.Account != "" {
		retdb = retdb.Where("account LIKE ?", "%"+filter.Account+"%")
	}
	if filter.Nickname != "" {
		retdb = retdb.Where("nickname LIKE ?", "%"+filter.Nickname+"%")
	}

	if filter.Sex != "" {
		retdb = retdb.Where("sex = ?", filter.Sex)
	}

	if filter.Country != "" {
		retdb = retdb.Where("country = ?", filter.Country)
	}

	if filter.Isjinyan != 0 {
		if filter.Isjinyan == 2 {
			retdb = retdb.Where("isjinyan = 1")
		} else {
			retdb = retdb.Where("isjinyan = 0")
		}
	}

	if filter.Serveraddr != "" {
		retdb = retdb.Where("serveraddr LIKE ?", "%"+filter.Serveraddr+"%")
	}

	if filter.Platform != "" {
		retdb = retdb.Where("platform LIKE ?", "%"+filter.Platform+"%")
	}

	if filter.Onlinebegindate != nil {
		retdb = retdb.Where("onlinedate >= ?", *filter.Onlinebegindate)
	}
	if filter.Onlineenddate != nil {
		retdb = retdb.Where("onlinedate <= ?", *filter.Onlineenddate)
	}

	return retdb
}

func (db *DBManager) GetOnlineCount(appname, zonename string, args ...*OnlineFilter) (uint64, error) {
	var count uint64
	retdb := db.sql.Table(db.sql.prefix + "app_data a")
	if appname != "" {
		retdb = retdb.Where("appname = ?", appname)
	}
	if zonename != "" {
		retdb = retdb.Where("zonename = ?", zonename)
	}

	if len(args) > 0 {
		filter := args[0]
		if filter != nil {
			retdb = filter.apply(retdb)
		}
	}

	retdb = retdb.Joins("join " + db.sql.prefix + "online b on b.dataid = a.id")

	retdb = retdb.Count(&count)
	return count, retdb.Error
}

func (db *DBManager) GetOnlineList(appname, zonename string, offset, count int, args ...*OnlineFilter) ([]*Online, error) {
	appdatalist := []*Online{}
	retdb := db.sql.Table(db.sql.prefix + "app_data a").Offset(offset).Limit(count)
	if appname != "" {
		retdb = retdb.Where("appname = ?", appname)
	}
	if zonename != "" {
		retdb = retdb.Where("zonename = ?", zonename)
	}

	if len(args) > 0 {
		filter := args[0]
		if filter != nil {
			retdb = filter.apply(retdb)
		}
	}

	retdb = retdb.Joins("join " + db.sql.prefix + "online b on b.dataid = a.id")

	retdb = retdb.Select("a.*, b.*").Scan(&appdatalist)
	return appdatalist, retdb.Error
}

func (db *DBManager) GetOnline(id uint64) (*Online, error) {
	appdata := &Online{}
	retdb := db.sql.Table(db.sql.prefix + "app_data a")
	retdb = retdb.Joins("join "+db.sql.prefix+"online b on b.dataid = a.id").Where("dataid = ?", id)
	retdb = retdb.Select("a.*, b.*").Limit(1).Scan(appdata)
	return appdata, retdb.Error
}

//多端登录的时候，会有多条online信息，因为每个端会可能会连接到不同的服务器
func (db *DBManager) GetUserOnlineInfoList(id uint64) ([]*Online, error) {
	var onlinelist []*Online
	retdb := db.sql.Model(online_table).Where("dataid = ?", id).Find(&onlinelist)
	return onlinelist, retdb.Error
}

func (db *DBManager) GetUserOnlineAddrList(id uint64) ([]string, error) {
	var addrlist []string
	retdb := db.sql.Model(online_table).Where("dataid = ?", id).Pluck("distinct serveraddr", &addrlist)
	return addrlist, retdb.Error
}

func (db *DBManager) SetUserOnline(tbl_online *Online) error {
	retdb := db.sql.Create(tbl_online)
	return retdb.Error
}

func (db *DBManager) SetUserOffline(id uint64, platform string) error {
	retdb := db.sql.Delete(online_table).Where("dataid = ? AND platform = ?", id, platform)
	return retdb.Error
}

func (db *DBManager) IsUserOnline(id uint64) (bool, error) {
	var count uint64
	retdb := db.sql.Model(online_table).Where("dataid = ?", id).Count(&count)
	return count > 0, retdb.Error
}

func (db *DBManager) ClearOnlineInfo(serveraddr string) error {
	retdb := db.sql.Delete(online_table, "serveraddr = ?", serveraddr)
	return retdb.Error
}
