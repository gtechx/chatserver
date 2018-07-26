package gtdb

import (
	"encoding/json"
	"time"

	"github.com/jinzhu/gorm"
)

var appdatabaned_table = &AppDataBaned{}
var appdatabaned_tablelist = []*AppDataBaned{}

type BanedAppData struct {
	ID       uint64 `redis:"id" json:"id,string"`
	Account  string `redis:"account" json:"account"`
	Appname  string `redis:"appname" json:"appname"`
	Zonename string `redis:"zonename" json:"zonename"`
	Nickname string `redis:"nickname" json:"nickname"`

	Why      string    `redis:"why" json:"why"`
	Dateline time.Time `redis:"dateline" json:"dateline"`
	Bandate  time.Time `redis:"bandate" json:"bandate"`
}

func (ps *BanedAppData) MarshalJSON() ([]byte, error) {
	// 定义一个该结构体的别名
	type Alias BanedAppData
	// 定义一个新的结构体
	tmpSt := struct {
		Alias
		EndDate    string `json:"enddate"`
		CreateDate string `json:"createdate"`
	}{
		Alias:      (Alias)(*ps),
		EndDate:    ps.Dateline.Format("01/02/2006 15:04:05"),
		CreateDate: ps.Bandate.Format("01/02/2006 15:04:05"),
	}
	return json.Marshal(tmpSt)
}

type BanedAppDataFilter struct {
	Account  string
	Appname  string
	Zonename string
	Nickname string

	Startdateline *time.Time
	Enddateline   *time.Time

	Banedbegindate *time.Time
	Banedenddate   *time.Time
}

func (filter *BanedAppDataFilter) apply(db *gorm.DB) *gorm.DB {
	retdb := db
	if filter.Account != "" {
		retdb = retdb.Where("account LIKE ?", "%"+filter.Account+"%")
	}
	if filter.Nickname != "" {
		retdb = retdb.Where("nickname LIKE ?", "%"+filter.Nickname+"%")
	}

	if filter.Startdateline != nil {
		retdb = retdb.Where("dataline >= ?", *filter.Startdateline)
	}

	if filter.Enddateline != nil {
		retdb = retdb.Where("dataline <= ?", *filter.Enddateline)
	}

	if filter.Banedbegindate != nil {
		retdb = retdb.Where("bandate >= ?", *filter.Banedbegindate)
	}
	if filter.Banedenddate != nil {
		retdb = retdb.Where("bandate <= ?", *filter.Banedenddate)
	}

	return retdb
}

func (db *DBManager) GetBanedAppDataCount(appname, zonename string, args ...*BanedAppDataFilter) (uint64, error) {
	var count uint64
	retdb := db.sql.Table("gtchat_app_data a")
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

	retdb = retdb.Joins("join gtchat_app_data_baneds b on b.dataid = a.id")

	retdb = retdb.Count(&count)
	return count, retdb.Error
}

func (db *DBManager) GetBanedAppDataList(appname, zonename string, offset, count int, args ...*BanedAppDataFilter) ([]*BanedAppData, error) {
	appdatalist := []*BanedAppData{}
	retdb := db.sql.Table("gtchat_app_data a").Offset(offset).Limit(count)
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

	retdb = retdb.Joins("join gtchat_app_data_baneds b on b.dataid = a.id")

	retdb = retdb.Select("a.*, b.*").Scan(&appdatalist)
	return appdatalist, retdb.Error
}

func (db *DBManager) GetBanedAppData(id uint64) (*BanedAppData, error) {
	appdata := &BanedAppData{}
	retdb := db.sql.Table("gtchat_app_data a")
	retdb = retdb.Joins("join gtchat_app_data_baneds b on b.dataid = a.id").Where("dataid = ?", id)
	retdb = retdb.Select("a.*, b.*").Limit(1).Scan(appdata)
	return appdata, retdb.Error
}

func (db *DBManager) BanAppDatas(ids []uint64, tbl_appdatabaned *AppDataBaned) error {
	tx := db.sql.Begin()
	for _, id := range ids {
		tbl_appdatabaned.Dataid = id
		if err := tx.Create(tbl_appdatabaned).Error; err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Model(appdata_table).Where("id = ?", id).Update("isbaned", true).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func (db *DBManager) UnbanAppDatas(ids []uint64) error {
	tx := db.sql.Begin()
	for _, id := range ids {
		if err := tx.Model(appdata_table).Where("id = ?", id).Update("isbaned", false).Error; err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Delete(appdatabaned_table).Where("dataid = ?", id).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func (db *DBManager) BanAppData(tbl_appdatabaned *AppDataBaned) error {
	tx := db.sql.Begin()
	if err := tx.Create(tbl_appdatabaned).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(appdata_table).Where("id = ?", tbl_appdatabaned.Dataid).Update("isbaned", true).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (db *DBManager) UnbanAppData(id uint64) error {
	tx := db.sql.Begin()
	if err := tx.Delete(appdatabaned_table).Where("dataid = ?", id).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(appdata_table).Where("id = ?", id).Update("isbaned", false).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
