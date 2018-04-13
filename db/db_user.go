package gtdb

import (
	. "github.com/gtechx/base/common"
)

//每个app之间可以是独立的数据，也可以共享数据，根据你的设置
func (db *DBManager) CreateAccount(tbl_account *Account) error {
	db.sql.Create(tbl_account)
	return db.sql.Error
}

func (db *DBManager) IsAccountExists(tbl_account *Account) (bool, error) {
	var count uint64
	db.sql.Model(tbl_account).Where("account = ?", tbl_account.Account).Count(&count)
	return count > 0, db.sql.Error
}

// func (db *DBManager) Updates(old interface{}, newval map[string]interface{}) error {
// 	db.sql.Model(old).Updates(newval)
// 	return db.sql.Error
// }

func (db *DBManager) UpdatePassword(tbl_account *Account) error {
	db.Model(tbl_account).Update("password", tbl_account.Password)
	return db.sql.Error
}

func (db *DBManager) CreateAppData(tbl_appdata *AppData) error {
	db.sql.Create(tbl_appdata)
	return db.sql.Error
}

func (db *DBManager) DeleteAppData(tbl_appdata *AppData) error {
	db.Delete(tbl_appdata, "appname = ? AND zonename = ? AND account = ?", tbl_appdata.AppName, tbl_appdata.ZoneName, tbl_appdata.Account)
	return db.sql.Error
}

func (db *DBManager) IsAppDataExists(tbl_appdata *AppData) (bool, error) {
	var count uint64
	db.Model(tbl_appdata).Where("id = ?", tbl_appdata.ID).Count(&count)
	return count > 0, db.Error
}

// func (db *DBManager) SetAppDataField(datakey *DataKey, fieldname string, value interface{}) error {
// 	db.Model(&Account{Account: account}).Update("password", password)
// 	return err
// }

// func (db *DBManager) GetAppDataField(datakey *DataKey, fieldname string) (interface{}, error) {
// 	conn := db.redisPool.Get()
// 	defer conn.Close()
// 	ret, err := conn.Do("HGET", datakey.KeyAppDataHsetByAppidZonenameAccount, fieldname)
// 	return ret, err
// }

// func (db *DBManager) SetMaxFriends(uid uint64, count int) error {
// 	conn := db.redisPool.Get()
// 	defer conn.Close()
// 	_, err := conn.Do("HSET", uid, "maxfriends", count)
// 	return err
// }

// func (db *DBManager) SetDesc(uid uint64, desc string) error {
// 	conn := db.redisPool.Get()
// 	defer conn.Close()
// 	_, err := conn.Do("HSET", uid, "desc", desc)
// 	return err
// }

// func (db *DBManager) IsUIDExists(uid uint64) (bool, error) {
// 	conn := db.redisPool.Get()
// 	defer conn.Close()
// 	ret, err := conn.Do("EXISTS", "uid:"+String(uid))
// 	return redis.Bool(ret, err)
// }

// func (db *DBManager) GetUIDByAccount(account string) (uint64, error) {
// 	conn := db.redisPool.Get()
// 	defer conn.Close()
// 	ret, err := conn.Do("HGET", "account:uid", account)
// 	return redis.Uint64(ret, err)
// }

// func (db *DBManager) GetAccountByUID(uid uint64) (string, error) {
// 	conn := db.redisPool.Get()
// 	defer conn.Close()
// 	ret, err := conn.Do("HGET", "uid:"+String(uid), "account")
// 	return redis.String(ret, err)
// }

func (db *DBManager) GetPassword(account string, out *Account) error {
	db.sql.Where("account = ?", account).First(out)
	return db.sql.Error
}

func (db *DBManager) SetUserOnline(tbl_online *Online) error {
	db.sql.Create(tbl_online)
	return db.sql.Error
}

func (db *DBManager) SetUserOffline(tbl_online *Online) error {
	db.Delete(tbl_online, "appname = ? AND zonename = ? AND account = ?", tbl_online.AppName, tbl_online.ZoneName, tbl_online.Account)
	return db.sql.Error
}

func (db *DBManager) IsUserOnline(tbl_online *Online) (bool, error) {
	var count uint64
	db.sql.Model(tbl_online).Where("appname = ? AND zonename = ? AND account = ?", tbl_online.AppName, tbl_online.ZoneName, tbl_online.Account).Count(&count)
	return count > 0, db.sql.Error
}

func (db *DBManager) SetUserState(tbl_online *Online) error {
	db.Model(tbl_online).Where("appname = ? AND zonename = ? AND account = ?", tbl_online.AppName, tbl_online.ZoneName, tbl_online.Account).Update("state", tbl_online.State)
	return db.sql.Error
}

func (db *DBManager) AddUserToBlack(tbl_black *Black) error {
	db.sql.Create(tbl_black)
	return db.sql.Error
}

func (db *DBManager) RemoveUserFromBlack(tbl_black *Black) error {
	db.Delete(tbl_black, "appname = ? AND zonename = ? AND account = ? AND other_account = ?", tbl_black.AppName, tbl_black.ZoneName, tbl_black.Account, tbl_black.OtherAccount)
	return db.sql.Error
}

func (db *DBManager) IsUserInBlack(tbl_black *Black) (bool, error) {
	var count uint64
	db.sql.Model(tbl_black).Where("appname = ? AND zonename = ? AND account = ? AND other_account = ?", tbl_black.AppName, tbl_black.ZoneName, tbl_black.Account, tbl_black.OtherAccount).Count(&count)
	return count > 0, db.sql.Error
}
