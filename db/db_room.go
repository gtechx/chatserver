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

func (db *DBManager) UpdateRoom(tbl_room *Room) error {
	retdb := db.sql.Save(tbl_room)
	return retdb.Error
}

func (db *DBManager) IsRoomExists(rid uint64) (bool, error) {
	var count uint64
	retdb := db.sql.Model(room_table).Where("rid = ?", rid).Count(&count)
	return count > 0, retdb.Error
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
	retdb := db.sql.Where("ownerid = ?", appdataid).Find(roomlist)
	return roomlist, retdb.Error
}

func (db *DBManager) GetRoomCountByOwner(appdataid uint64) (uint64, error) {
	var count uint64
	retdb := db.sql.Model(room_table).Where("ownerid = ?", appdataid).Count(&count)
	return count, retdb.Error
}

func (db *DBManager) GetRoomListByJoined(appdataid uint64) ([]*Room, error) {
	roomlist := []*Room{}
	retdb := db.sql.Table(db.sql.prefix + "room a")
	retdb = retdb.Joins("join "+db.sql.prefix+"room_user b on b.rid = a.rid").Where("dataid = ?", appdataid)
	retdb = retdb.Select("a.*, b.msgsetting").Scan(roomlist)
	return roomlist, retdb.Error
}

func (db *DBManager) GetRoomCountByJoined(appdataid uint64) (uint64, error) {
	var count uint64
	retdb := db.sql.Table(db.sql.prefix + "room a")
	retdb = retdb.Joins("join "+db.sql.prefix+"room_user b on b.rid = a.rid").Where("dataid = ?", appdataid)
	retdb = retdb.Select("a.*").Count(&count)
	return count, retdb.Error
}

//room user op
func (db *DBManager) AddRoomUser(tbl_roomuser *RoomUser) error {
	retdb := db.sql.Create(tbl_roomuser)
	return retdb.Error
}

func (db *DBManager) UpdateRoomUser(tbl_roomuser *RoomUser) error {
	retdb := db.sql.Save(tbl_roomuser)
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
	retdb := db.sql.Table(db.sql.prefix + "app_data a")
	retdb = retdb.Joins("join "+db.sql.prefix+"room_user b on b.dataid = a.id").Where("rid = ?", rid)
	retdb = retdb.Joins("left join " + db.sql.prefix + "online c on c.dataid = a.id")
	retdb = retdb.Select("a.*, b.*, c.dataid is not null as isonline").Scan(roomuserlist)
	return roomuserlist, retdb.Error
}

func (db *DBManager) GetRoomUserIds(rid uint64) ([]uint64, error) {
	ids := []uint64{}
	retdb := db.sql.Model(roomuser_table).Select("a.id").Where("rid = ?", rid).Scan(&ids)
	return ids, retdb.Error
}

func (db *DBManager) GetRoomUserOnlineIds(rid uint64) ([]uint64, error) {
	ids := []uint64{}
	retdb := db.sql.Table(db.sql.prefix+"room_user a").Where("rid = ?", rid)
	retdb = retdb.Joins("join " + db.sql.prefix + "online b on b.dataid = a.dataid")
	retdb = retdb.Select("a.dataid").Scan(&ids)
	return ids, retdb.Error
}

func (db *DBManager) GetRoomUserCount(rid uint64) (uint64, error) {
	var count uint64
	retdb := db.sql.Model(roomuser_table).Where("rid = ?", rid).Count(&count)
	return count, retdb.Error
}

func (db *DBManager) IsUserInRoom(rid, appdataid uint64) (bool, error) {
	var count uint64
	retdb := db.sql.Model(roomuser_table).Where("rid = ?", rid).Where("dataid = ?", appdataid).Count(&count)
	return count > 0, retdb.Error
}

func (db *DBManager) GetRoomAdminIds(rid uint64) ([]uint64, error) {
	ids := []uint64{}
	retdb := db.sql.Model(roomuser_table).Where("rid = ?", rid).Where("isadmin = 1").Scan(&ids)
	return ids, retdb.Error
}

func (db *DBManager) IsRoomAdmin(rid, appdataid uint64) (bool, error) {
	var isadmin bool
	retdb := db.sql.Model(roomuser_table).Select("isadmin").Where("rid = ?", rid).Where("dataid = ?", appdataid).Scan(&isadmin)
	return isadmin, retdb.Error
}

//踢出玩家
func (db *DBManager) BanUserInRoom(rid, appdataid uint64) error {
	retdb := db.sql.Delete(roomuser_table, "rid = ? AND dataid = ?", rid, appdataid)
	return retdb.Error
}

