package gtdb

import (
	"time"

	"github.com/garyburd/redigo/redis"
	. "github.com/gtechx/base/common"
)

//[set]app aid set
//[hset]app:aid:uid aid owner desc regdate
//[hset]app:aid:uid:config
app_table := &App{}
app_tablelist := []*App{}

appzone_table := &AppZone{}
appzone_tablelist := []*AppZone{}

//app op
func (db *DBManager) CreateApp(tbl_app *App) error {
	db.sql.Create(tbl_app)
	return db.sql.Error
}

func (db *DBManager) DeleteApp(appname string) error {
	db.sql.Delete(app_table, "name = ?", appname)
	return db.sql.Error
}

func (db *DBManager) IsAppExists(appname string) (bool, error) {
	var count uint64
	db.sql.Model(app_table).Where("name = ?", appname).Count(&count)
	return count > 0, db.sql.Error
}

func (db *DBManager) SetAppField(appname, fieldname string, val interface{}) error {
	db.Model(app_table).Where("name = ?", appname).Update(fieldname, val)
	return db.sql.Error
}

func (db *DBManager) GetAppField(appname, fieldname string) (*App, error) {
	app := &App{}
	db.sql.Select(fieldname).Where("name = ?", appname).First(app)
	return app, db.sql.Error
}

func (db *DBManager) GetAppCount() (uint64, error) {
	var count uint64
	db.sql.Find(&app_tablelist).Count(&count)
	return count, db.sql.Error
}

func (db *DBManager) GetAppCountByAccount(account string) (uint64, error) {
	var count uint64
	db.sql.Find(&app_tablelist).Where("owner = ?", account).Count(&count)
	return count, db.sql.Error
}

func (db *DBManager) GetApp(appname string) (*App, error) {
	app := &App{}
	db.sql.First(app, appname)
	return app, db.sql.Error
}

func (db *DBManager) GetAppList(offset, count int) ([]*App, error) {
	applist := []*App{}
	db.sql.Offset(offset).Limit(count).Find(&applist)
	return applist, err
}

func (db *DBManager) GetAppByAccount(account string, offset, count int) ([]*App, error) {
	applist := []*App{}
	db.sql.Offset(offset).Limit(count).Where("owner = ?", account).Find(&applist)
	return applist, err
}

func (db *DBManager) GetAppOwner(appname string) (string, error) {
	app := &App{}
	db.sql.Where("name = ?", appname).First(app)
	return app.Owner, db.sql.Error
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
	db.sql.Create(tbl_appzone)
	return db.sql.Error
}

func (db *DBManager) RemoveAppZone(appname, zonename string) error {
	db.sql.Delete(appzone_table, "name = ? AND owner = ?", zonename, appname)
	return db.sql.Error
}

func (db *DBManager) GetAppZones(appname string) ([]*AppZone, error) {
	zonelist := []*AppZone{}
	db.sql.Where("owner = ?", appname).Find(&zonelist)
	return zonelist, err
}

func (db *DBManager) IsAppZoneExists(appname, zonename string) (bool, error) {
	var count uint64
	db.sql.Model(appzone_table).Where("name = ? AND owner = ?", zonename, appname).Count(&count)
	return count > 0, db.sql.Error
}

// func (db *DBManager) GetAppZoneName(datakey *DataKey) (string, error) {
// 	conn := db.redisPool.Get()
// 	defer conn.Close()
// 	ret, err := conn.Do("HGET", "app:"+String(appid), zone)
// 	return redis.String(ret, err)
// }

func (db *DBManager) AddShareApp(appname, otherappname string) error {
	db.Model(app_table).Where("name = ?", appname).Update("share", otherappname)
	return db.sql.Error
}

func (db *DBManager) RemoveShareApp(appname string) error {
	db.Model(app_table).Where("name = ?", appname).Update("share", "")
	return db.sql.Error
}

func (db *DBManager) IsShareWithOtherApp(appname string) (bool, error) {
	app := &App{}
	db.sql.First(app, appname)
	return app.Share != "", db.sql.Error
}

func (db *DBManager) GetShareApp(appname string) (string, error) {
	app := &App{}
	db.sql.First(app, appname)
	return app.Share, db.sql.Error
}

func (db *DBManager) GetShareAppList(appname string) ([]string, error) {
	sharelist := []*AppShare{}
	db.sql.Select("other_name").Where("name = ?", appname).Find(&sharelist)

	slist := make([]string, len(sharelist))
	for i, name := range sharelist {
		slist[i] = name
	}
	return slist, db.sql.Error
}
