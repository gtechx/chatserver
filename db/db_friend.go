package gtdb

import (
	"github.com/garyburd/redigo/redis"
	. "github.com/gtechx/base/common"
)

var defaultGroupName string = "我的好友"
var userOnlineKeyName string = "user:online"

friend_table := &Friend{}
friend_tablelist := []*Friend{}

group_table := &Group{}
group_tablelist := []*Group{}

func (db *DBManager) AddFriendRequest(appname, zonename, account, otheraccount, group string) error {
	conn := db.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("HSET", keyJoin("hset:app:data:friend:request", appname, zonename, account), otheraccount, group)
	return err
}

func (db *DBManager) RemoveFriendRequest(appname, zonename, account, otheraccount string) error {
	conn := db.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("HDEL", keyJoin("hset:app:data:friend:request", appname, zonename, account), otheraccount)
	return err
}

func (db *DBManager) AddFriend(tbl_friend *Friend) error {
	db.sql.Create(tbl_friend)
	return db.sql.Error
}

func (db *DBManager) RemoveFriend(appname, zonename, account, otheraccount string) error {
	db.sql.Delete(friend_table, "appname = ? AND zonename = ? AND account = ? AND otheraccount = ?", appname, zonename, account, otheraccount)
	return db.sql.Error
}

func (db *DBManager) GetFriend(appname, zonename, account, otheraccount string) (*Friend, error) {
	friend := &Friend{}
	db.sql.Where("appname = ? AND zonename = ? AND account = ? AND otheraccount = ?", appname, zonename, account, otheraccount).First(friend)
	return friend, db.sql.Error
}

func (db *DBManager) GetFriendList(appname, zonename, account string, offset, count int) ([]*Friend, error) {
	friendlist := []*Friend{}
	db.sql.Offset(offset).Limit(count).Where("appname = ? AND zonename = ? AND account = ?", appname, zonename, account).Find(&friendlist)
	return friendlist, err
}

func (db *DBManager) GetFriendListByGroup(appname, zonename, account, group string) ([]*Friend, error) {
	friendlist := []*App{}
	db.sql.Where("appname = ? AND zonename = ? AND account = ? AND group = ?", appname, zonename, account, group).Find(&friendlist)
	return friendlist, err
}

func (db *DBManager) IsFriend(appname, zonename, account, otheraccount string) (bool, error) {
	var count int
	db.sql.Model(friend_table).Where("appname = ? AND zonename = ? AND account = ? AND otheraccount = ?", appname, zonename, account, otheraccount).Count(&count)
	return count > 0, db.sql.Error
}

// func (db *DBManager) GetGroupFriendIn(datakey *DataKey, otheraccount string) (string, error) {
// 	conn := db.redisPool.Get()
// 	defer conn.Close()
// 	ret, err := conn.Do("HGET", datakey.KeyAppDataHsetFriendByAppidZonenameAccount, otheraccount)
// 	return redis.String(ret, err)
// }

func (db *DBManager) AddGroup(tbl_group *Group) error {
	db.sql.Create(tbl_group)
	return db.sql.Error
}

func (db *DBManager) RemoveGroup(appname, zonename, account, group string) error {
	db.sql.Delete(group_table, "appname = ? AND zonename = ? AND account = ? AND name = ?", appname, zonename, account, group)
	return db.sql.Error
}

func (db *DBManager) GetGroupList(appname, zonename, account string) ([]*Group, error) {
	grouplist := []*Group{}
	db.sql.Where("appname = ? AND zonename = ? AND account = ?", appname, zonename, account).Find(&grouplist)
	return grouplist, err
}

func (db *DBManager) IsGroupExists(appname, zonename, account, group string) (bool, error) {
	var count int
	db.sql.Model(group_table).Where("appname = ? AND zonename = ? AND account = ? AND name = ?", appname, zonename, account, group).Count(&count)
	return count > 0, db.sql.Error
}

func (db *DBManager) IsFriendInGroup(appname, zonename, account, otheraccount, group string) (bool, error) {
	var count int
	db.sql.Model(friend_table).Where("appname = ? AND zonename = ? AND account = ? AND otheraccount = ? AND group = ?", appname, zonename, account, group).Count(&count)
	return count > 0, db.sql.Error
}

func (db *DBManager) MoveFriendToGroup(appname, zonename, account, srcgroup, destgroup string) error {
	db.Model(friend_table).Where("appname = ? AND zonename = ? AND account = ? AND otheraccount = ? AND group = ?", appname, zonename, account, srcgroup).Update("group", destgroup)
	return db.sql.Error
}

// tx := db.sql.Begin()
// // 注意，一旦你在一个事务中，使用tx作为数据库句柄

// if err := db.sql.Create(&Animal{Name: "Giraffe"}).Error; err != nil {
// 	tx.Rollback()
// 	return err
// }

// if err := tx.Create(&Animal{Name: "Lion"}).Error; err != nil {
// 	tx.Rollback()
// 	return err
// }

// tx.Commit()
// return db.sql.Error

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
