package gtdb

//. "github.com/gtechx/base/common"

var defaultGroupName string = "我的好友"
var userOnlineKeyName string = "user:online"

var friend_table = &Friend{}
var friend_tablelist = []*Friend{}

var group_table = &Group{}
var group_tablelist = []*Group{}

// func (db *DBManager) AddFriendRequest(id, otherid uint64, group string) error {
// 	conn := db.rd.Get()
// 	defer conn.Close()
// 	_, err := conn.Do("HSET", keyJoin("hset:app:data:friend:request", id), otherid, group)
// 	return err
// }

// func (db *DBManager) RemoveFriendRequest(id, otherid uint64) error {
// 	conn := db.rd.Get()
// 	defer conn.Close()
// 	_, err := conn.Do("HDEL", keyJoin("hset:app:data:friend:request", id), otherid)
// 	return err
// }

func (db *DBManager) AddFriend(tbl_from, tbl_to *Friend) error {
	// retdb := db.sql.Create(tbl_friend)
	// return retdb.Error
	tx := db.sql.Begin()
	if err := tx.Create(tbl_from).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Create(tbl_to).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (db *DBManager) RemoveFriend(id, otherid uint64) error {
	tx := db.sql.Begin()
	if err := tx.Delete(friend_table, "dataid = ? AND otherdataid = ?", id, otherid).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Delete(friend_table, "dataid = ? AND otherdataid = ?", otherid, id).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (db *DBManager) GetFriend(id, otherid uint64) (*Friend, error) {
	friend := &Friend{}
	retdb := db.sql.Where("dataid = ? AND otherdataid = ?", id, otherid).First(friend)
	return friend, retdb.Error
}

func (db *DBManager) GetFriendOnlineList(id uint64) ([]*Online, error) {
	onlinelist := []*Online{}
	retdb := db.sql.Model(online_table).Joins("join gtchat_friends on gtchat_friends.dataid = ? AND gtchat_friends.otherdataid = gtchat_onlines.dataid", id)
	retdb = retdb.Find(&onlinelist)
	return onlinelist, retdb.Error
}

func (db *DBManager) GetOnlineFriendIdList(id uint64) ([]uint64, error) {
	var friendidlist []uint64
	retdb := db.sql.Table("gtchat_friends").Where("gtchat_friends.dataid = ?", id).Select("gtchat_friends.otherdataid").Joins("join gtchat_onlines on gtchat_friends.otherdataid = gtchat_onlines.dataid").Scan(&friendidlist)
	return friendidlist, retdb.Error
}

func (db *DBManager) GetOfflineFriendIdList(id uint64) ([]uint64, error) {
	var friendidlist []uint64
	retdb := db.sql.Table("gtchat_friends").Where("gtchat_friends.dataid = ?", id).Select("gtchat_friends.otherdataid").Joins("join gtchat_onlines on gtchat_friends.otherdataid != gtchat_onlines.dataid").Scan(&friendidlist)
	return friendidlist, retdb.Error
}

func (db *DBManager) GetFriendIdList(id uint64) ([]uint64, error) {
	var friendidlist []uint64
	retdb := db.sql.Table("gtchat_friends").Where("gtchat_friends.dataid = ?", id).Pluck("otherdataid", &friendidlist) //.Select("friends.otherdataid").Scan(&friendidlist)
	return friendidlist, retdb.Error
}

func (db *DBManager) GetFriendInfoList(id uint64) ([]*FriendJson, error) {
	friendlist := []*FriendJson{}
	retdb := db.sql.Table("gtchat_friends").Where("gtchat_friends.dataid = ?", id).Select("gtchat_friends.otherdataid as dataid, gtchat_friends.group, gtchat_friends.comment, gtchat_app_data.nickname, gtchat_app_data.desc").Joins("join gtchat_app_data on gtchat_friends.otherdataid = gtchat_app_data.id").Find(&friendlist)
	return friendlist, retdb.Error
}

func (db *DBManager) GetFriendList(id uint64, offset, count int) ([]*Friend, error) {
	friendlist := []*Friend{}
	retdb := db.sql.Offset(offset).Limit(count).Where("dataid = ?", id).Find(&friendlist)
	return friendlist, retdb.Error
}

func (db *DBManager) GetFriendListByGroup(id uint64, group string) ([]*Friend, error) {
	friendlist := []*Friend{}
	retdb := db.sql.Where("dataid = ? AND group = ?", id, group).Find(&friendlist)
	return friendlist, retdb.Error
}

func (db *DBManager) IsFriend(id, otherid uint64) (bool, error) {
	var count int
	retdb := db.sql.Model(friend_table).Where("dataid = ? AND otherdataid = ?", id, otherid).Count(&count)
	return count > 0, retdb.Error
}

// func (db *DBManager) GetGroupFriendIn(datakey *DataKey, otheraccount string) (string, error) {
// 	conn := db.redisPool.Get()
// 	defer conn.Close()
// 	ret, err := conn.Do("HGET", datakey.KeyAppDataHsetFriendByAppidZonenameAccount, otheraccount)
// 	return redis.String(ret, err)
// }

func (db *DBManager) AddGroup(tbl_group *Group) error {
	retdb := db.sql.Create(tbl_group)
	return retdb.Error
}

func (db *DBManager) RemoveGroup(id uint64, group string) error {
	retdb := db.sql.Delete(group_table, "dataid = ? AND group = ?", id, group)
	return retdb.Error
}

func (db *DBManager) GetGroupList(id uint64) ([]*Group, error) {
	grouplist := []*Group{}
	retdb := db.sql.Where("dataid = ?", id).Find(&grouplist)
	return grouplist, retdb.Error
}

func (db *DBManager) IsGroupExists(id uint64, group string) (bool, error) {
	var count int
	retdb := db.sql.Model(group_table).Where("dataid = ? AND group = ?", id, group).Count(&count)
	return count > 0, retdb.Error
}

func (db *DBManager) IsInGroup(id, otherid uint64, group string) (bool, error) {
	var count int
	retdb := db.sql.Model(friend_table).Where("dataid = ? AND otherdataid = ? AND group = ?", id, otherid, group).Count(&count)
	return count > 0, retdb.Error
}

func (db *DBManager) MoveToGroup(id, otherid uint64, destgroup string) error {
	retdb := db.sql.Model(friend_table).Where("dataid = ? AND otherdataid = ?", id, otherid).Update("group", destgroup)
	return retdb.Error
}

func (db *DBManager) SetComment(id, otherid uint64, comment string) error {
	retdb := db.sql.Model(friend_table).Where("dataid = ? AND otherdataid = ?", id, otherid).Update("comment", comment)
	return retdb.Error
}

// tx := retdb := db.sql.Begin()
// // 注意，一旦你在一个事务中，使用tx作为数据库句柄

// if err := retdb := db.sql.Create(&Animal{Name: "Giraffe"}).Error; err != nil {
// 	tx.Rollback()
// 	return err
// }

// if err := tx.Create(&Animal{Name: "Lion"}).Error; err != nil {
// 	tx.Rollback()
// 	return err
// }

// tx.Commit()
// return retdb.Error

// func (db *DBManager) BanFriend(uid, fuid uint64) {

// }

// func (db *DBManager) UnBanFriend(uid, fuid uint64) {

// }

// func (db *DBManager) SetFriendVerifyType(uid uint64, vtype byte) error {
// 	conn := db.redisPool.Get()
// 	defer conn.Close()
// 	_, err := conn.Do("HSET", uid, "verifytype", vtype)
// 	return err
// }

// func (db *DBManager) GetFriendVerifyType(uid uint64) (byte, error) {
// 	conn := db.redisPool.Get()
// 	defer conn.Close()

// 	ret, err := conn.Do("HGET", uid, "verifytype")

// 	return Byte(ret), err
// }