func (db *DBManager) JinyanUserInRoom(rid, appdataid uint64) error {
	retdb := db.sql.Model(roomuser_table).Where("rid = ?", rid).Where("dataid = ?", appdataid).Update("isjinyan", true)
	return retdb.Error
}

func (db *DBManager) UnJinyanUserInRoom(rid, appdataid uint64) error {
	retdb := db.sql.Model(roomuser_table).Where("rid = ?", rid).Where("dataid = ?", appdataid).Update("isjinyan", false)
	return retdb.Error
}

func (db *DBManager) AddRoomAdmin(rid, appdataid uint64) error {
	retdb := db.sql.Model(roomuser_table).Where("rid = ?", rid).Where("dataid = ?", appdataid).Update("isadmin", true)
	return retdb.Error
}

func (db *DBManager) RemoveRoomAdmin(rid, appdataid uint64) error {
	retdb := db.sql.Model(roomuser_table).Where("rid = ?", rid).Where("dataid = ?", appdataid).Update("isadmin", false)
	return retdb.Error
}

func (db *DBManager) SetRoomUserDisplayName(rid, appdataid uint64, displayname string) error {
	retdb := db.sql.Model(roomuser_table).Where("rid = ?", rid).Where("dataid = ?", appdataid).Update("displayname", displayname)
	return retdb.Error
}

func (db *DBManager) SetRoomMsgSetting(rid, appdataid uint64, msgsetting byte) error {
	retdb := db.sql.Model(roomuser_table).Where("rid = ?", rid).Where("dataid = ?", appdataid).Update("msgsetting", msgsetting)
	return retdb.Error
}

func (db *DBManager) GetRoomMsgSetting(rid, appdataid uint64) (byte, error) {
	var msgsetting byte
	retdb := db.sql.Model(roomuser_table).Select("msgsetting").Where("rid = ?", rid).Where("dataid = ?", appdataid).Scan(&msgsetting)
	return msgsetting, retdb.Error
}

//room property
func (db *DBManager) SetRoomNotice(rid uint64, notice string) error {
	retdb := db.sql.Model(room_table).Where("rid = ?", rid).Update("notice", notice)
	return retdb.Error
}

func (db *DBManager) GetRoomNotice(rid uint64) (string, error) {
	var notice string
	retdb := db.sql.Model(room_table).Select("notice").Where("rid = ?", rid).Scan(&notice)
	return notice, retdb.Error
}

func (db *DBManager) SetRoomName(rid uint64, roomname string) error {
	retdb := db.sql.Model(room_table).Where("rid = ?", rid).Update("roomname", roomname)
	return retdb.Error
}

func (db *DBManager) GetRoomName(rid uint64) (string, error) {
	var roomname string
	retdb := db.sql.Model(room_table).Select("roomname").Where("rid = ?", rid).Scan(&roomname)
	return roomname, retdb.Error
}

func (db *DBManager) SetRoomType(rid uint64, roomtype byte) error {
	retdb := db.sql.Model(room_table).Where("rid = ?", rid).Update("roomtype", roomtype)
	return retdb.Error
}

func (db *DBManager) GetRoomType(rid uint64) (byte, error) {
	var roomtype byte
	retdb := db.sql.Model(room_table).Select("roomtype").Where("rid = ?", rid).Scan(&roomtype)
	return roomtype, retdb.Error
}

func (db *DBManager) SetRoomPassword(rid uint64, password string) error {
	retdb := db.sql.Model(room_table).Where("rid = ?", rid).Update("password", password)
	return retdb.Error
}

func (db *DBManager) GetRoomPassword(rid uint64) (string, error) {
	var password string
	retdb := db.sql.Model(room_table).Select("password").Where("rid = ?", rid).Scan(&password)
	return password, retdb.Error
}

func (db *DBManager) SetRoomJoinType(rid uint64, jointype byte) error {
	retdb := db.sql.Model(room_table).Where("rid = ?", rid).Update("jointype", jointype)
	return retdb.Error
}

func (db *DBManager) GetRoomJoinType(rid uint64) (byte, error) {
	var jointype byte
	retdb := db.sql.Model(room_table).Select("jointype").Where("rid = ?", rid).Scan(&jointype)
	return jointype, retdb.Error
}

func (db *DBManager) SetRoomJieshao(rid uint64, jieshao string) error {
	retdb := db.sql.Model(room_table).Where("rid = ?", rid).Update("jieshao", jieshao)
	return retdb.Error
}

func (db *DBManager) GetRoomJieshao(rid uint64) (string, error) {
	var jieshao string
	retdb := db.sql.Model(room_table).Select("jieshao").Where("rid = ?", rid).Scan(&jieshao)
	return jieshao, retdb.Error
}
