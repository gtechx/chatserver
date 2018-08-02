package gtdb

import (
	"time"

	"github.com/jinzhu/gorm"
)

var appdatajinyan_table = &AppDataJinyan{}
var appdatajinyan_tablelist = []*AppDataJinyan{}

// type JinyanAppData struct {
// 	ID       uint64 `redis:"id" json:"id,string"`
// 	Account  string `redis:"account" json:"account"`
// 	Appname  string `redis:"appname" json:"appname"`
// 	Zonename string `redis:"zonename" json:"zonename"`
// 	Nickname string `redis:"nickname" json:"nickname"`

// 	Why        string    `redis:"why" json:"why"`
// 	Dateline   time.Time `redis:"dateline" json:"-"`
// 	Jinyandate time.Time `redis:"jinyandate" json:"-"`
// }

// func (ps *JinyanAppData) MarshalJSON() ([]byte, error) {
// 	// 定义一个该结构体的别名
// 	type Alias JinyanAppData
// 	// 定义一个新的结构体
// 	tmpSt := struct {
// 		Alias
// 		EndDate    string `json:"enddate"`
// 		CreateDate string `json:"createdate"`
// 	}{
// 		Alias:      (Alias)(*ps),
// 		EndDate:    ps.Dateline.Format("01/02/2006 15:04:05"),
// 		CreateDate: ps.Jinyandate.Format("01/02/2006 15:04:05"),
// 	}
// 	return json.Marshal(tmpSt)
// }

type AppDataJinyanFilter struct {
	Account  string
	Appname  string
	Zonename string
	Nickname string

	Startdateline *time.Time
	Enddateline   *time.Time

	Jinyanbegindate *time.Time
	Jinyanenddate   *time.Time
}

func (filter *AppDataJinyanFilter) apply(db *gorm.DB) *gorm.DB {
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

	if filter.Jinyanbegindate != nil {
		retdb = retdb.Where("jinyandate >= ?", *filter.Jinyanbegindate)
	}
	if filter.Jinyanenddate != nil {
		retdb = retdb.Where("jinyandate <= ?", *filter.Jinyanenddate)
	}

	return retdb
}

func (db *DBManager) GetAppDataJinyanCount(appname, zonename string, args ...*AppDataJinyanFilter) (uint64, error) {
	var count uint64
	retdb := db.sql.Table(db.sql.prefix + "app_data a")
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

	retdb = retdb.Joins("join " + db.sql.prefix + "app_data_jinyan b on b.dataid = a.id")

	retdb = retdb.Count(&count)
	return count, retdb.Error
}

func (db *DBManager) GetAppDataJinyanList(appname, zonename string, offset, count int, args ...*AppDataJinyanFilter) ([]*AppDataJinyan, error) {
	appdatalist := []*AppDataJinyan{}
	retdb := db.sql.Table(db.sql.prefix + "app_data a").Offset(offset).Limit(count)
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

	retdb = retdb.Joins("join " + db.sql.prefix + "app_data_jinyan b on b.dataid = a.id")

	retdb = retdb.Select("a.*, b.*").Scan(&appdatalist)
	return appdatalist, retdb.Error
}

func (db *DBManager) GetAppDataJinyan(id uint64) (*AppDataJinyan, error) {
	appdata := &AppDataJinyan{}
	retdb := db.sql.Table(db.sql.prefix + "app_data a")
	retdb = retdb.Joins("join "+db.sql.prefix+"app_data_jinyan b on b.dataid = a.id").Where("dataid = ?", id)
	retdb = retdb.Select("a.*, b.*").Limit(1).Scan(appdata)
	return appdata, retdb.Error
}

func (db *DBManager) JinyanAppDatas(ids []uint64, tbl_appdatajinyan *AppDataJinyan) error {
	tx := db.sql.Begin()
	for _, id := range ids {
		tbl_appdatajinyan.Dataid = id
		if err := tx.Create(tbl_appdatajinyan).Error; err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Model(appdata_table).Where("id = ?", id).Update("isjinyan", true).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func (db *DBManager) UnJinyanAppDatas(ids []uint64) error {
	tx := db.sql.Begin()
	for _, id := range ids {
		if err := tx.Model(appdata_table).Where("id = ?", id).Update("isjinyan", false).Error; err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Delete(appdatajinyan_table).Where("dataid = ?", id).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func (db *DBManager) JinyanAppData(tbl_appdatajinyan *AppDataJinyan) error {
	tx := db.sql.Begin()
	if err := tx.Create(tbl_appdatajinyan).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(appdata_table).Where("id = ?", tbl_appdatajinyan.Dataid).Update("isjinyan", true).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (db *DBManager) UnJinyanAppData(id uint64) error {
	tx := db.sql.Begin()
	if err := tx.Delete(appdatajinyan_table).Where("dataid = ?", id).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(appdata_table).Where("id = ?", id).Update("isjinyan", false).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
