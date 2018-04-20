package gtdb

import "time"
import "github.com/jinzhu/gorm"

//. "github.com/gtechx/base/common"

var appdata_table = &AppData{}
var appdata_tablelist = []*AppData{}

func (db *DBManager) CreateAppData(tbl_appdata *AppData) error {
	retdb := db.sql.Create(tbl_appdata)
	return retdb.Error
}

func (db *DBManager) UpdateAppData(tbl_appdata *AppData) error {
	retdb := db.sql.Save(tbl_appdata)
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

func (db *DBManager) IsNicknameExists(appname, zonename, nickname string) (bool, error) {
	var count uint64
	retdb := db.sql.Model(appdata_table).Where("appname = ?", appname).Where("zonename = ?", zonename).Where("nickname = ?", nickname).Count(&count)
	return count > 0, retdb.Error
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
