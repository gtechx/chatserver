package gtdata

import (
	"github.com/garyburd/redigo/redis"
	. "github.com/gtechx/base/common"
)

var defaultGroupName string = "我的好友"
var userOnlineKeyName string = "user:online"

func (rdm *RedisDataManager) AddFriendRequest(datakey *DataKey, otheraccount, group string) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("HSET", datakey.KeyAppDataHsetFriendrequestGroupByAppidZonenameAccount, otheraccount, group)
	return err
}

func (rdm *RedisDataManager) RemoveFriendRequest(datakey *DataKey, otheraccount string) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("HDEL", datakey.KeyAppDataHsetFriendrequestGroupByAppidZonenameAccount, otheraccount)
	return err
}

func (rdm *RedisDataManager) AddFriend(datakey *DataKey, otheraccount, group string) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	conn.Send("MULTI")
	conn.Send("HSET", datakey.KeyAppDataHsetFriendByAppidZonenameAccount, otheraccount, group)
	conn.Send("SADD", datakey.KeyAppDataHsetFriendByAppidZonenameAccount+":"+group, otheraccount)
	_, err := conn.Do("EXEC")
	return err
}

func (rdm *RedisDataManager) RemoveFriend(datakey *DataKey, otheraccount, group string) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	conn.Send("MULTI")
	conn.Send("HDEL", datakey.KeyAppDataHsetFriendByAppidZonenameAccount, otheraccount)
	conn.Send("SREM", datakey.KeyAppDataHsetFriendByAppidZonenameAccount+":"+group, otheraccount)
	_, err := conn.Do("EXEC")
	return err
}

func (rdm *RedisDataManager) GetFriendList(datakey *DataKey, group string) ([]string, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("SMEMBERS", datakey.KeyAppDataHsetFriendByAppidZonenameAccount+":"+group)

	retarr, err := redis.Values(ret, err)

	if err != nil {
		return nil, err
	}

	userlist := []string{}
	for _, account := range retarr {
		userlist = append(userlist, String(account))
	}

	return userlist, err
}

func (rdm *RedisDataManager) IsFriend(datakey *DataKey, otheraccount string) (bool, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("HEXISTS", datakey.KeyAppDataHsetFriendByAppidZonenameAccount, otheraccount)
	return redis.Bool(ret, err)
}

func (rdm *RedisDataManager) GetGroupFriendIn(datakey *DataKey, otheraccount string) (string, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("HGET", datakey.KeyAppDataHsetFriendByAppidZonenameAccount, otheraccount)
	return redis.String(ret, err)
}

func (rdm *RedisDataManager) AddGroup(datakey *DataKey, group string) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("SADD", datakey.KeyAppDataSetGroupByAppidZonenameAccount, group)
	return err
}

func (rdm *RedisDataManager) RemoveGroup(datakey *DataKey, group string) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("SREM", datakey.KeyAppDataSetGroupByAppidZonenameAccount, group)
	return err
}

func (rdm *RedisDataManager) GetGroupList(datakey *DataKey) ([]string, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("SMEMBERS", datakey.KeyAppDataSetGroupByAppidZonenameAccount)

	retarr, err := redis.Values(ret, err)

	if err != nil {
		return nil, err
	}

	grouplist := []string{}
	for _, group := range retarr {
		grouplist = append(grouplist, String(group))
	}

	return grouplist, err
}

func (rdm *RedisDataManager) IsGroupExists(datakey *DataKey, group string) (bool, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("SISMEMBER", datakey.KeyAppDataSetGroupByAppidZonenameAccount, group)
	return redis.Bool(ret, err)
}

func (rdm *RedisDataManager) IsFriendInGroup(datakey *DataKey, otheraccount, group string) (bool, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("SISMEMBER", datakey.KeyAppDataSetGroupByAppidZonenameAccount+":"+group, otheraccount)
	return redis.Bool(ret, err)
}

func (rdm *RedisDataManager) MoveFriendToGroup(datakey *DataKey, otheraccount, srcgroup, destgroup string) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()

	conn.Send("MULTI")
	conn.Send("SREM", datakey.KeyAppDataSetGroupByAppidZonenameAccount+":"+srcgroup, otheraccount)
	conn.Send("SADD", datakey.KeyAppDataSetGroupByAppidZonenameAccount+":"+destgroup, otheraccount)
	_, err := conn.Do("EXEC")

	return err
}

// func (rdm *RedisDataManager) BanFriend(uid, fuid uint64) {

// }

// func (rdm *RedisDataManager) UnBanFriend(uid, fuid uint64) {

// }

// func (rdm *RedisDataManager) SetFriendVerifyType(uid uint64, vtype byte) error {
// 	conn := rdm.redisPool.Get()
// 	defer conn.Close()
// 	_, err := conn.Do("HSET", uid, "verifytype", vtype)
// 	return err
// }

// func (rdm *RedisDataManager) GetFriendVerifyType(uid uint64) (byte, error) {
// 	conn := rdm.redisPool.Get()
// 	defer conn.Close()

// 	ret, err := conn.Do("HGET", uid, "verifytype")

// 	return Byte(ret), err
// }
