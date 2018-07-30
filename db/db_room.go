package gtdb

// import (
// 	//"errors"

// 	//. "github.com/gtechx/base/common"
// )
var room_table = &Room{}
var room_tablelist = []*Room{}

//room op
func (db *DBManager) CreateRoom(tbl_room *Room) error {
	retdb := db.sql.Create(tbl_room)
	return retdb.Error
}

func (db *DBManager) DeleteRoom(rid uint64) error {
	retdb := db.sql.Delete(&Room{Rid: rid}, "rid = ?", rid)
	return retdb.Error
}

func (db *DBManager) getRoomList() {

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
