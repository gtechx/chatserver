package gtdb

import (
	"time"

	"github.com/jinzhu/gorm"
)

//. "github.com/gtechx/base/common"

//[set]app aid set
//[hset]app:aid:uid aid owner desc regdate
//[hset]app:aid:uid:config
var app_table = &App{}
var app_tablelist = []*App{}

var appzone_table = &AppZone{}
var appzone_tablelist = []*AppZone{}

type AppFilter struct {
	Appname         string
	Desc            string
	Share           string
	Createbegindate *time.Time
	Createenddate   *time.Time
}

func (filter *AppFilter) apply(db *gorm.DB) *gorm.DB {
	retdb := db
	if filter.Appname != "" {
		retdb = retdb.Where("appname LIKE ?", "%"+filter.Appname+"%")
	}
	if filter.Desc != "" {
		retdb = retdb.Where("desc LIKE ?", "%"+filter.Desc+"%")
	}
	if filter.Share != "" {
		retdb = retdb.Where("share = ?", filter.Share)
	}

	if filter.Createbegindate != nil {
		retdb = retdb.Where("created_at >= ?", *filter.Createbegindate)
	}
	if filter.Createenddate != nil {
		retdb = retdb.Where("created_at <= ?", *filter.Createenddate)
	}
	return retdb
}

//app op
func (db *DBManager) CreateApp(tbl_app *App) error {
	retdb := db.sql.Create(tbl_app)
	return retdb.Error
}

func (db *DBManager) UpdateApp(tbl_app *App) error {
	retdb := db.sql.Save(tbl_app)
	return retdb.Error
}

func (db *DBManager) DeleteApp(appname string) error {
	retdb := db.sql.Delete(&App{Appname: appname}, "appname = ?", appname)
	return retdb.Error
}

