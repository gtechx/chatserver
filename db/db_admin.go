package gtdb

import (
	"github.com/garyburd/redigo/redis"
	. "github.com/gtechx/base/common"
)

//[hashes]admin pair(uid, privilege) --管理员权限
//[sets]online uid --在线用户uid
admin_table := &Admin{}

func (db *DBManager) IsAdmin(account string) (bool, error) {
	var count uint64
	db.sql.Model(admin_table).Where("account = ?", account).Count(&count)
	return count > 0, db.sql.Error
}

func (db *DBManager) AddAdmin(tbl_admin *Admin) error {
	db.sql.Create(tbl_admin)
	return db.sql.Error
}

func (db *DBManager) RemoveAdmin(account string) error {
	db.sql.Delete(admin_table, "account = ?", account)
	return db.sql.Error
}

func (db *DBManager) GetAdmin(account string) (*Admin, error) {
	admin := &Admin{}
	db.sql.First(admin, account)
	return admin, db.sql.Error
}

func (db *DBManager) UpdateAdmin(tbl_admin *Admin) error {
	db.sql.Save(&user)
	return db.sql.Error
}

func (db *DBManager) GetAdminList(offset, count int) ([]*Admin, error) {
	adminlist := []*Admin{}
	db.sql.Offset(offset).Limit(count).Find(&adminlist)
	return adminlist, err
}

func (db *DBManager) GetUserOnlineByAppname(offset, count int) ([]*Online, error) {
	onlinelist := []*Online{}
	db.sql.Offset(offset).Limit(count).Find(&onlinelist)
	return onlinelist, err
}

func (db *DBManager) GetUserOnlineByAppname(appname string, offset, count int) ([]*Online, error) {
	onlinelist := []*Online{}
	db.sql.Offset(offset).Limit(count).Where("app_name = ?", appname).Find(&onlinelist)
	return onlinelist, err
}

func (db *DBManager) GetUserOnlineByZonename(zonename string, offset, count int) ([]*Online, error) {
	onlinelist := []*Online{}
	db.sql.Offset(offset).Limit(count).Where("zone_name = ?", zonename).Find(&onlinelist)
	return onlinelist, err
}

func (db *DBManager) GetUserOnlineByAppnameZonename(appname, zonename string, offset, count int) ([]*Online, error) {
	onlinelist := []*Online{}
	db.sql.Offset(offset).Limit(count).Where("app_name = ? zone_name = ?", appname, zonename).Find(&onlinelist)
	return onlinelist, err
}
