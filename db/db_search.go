package gtdb

import "github.com/jinzhu/gorm"

func (db *DBManager) SearchUserById(id uint64) (*AppDataPublicWithID, error) {
	ret := &AppDataPublicWithID{}
	retdb := db.sql.Model(appdata_table).Where("id = ?", id).Scan(ret)
	if retdb.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return ret, retdb.Error
}

func (db *DBManager) SearchUserByNickname(nickname string) ([]*AppDataPublicWithID, error) {
	ret := []*AppDataPublicWithID{}
	retdb := db.sql.Model(appdata_table).Where("nickname like ?", "%"+nickname+"%").Scan(&ret)
	if retdb.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return ret, retdb.Error
}

func (db *DBManager) SearchRoom(roomname string) ([]*Room, error) {
	ret := []*Room{}
	retdb := db.sql.Model(room_table).Where("roomname like ?", "%"+roomname+"%").Scan(&ret)
	if retdb.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return ret, retdb.Error
}