func (db *DBManager) DeleteApps(appnames []string) error {
	tx := db.sql.Begin()
	for _, appname := range appnames {
		if err := tx.Delete(&App{Appname: appname}, "appname = ?", appname).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func (db *DBManager) IsAppExists(appname string) (bool, error) {
	var count uint64
	retdb := db.sql.Model(app_table).Where("appname = ?", appname).Count(&count)
	return count > 0, retdb.Error
}

func (db *DBManager) SetAppField(appname, fieldname string, val interface{}) error {
	retdb := db.sql.Model(app_table).Where("appname = ?", appname).Update(fieldname, val)
	return retdb.Error
}

func (db *DBManager) GetAppField(appname, fieldname string) (*App, error) {
	app := &App{}
	retdb := db.sql.Select(fieldname).Where("appname = ?", appname).First(app)
	return app, retdb.Error
}

func (db *DBManager) GetAppCount(args ...*AppFilter) (uint64, error) {
	var count uint64
	retdb := db.sql.Model(app_table)
	if len(args) > 0 {
		filter := args[0]
		if filter != nil {
			retdb = filter.apply(retdb)
		}
	}
	retdb = retdb.Count(&count)
	return count, retdb.Error
}

func (db *DBManager) GetAppCountByAccount(account string, args ...*AppFilter) (uint64, error) {
	var count uint64
	retdb := db.sql.Model(app_table).Where("owner = ?", account)
	if len(args) > 0 {
		filter := args[0]
		if filter != nil {
			retdb = filter.apply(retdb)
		}
	}
	retdb = retdb.Count(&count)
	return count, retdb.Error
}

func (db *DBManager) GetApp(appname string) (*App, error) {
	app := &App{}
	retdb := db.sql.First(app, "appname = ?", appname)
	return app, retdb.Error
}

func (db *DBManager) GetAppList(offset, count int, args ...*AppFilter) ([]*App, error) {
	applist := []*App{}
	retdb := db.sql.Offset(offset).Limit(count)
	if len(args) > 0 {
		filter := args[0]
		if filter != nil {
			retdb = filter.apply(retdb)
		}
	}
	retdb = retdb.Find(&applist)
	return applist, retdb.Error
}

func (db *DBManager) GetAppListByAccount(account string, offset, count int, args ...*AppFilter) ([]*App, error) {
	applist := []*App{}
	retdb := db.sql.Offset(offset).Limit(count).Where("owner = ?", account)
	if len(args) > 0 {
		filter := args[0]
		if filter != nil {
			retdb = filter.apply(retdb)
		}
	}
	retdb = retdb.Find(&applist)
	return applist, retdb.Error
}

func (db *DBManager) GetAppOwner(appname string) (string, error) {
	app := &App{}
	retdb := db.sql.First(app, "appname = ?", appname)
	return app.Owner, retdb.Error
}

// func (db *DBManager) SetAppType(appid uint64, typestr string) error {
// 	conn := db.redisPool.Get()
// 	defer conn.Close()
// 	_, err := conn.Do("HSET", "app:"+String(appid), "type", typestr)
// 	return err
// }

// func (db *DBManager) GetAppType(appid uint64) (string, error) {
// 	conn := db.redisPool.Get()
// 	defer conn.Close()
// 	ret, err := conn.Do("HGET", "app:"+String(appid), "type")
// 	return redis.String(ret, err)
// }

// func (db *DBManager) IsAppIDExists(datakey *DataKey) (bool, error) {
// 	conn := db.redisPool.Get()
// 	defer conn.Close()
// 	ret, err := conn.Do("SISMEMBER", datakey.KeyAppSet, datakey.Appid)
// 	return redis.Bool(ret, err)
// }

func (db *DBManager) AddAppZone(tbl_appzone *AppZone) error {
	retdb := db.sql.Create(tbl_appzone)
	return retdb.Error
}

func (db *DBManager) RemoveAppZone(appname, zonename string) error {
	retdb := db.sql.Delete(&AppZone{Zonename: zonename}, "zonename = ? AND owner = ?", zonename, appname)
	return retdb.Error
}

func (db *DBManager) RemoveAppZones(appname string, zonenames []string) error {
	tx := db.sql.Begin()
	for _, zonename := range zonenames {
		if err := tx.Delete(&AppZone{Zonename: zonename}, "zonename = ? AND owner = ?", zonename, appname).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func (db *DBManager) GetAppZoneList(appname string) ([]*AppZone, error) {
	zonelist := []*AppZone{}
	retdb := db.sql.Where("owner = ?", appname).Find(&zonelist)
	return zonelist, retdb.Error
}

func (db *DBManager) IsAppZoneExists(appname, zonename string) (bool, error) {
	var count uint64
	retdb := db.sql.Model(appzone_table).Where("zonename = ? AND owner = ?", zonename, appname).Count(&count)
	return count > 0, retdb.Error
}

// func (db *DBManager) GetAppZoneName(datakey *DataKey) (string, error) {
// 	conn := db.redisPool.Get()
// 	defer conn.Close()
// 	ret, err := conn.Do("HGET", "app:"+String(appid), zone)
// 	return redis.String(ret, err)
// }

func (db *DBManager) SetShareApp(appname, otherappname string) error {
	retdb := db.sql.Model(app_table).Where("appname = ?", appname).Update("share", otherappname)
	return retdb.Error
}

func (db *DBManager) RemoveShareApp(appname string) error {
	retdb := db.sql.Model(app_table).Where("appname = ?", appname).Update("share", "")
	return retdb.Error
}

func (db *DBManager) IsShareWithOtherApp(appname string) (bool, error) {
	app := &App{}
	retdb := db.sql.First(app, "appname = ?", appname)
	return app.Share != "", retdb.Error
}

func (db *DBManager) GetShareApp(appname string) (string, error) {
	app := &App{}
	retdb := db.sql.First(app, "appname = ?", appname)
	return app.Share, retdb.Error
}

func (db *DBManager) GetShareAppList(appname string) ([]string, error) {
	sharelist := []*AppShare{}
	retdb := db.sql.Select("other_name").Where("appname = ?", appname).Find(&sharelist)

	slist := make([]string, len(sharelist))
	for i, share := range sharelist {
		slist[i] = share.Othername
	}
	return slist, retdb.Error
}
