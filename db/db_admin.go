package gtdb

import (
	"time"

	"github.com/jinzhu/gorm"
)

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
	retdb := db.sql.First(admin, "account = ?", account)
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

func (db *DBManager) BanAccounts(accounts []string) error {
	tx := db.sql.Begin()
	accdb := tx.Model(account_table)
	for _, account := range accounts {
		if err := accdb.Where("account = ?", account).Update("isbaned", true).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func (db *DBManager) UnbanAccounts(accounts []string) error {
	tx := db.sql.Begin()
	accdb := tx.Model(account_table)
	for _, account := range accounts {
		if err := accdb.Where("account = ?", account).Update("isbaned", false).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func (db *DBManager) BanAccount(account string) error {
	retdb := db.sql.Model(account_table).Where("account = ?", account).Update("isbaned", true)
	return retdb.Error
}

func (db *DBManager) UnbanAccount(account string) error {
	retdb := db.sql.Model(account_table).Where("account = ?", account).Update("isbaned", false)
	return retdb.Error
}

// func (db *DBManager) GetAccountList(offset, count int) ([]*Account, error) {
// 	accountlist := []*Account{}
// 	retdb := db.sql.Offset(offset).Limit(count).Where("account != ?", "admin").Find(&accountlist)
// 	return accountlist, retdb.Error
// }

type AccountFilter struct {
	Account         string
	Email           string
	Regip           string
	Createbegindate *time.Time
	Createenddate   *time.Time
}

func (filter *AccountFilter) apply(db *gorm.DB) *gorm.DB {
	retdb := db
	if filter.Account != "" {
		retdb = retdb.Where("account LIKE ?", "%"+filter.Account+"%")
	}
	if filter.Email != "" {
		retdb = retdb.Where("email LIKE ?", "%"+filter.Email+"%")
	}
	if filter.Regip != "" {
		retdb = retdb.Where("regip LIKE ?", "%"+filter.Regip+"%")
	}
	if filter.Createbegindate != nil {
		retdb = retdb.Where("created_at >= ?", *filter.Createbegindate)
	}
	if filter.Createenddate != nil {
		retdb = retdb.Where("created_at <= ?", *filter.Createenddate)
	}
	return retdb
}

func (db *DBManager) GetAccountCount(args ...*AccountFilter) (uint64, error) {
	var count uint64
	retdb := db.sql.Model(account_table).Where("account != ?", "admin")
	if len(args) > 0 {
		filter := args[0]
		if filter != nil {
			retdb = filter.apply(retdb)
		}
	}
	retdb = retdb.Count(&count)
	return count, retdb.Error
}

func (db *DBManager) GetAccountList(offset, count int, args ...*AccountFilter) ([]*Account, error) {
	accountlist := []*Account{}
	retdb := db.sql.Offset(offset).Limit(count).Where("account != ?", "admin")

	if len(args) > 0 {
		filter := args[0]
		if filter != nil {
			retdb = filter.apply(retdb)
		}
	}

	retdb = retdb.Find(&accountlist)
	return accountlist, retdb.Error
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
