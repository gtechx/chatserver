package gtdb

import "github.com/jinzhu/gorm"

func (db *DBManager) SearchUserById(id uint64) (*SearchUserJson, error) {
	ret := &SearchUserJson{}
	retdb := db.sql.Model(appdata_table).Select("id as dataid, nickname, country").Where("id = ?", id).Find(ret)
	if retdb.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return ret, retdb.Error
}

func (db *DBManager) SearchUserByNickname(nickname string) ([]*SearchUserJson, error) {
	ret := []*SearchUserJson{}
	retdb := db.sql.Model(appdata_table).Select("id as dataid, nickname, country").Where("nickname = ?", "%"+nickname+"%").Find(ret)
	return ret, retdb.Error
}
