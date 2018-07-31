package gtdb

// import (
// 	//"errors"

// 	//. "github.com/gtechx/base/common"
// )
//怎样更新房间内玩家的在线状态？
//让客户端自己去请求更新房间玩家在线状态。
//就是客户端打开房间查看玩家列表的时候，才需要更新房间玩家在线状态。
var room_table = &Room{}
var room_tablelist = []*Room{}

var roomuser_table = &RoomUser{}
var roomuser_tablelist = []*RoomUser{}

//room op
func (db *DBManager) CreateRoom(tbl_room *Room) error {
	tx := db.sql.Begin()
	if err := tx.Create(tbl_room).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Create(&RoomUser{Rid: tbl_room.Rid, Dataid: tbl_room.Ownerid, Isowner: true, Isadmin: true}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (db *DBManager) DeleteRoom(rid uint64) error {
	tx := db.sql.Begin()
	if err := tx.Delete(&Room{Rid: rid}, "rid = ?", rid).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Delete(&RoomUser{Rid: rid}, "rid = ?", rid).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (db *DBManager) GetRoom(rid uint64) (*Room, error) {
	room := &Room{}
	retdb := db.sql.Where("rid = ?", rid).First(room)
	return room, retdb.Error
}

func (db *DBManager) GetRoomListByOwner(appdataid uint64) ([]*Room, error) {
	roomlist := []*Room{}
	retdb := db.sql.Where("dataid = ?", appdataid).Find(roomlist)
	return roomlist, retdb.Error
}

func (db *DBManager) GetRoomCountByOwner(appdataid uint64) (uint64, error) {
	var count uint64
	retdb := db.sql.Model(room_table).Where("dataid = ?", appdataid).Count(&count)
	return count, retdb.Error
}

func (db *DBManager) GetRoomListByJoined(appdataid uint64) ([]*Room, error) {
	roomlist := []*Room{}
	retdb := db.sql.Table("gtchat_rooms a")
	retdb = retdb.Joins("join gtchat_room_users b on b.rid = a.rid").Where("dataid = ?", appdataid)
	retdb = retdb.Select("a.*").Scan(roomlist)
	return roomlist, retdb.Error
}

func (db *DBManager) GetRoomCountByJoined(appdataid uint64) (uint64, error) {
	var count uint64
	retdb := db.sql.Table("gtchat_rooms a")
	retdb = retdb.Joins("join gtchat_room_users b on b.rid = a.rid").Where("dataid = ?", appdataid)
	retdb = retdb.Select("a.*").Count(&count)
	return count, retdb.Error
}

//room user op
func (db *DBManager) AddRoomUser(tbl_roomuser *RoomUser) error {
	retdb := db.sql.Create(tbl_roomuser)
	return retdb.Error
}

func (db *DBManager) RemoveRoomUser(rid, appdataid uint64) error {
	retdb := db.sql.Delete(&RoomUser{}, "rid = ? and dataid = ?", rid, appdataid)
	return retdb.Error
}

func (db *DBManager) GetRoomUser(rid, appdataid uint64) (*RoomUser, error) {
	roomuser := &RoomUser{}
	retdb := db.sql.Where("rid = ?", rid).Where("dataid = ?", appdataid).First(roomuser)
	return roomuser, retdb.Error
}

func (db *DBManager) GetRoomUserList(rid uint64) ([]*RoomUser, error) {
	roomuserlist := []*RoomUser{}
	retdb := db.sql.Where("rid = ?", rid).Find(roomuserlist)
	return roomuserlist, retdb.Error
}

func (db *DBManager) GetRoomUserCount(rid uint64) (uint64, error) {
	var count uint64
	retdb := db.sql.Model(roomuser_table).Where("rid = ?", rid).Count(&count)
	return count, retdb.Error
}

func (db *DBManager) GetRoomList() {

}

func (db *DBManager) getRoomType() {

}

func (db *DBManager) getRoomPassword() {

}

func (db *DBManager) setRoomPassword() {

}

func (db *DBManager) isRoomExist() {

}

func (db *DBManager) isUserInRoom() {

}

func (db *DBManager) addUserToRoom() {

}

func (db *DBManager) removeUserToRoom() {

}

//踢出玩家
func (db *DBManager) banUserInRoom() {

}

func (db *DBManager) JinyanUserInRoom() {

}

func (db *DBManager) UnJinyanUserInRoom() {

}

func (db *DBManager) setRoomDescription() {

}

func (db *DBManager) getRoomDescription() {

}

func (db *DBManager) setRoomVerify() {

}

func (db *DBManager) getRoomVerify() {

}

func (db *DBManager) setRoomVerifyType() {

}

func (db *DBManager) getRoomVerifyType() {

}
