package gtdb

//. "github.com/gtechx/base/common"

var account_table = &Account{}
var account_tablelist = []*Account{}

var online_table = &Online{}
var online_tablelist = []*Online{}

//每个app之间可以是独立的数据，也可以共享数据，根据你的设置
func (db *DBManager) CreateAccount(tbl_account *Account) error {
	retdb := db.sql.Create(tbl_account)
	return retdb.Error
}

func (db *DBManager) UpdateAccount(tbl_account *Account) error {
	retdb := db.sql.Save(tbl_account)
	return retdb.Error
}

func (db *DBManager) DeleteAccount(account string) error {
	retdb := db.sql.Delete(&Account{Account: account}, "Account = ?", account)
	return retdb.Error
}

func (db *DBManager) DeleteAccounts(accounts []string) error {
	tx := db.sql.Begin()
	for _, account := range accounts {
		if err := tx.Delete(&Account{Account: account}, "Account = ?", account).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func (db *DBManager) IsAccountExists(account string) (bool, error) {
	var count uint64
	retdb := db.sql.Model(account_table).Where("account = ?", account).Count(&count)
	return count > 0, retdb.Error
}

// func (db *DBManager) Updates(old interface{}, newval map[string]interface{}) error {
// 	db.sql.Model(old).Updates(newval)
// 	return db.sql
// }

func (db *DBManager) UpdatePassword(account, password string) error {
	retdb := db.sql.Model(account_table).Where("account = ?", account).Update("password", password)
	return retdb.Error
}

// func (db *DBManager) SetAppDataField(datakey *DataKey, fieldname string, value interface{}) error {
// 	db.Model(&Account{Account: account}).Update("password", password)
// 	return retdb.Error
// }

// func (db *DBManager) GetAppDataField(datakey *DataKey, fieldname string) (interface{}, error) {
// 	conn := db.redisPool.Get()
// 	defer conn.Close()
// 	ret, retdb := conn.Do("HGET", datakey.KeyAppDataHsetByAppidZonenameAccount, fieldname)
// 	return ret, retdb.Error
// }

// func (db *DBManager) SetMaxFriends(uid uint64, count int) error {
// 	conn := db.redisPool.Get()
// 	defer conn.Close()
// 	_, retdb := conn.Do("HSET", uid, "maxfriends", count)
// 	return retdb.Error
// }

// func (db *DBManager) SetDesc(uid uint64, desc string) error {
// 	conn := db.redisPool.Get()
// 	defer conn.Close()
// 	_, retdb := conn.Do("HSET", uid, "desc", desc)
// 	return retdb.Error
// }

// func (db *DBManager) IsUIDExists(uid uint64) (bool, error) {
// 	conn := db.redisPool.Get()
// 	defer conn.Close()
// 	ret, retdb := conn.Do("EXISTS", "uid:"+String(uid))
// 	return redis.Bool(ret, retdb.Error)
// }

// func (db *DBManager) GetUIDByAccount(account string) (uint64, error) {
// 	conn := db.redisPool.Get()
// 	defer conn.Close()
// 	ret, retdb := conn.Do("HGET", "account:uid", account)
// 	return redis.Uint64(ret, retdb.Error)
// }

// func (db *DBManager) GetAccountByUID(uid uint64) (string, error) {
// 	conn := db.redisPool.Get()
// 	defer conn.Close()
// 	ret, retdb := conn.Do("HGET", "uid:"+String(uid), "account")
// 	return redis.String(ret, retdb.Error)
// }

func (db *DBManager) GetAccount(account string) (*Account, error) {
	tbl_acc := &Account{}
	retdb := db.sql.Where("account = ?", account).First(tbl_acc)
	return tbl_acc, retdb.Error
}

func (db *DBManager) SetUserState(id uint64, state string) error {
	retdb := db.sql.Model(online_table).Where("dataid = ?", id).Update("state", state)
	return retdb.Error
}

var black_table = &Black{}
var black_tablelist = []*Black{}

func (db *DBManager) AddBlack(tbl_black *Black) error {
	retdb := db.sql.Create(tbl_black)
	return retdb.Error
}

func (db *DBManager) RemoveFromBlack(id, otherid uint64) error {
	retdb := db.sql.Delete(black_table, "dataid = ? AND otherdataid = ?", id, otherid)
	return retdb.Error
}

func (db *DBManager) IsInBlack(id, otherid uint64) (bool, error) {
	var count uint64
	retdb := db.sql.Model(black_table).Where("dataid = ? AND otherdataid = ?", id, otherid).Count(&count)
	return count > 0, retdb.Error
}

func (db *DBManager) GetBlackInfoList(id uint64) ([]*FriendJson, error) {
	blacklist := []*FriendJson{}
	retdb := db.sql.Table("gtchat_blacks").Where("gtchat_blacks.dataid = ?", id).Select("gtchat_blacks.otherdataid as dataid, gtchat_friends.comment, gtchat_app_data.nickname, gtchat_app_data.desc").Joins("join gtchat_app_data on gtchat_blacks.otherdataid = gtchat_app_data.id").Find(&blacklist)
	return blacklist, retdb.Error
}
