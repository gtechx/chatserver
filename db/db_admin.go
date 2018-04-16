package gtdb

//. "github.com/gtechx/base/common"

//[hashes]admin pair(uid, privilege) --管理员权限
//[sets]online uid --在线用户uid
var admin_table = &Admin{}

func (db *DBManager) IsAdmin(account string) (bool, error) {
	var count uint64
	retdb := db.sql.Model(admin_table).Where("account = ?", account).Count(&count)
	return count > 0, retdb.Error
}

func (db *DBManager) AddAdmin(tbl_admin *Admin) error {
	retdb := db.sql.Create(tbl_admin)
	return retdb.Error
}

func (db *DBManager) RemoveAdmin(account string) error {
	retdb := db.sql.Delete(admin_table, "account = ?", account)
	return retdb.Error
}

func (db *DBManager) GetAdmin(account string) (*Admin, error) {
	admin := &Admin{}
	retdb := db.sql.First(admin, account)
	return admin, retdb.Error
}

func (db *DBManager) UpdateAdmin(tbl_admin *Admin) error {
	retdb := db.sql.Save(tbl_admin)
	return retdb.Error
}

func (db *DBManager) GetAdminList(offset, count int) ([]*Admin, error) {
	adminlist := []*Admin{}
	retdb := db.sql.Offset(offset).Limit(count).Find(&adminlist)
	return adminlist, retdb.Error
}

func (db *DBManager) GetUserOnline(offset, count int) ([]*Online, error) {
	onlinelist := []*Online{}
	retdb := db.sql.Offset(offset).Limit(count).Find(&onlinelist)
	return onlinelist, retdb.Error
}

func (db *DBManager) GetUserOnlineByAppname(appname string, offset, count int) ([]*Online, error) {
	onlinelist := []*Online{}
	retdb := db.sql.Offset(offset).Limit(count).Where("app_name = ?", appname).Find(&onlinelist)
	return onlinelist, retdb.Error
}

func (db *DBManager) GetUserOnlineByZonename(zonename string, offset, count int) ([]*Online, error) {
	onlinelist := []*Online{}
	retdb := db.sql.Offset(offset).Limit(count).Where("zone_name = ?", zonename).Find(&onlinelist)
	return onlinelist, retdb.Error
}

func (db *DBManager) GetUserOnlineByAppnameZonename(appname, zonename string, offset, count int) ([]*Online, error) {
	onlinelist := []*Online{}
	retdb := db.sql.Offset(offset).Limit(count).Where("app_name = ? zone_name = ?", appname, zonename).Find(&onlinelist)
	return onlinelist, retdb.Error
}
