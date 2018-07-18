package gtdb

import "time"
import "github.com/jinzhu/gorm"
import . "github.com/gtechx/base/common"

//. "github.com/gtechx/base/common"

var appdata_table = &AppData{}
var appdata_tablelist = []*AppData{}

func (db *DBManager) CreateAppData(tbl_appdata *AppData) error {
	tx := db.sql.Begin()
	tmpdb := tx.Create(tbl_appdata)
	if err := tmpdb.Error; err != nil {
		tx.Rollback()
		return err
	}
	var count uint64
	if err := tx.Model(&AccountApp{}).Where("account = ?", tbl_appdata.Account).Where("appname = ?", tbl_appdata.Appname).Count(&count).Error; err != nil {
		tx.Rollback()
		return err
	}
	if count == 0 {
		if err := tx.Create(tbl_appdata.toAccountApp()).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Model(&AccountZone{}).Where("account = ?", tbl_appdata.Account).Where("appname = ?", tbl_appdata.Appname).Where("zonename = ?", tbl_appdata.Zonename).Count(&count).Error; err != nil {
		tx.Rollback()
		return err
	}
	if count == 0 {
		if err := tx.Create(tbl_appdata.toAccountZone()).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	id := Uint64(tmpdb.Value)

	if err := tx.Create(&Group{Groupname: "MyFriends", Dataid: id}).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (db *DBManager) UpdateAppData(tbl_appdata *AppData) error {
	retdb := db.sql.Save(tbl_appdata)
	return retdb.Error
}

func (db *DBManager) UpdateAppDataByMap(data map[string]interface{}) error {
	retdb := db.sql.Model(appdata_table).Updates(data)
	return retdb.Error
}

func (db *DBManager) UpdateLastLoginInfo(id uint64, ip string, date time.Time) error {
	retdb := db.sql.Model(appdata_table).Where("id = ?", id).Updates(map[string]interface{}{"lastip": ip, "lastlogin": date})
	return retdb.Error
}

func (db *DBManager) DeleteAppData(id uint64) error {
	retdb := db.sql.Delete(&AppData{ID: id}, "id = ?", id)
	return retdb.Error
}

func (db *DBManager) DeleteAppDatas(ids []uint64) error {
	tx := db.sql.Begin()
	for _, id := range ids {
		if err := tx.Delete(&AppData{ID: id}, "id = ?", id).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func (db *DBManager) GetAppData(id uint64) (*AppData, error) {
	tbl_appdata := &AppData{}
	retdb := db.sql.Where("id = ?", id).First(tbl_appdata)
	return tbl_appdata, retdb.Error
}

func (db *DBManager) IsAppDataExists(id uint64) (bool, error) {
	var count uint64
	retdb := db.sql.Model(appdata_table).Where("id = ?", id).Count(&count)
	return count > 0, retdb.Error
}

func (db *DBManager) IsNicknameExists(appname, zonename, account, nickname string) (bool, error) {
	var count uint64
	retdb := db.sql.Model(appdata_table).Where("appname = ?", appname).Where("zonename = ?", zonename).Where("account = ?", account).Where("nickname = ?", nickname).Count(&count)
	return count > 0, retdb.Error
}

func (db *DBManager) GetAppDataIdList(appname, zonename, account string) ([]uint64, error) {
	appdatalist := []*AppData{}
	retdb := db.sql.Model(appdata_table).Select("id").Where("appname = ?", appname).Where("zonename = ?", zonename).Where("account = ?", account).Find(&appdatalist)

	if retdb.Error != nil {
		return nil, retdb.Error
	}

	idlist := make([]uint64, len(appdatalist))
	for i, appdata := range appdatalist {
		idlist[i] = appdata.ID
	}

	return idlist, retdb.Error
}

// func (db *DBManager) GetAppDataCountByApp(appname string) (uint64, error) {
// 	var count uint64
// 	retdb := db.sql.Model(appdata_table).Where("appname = ?", appname).Count(&count)
// 	return count, retdb.Error
// }

// func (db *DBManager) GetAppDataCountByAppZone(appname, zonename string) (uint64, error) {
// 	var count uint64
// 	retdb := db.sql.Model(appdata_table).Where("appname = ?", appname).Where("zonename = ?", zonename).Count(&count)
// 	return count, retdb.Error
// }

// func (db *DBManager) GetAppDataList(appname, zonename, account string, offset, count int) ([]*AppData, error) {
// 	appdatalist := []*AppData{}
// 	retdb := db.sql.Offset(offset).Limit(count)
// 	if appname != "" {
// 		retdb = retdb.Where("appname = ?", appname)
// 	}
// 	if zonename != "" {
// 		retdb = retdb.Where("zonename = ?", zonename)
// 	}
// 	if account != "" {
// 		retdb = retdb.Where("account = ?", account)
// 	}
// 	retdb = retdb.Find(&appdatalist)
// 	return appdatalist, retdb.Error
// }

// func (db *DBManager) GetAppDataListByAccount(appname, zonename, account string, offset, count int) ([]*AppData, error) {
// 	appdatalist := []*AppData{}
// 	retdb := db.sql.Offset(offset).Limit(count).Where("appname = ?", appname)
// 	if zonename != "" {
// 		retdb = retdb.Where("zonename = ?", zonename)
// 	}
// 	if account != "" {
// 		retdb = retdb.Where("account = ?", account)
// 	}
// 	retdb = retdb.Find(&appdatalist)
// 	return appdatalist, retdb.Error
// }

type AppDataFilter struct {
	//Account            string
	Nickname           string
	Desc               string
	Sex                string
	Country            string
	Regip              string
	Lastip             string
	Birthdaybegindate  *time.Time
	Birthdayenddate    *time.Time
	Lastloginbegindate *time.Time
	Lastloginenddate   *time.Time
	Createbegindate    *time.Time
	Createenddate      *time.Time
}

func (filter *AppDataFilter) apply(db *gorm.DB) *gorm.DB {
	retdb := db
	if filter.Nickname != "" {
		retdb = retdb.Where("nickname LIKE ?", "%"+filter.Nickname+"%")
	}
	if filter.Desc != "" {
		retdb = retdb.Where("desc LIKE ?", "%"+filter.Desc+"%")
	}
	if filter.Sex != "" {
		retdb = retdb.Where("sex = ?", filter.Sex)
	}
	if filter.Country != "" {
		retdb = retdb.Where("country = ?", filter.Country)
	}
	if filter.Regip != "" {
		retdb = retdb.Where("regip LIKE ?", "%"+filter.Regip+"%")
	}
	if filter.Lastip != "" {
		retdb = retdb.Where("lastip LIKE ?", "%"+filter.Lastip+"%")
	}

	if filter.Birthdaybegindate != nil {
		retdb = retdb.Where("birthday >= ?", *filter.Birthdaybegindate)
	}
	if filter.Birthdayenddate != nil {
		retdb = retdb.Where("birthday <= ?", *filter.Birthdayenddate)
	}
	if filter.Lastloginbegindate != nil {
		retdb = retdb.Where("lastlogin >= ?", *filter.Lastloginbegindate)
	}
	if filter.Lastloginenddate != nil {
		retdb = retdb.Where("lastlogin <= ?", *filter.Lastloginenddate)
	}
	if filter.Createbegindate != nil {
		retdb = retdb.Where("created_at >= ?", *filter.Createbegindate)
	}
	if filter.Createenddate != nil {
		retdb = retdb.Where("created_at <= ?", *filter.Createenddate)
	}
	return retdb
}

func (db *DBManager) GetAppDataCount(appname, zonename, account string, args ...*AppDataFilter) (uint64, error) {
	var count uint64
	retdb := db.sql.Model(appdata_table)
	if appname != "" {
		retdb = retdb.Where("appname = ?", appname)
	}
	if zonename != "" {
		retdb = retdb.Where("zonename = ?", zonename)
	}
	if account != "" {
		retdb = retdb.Where("account = ?", account)
	}
	if len(args) > 0 {
		filter := args[0]
		if filter != nil {
			retdb = filter.apply(retdb)
		}
	}
	retdb = retdb.Count(&count)
	return count, retdb.Error
}

func (db *DBManager) GetAppDataList(appname, zonename, account string, offset, count int, args ...*AppDataFilter) ([]*AppData, error) {
	appdatalist := []*AppData{}
	retdb := db.sql.Offset(offset).Limit(count)
	if appname != "" {
		retdb = retdb.Where("appname = ?", appname)
	}
	if zonename != "" {
		retdb = retdb.Where("zonename = ?", zonename)
	}
	if account != "" {
		retdb = retdb.Where("account = ?", account)
	}
	if len(args) > 0 {
		filter := args[0]
		if filter != nil {
			retdb = filter.apply(retdb)
		}
	}

	retdb = retdb.Find(&appdatalist)
	return appdatalist, retdb.Error
}

func (db *DBManager) GetAccountAppList(account string) ([]*AccountApp, error) {
	accountapplist := []*AccountApp{}
	retdb := db.sql.Where("account = ?", account).Find(&accountapplist)
	return accountapplist, retdb.Error
}

func (db *DBManager) GetAccountZoneList(account, appname string) ([]*AccountZone, error) {
	accountzonelist := []*AccountZone{}
	retdb := db.sql.Where("account = ?", account).Where("appname = ?", appname).Find(&accountzonelist)
	return accountzonelist, retdb.Error
}

func (db *DBManager) BanAppDatas(ids []uint64) error {
	tx := db.sql.Begin()
	appdatadb := tx.Model(appdata_table)
	for _, id := range ids {
		if err := appdatadb.Where("id = ?", id).Update("isbaned", true).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func (db *DBManager) UnbanAppDatas(ids []uint64) error {
	tx := db.sql.Begin()
	appdatadb := tx.Model(appdata_table)
	for _, id := range ids {
		if err := appdatadb.Where("id = ?", id).Update("isbaned", false).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func (db *DBManager) BanAppData(id uint64) error {
	retdb := db.sql.Model(appdata_table).Where("id = ?", id).Update("isbaned", true)
	return retdb.Error
}

func (db *DBManager) UnbanAppData(id uint64) error {
	retdb := db.sql.Model(appdata_table).Where("id = ?", id).Update("isbaned", false)
	return retdb.Error
}
