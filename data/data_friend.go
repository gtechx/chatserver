package gtdata

import (
	"github.com/garyburd/redigo/redis"
	. "github.com/gtechx/base/common"
)

var defaultGroupName string = "我的好友"
var userOnlineKeyName string = "user:online"

func (rdm *RedisDataManager) AddFriendRequest(entity *EntityKey, otheruid uint64, group string) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("HSET", entity.KeyFriendRequest, otheruid, group)
	return err
}

func (rdm *RedisDataManager) RemoveFriendRequest(entity *EntityKey, otheruid uint64) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("HDEL", entity.KeyFriendRequest, otheruid)
	return err
}

func (rdm *RedisDataManager) AddFriend(entity *EntityKey, otheruid uint64, group string) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()

	conn.Send("MULTI")
	conn.Send("HSET", entity.KeyFriend, otheruid, group)
	conn.Send("SADD", entity.KeyGroup+":"+group, otheruid)
	_, err := conn.Do("EXEC")

	return err
}

func (rdm *RedisDataManager) RemoveFriend(entity *EntityKey, otheruid uint64, group string) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()

	conn.Send("MULTI")
	conn.Send("HDEL", entity.KeyFriend, otheruid)
	conn.Send("SREM", entity.KeyGroup+":"+group, otheruid)
	_, err := conn.Do("EXEC")

	return err
}

func (rdm *RedisDataManager) GetFriendList(entity *EntityKey, group string) ([]uint64, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("SMEMBERS", entity.KeyGroup+":"+group)

	if err != nil {
		return nil, err
	}

	retarr, err := redis.Values(ret, nil)

	if err != nil {
		return nil, err
	}

	userlist := []uint64{}
	for _, uid := range retarr {
		userlist = append(userlist, Uint64(uid))
	}

	return userlist, err
}

func (rdm *RedisDataManager) IsFriend(entity *EntityKey, otheruid uint64) (bool, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("HEXISTS", entity.KeyFriend, otheruid)
	return redis.Bool(ret, err)
}

func (rdm *RedisDataManager) GetGroupOfFriend(entity *EntityKey, otheruid uint64) (string, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("HGET", entity.KeyFriend, otheruid)
	return redis.String(ret, err)
}

func (rdm *RedisDataManager) AddGroup(entity *EntityKey, group string) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("SADD", entity.KeyGroup, group)
	return err
}

func (rdm *RedisDataManager) RemoveGroup(entity *EntityKey, group string) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("SREM", entity.KeyGroup, group)
	return err
}

func (rdm *RedisDataManager) GetGroupList(entity *EntityKey) ([]string, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("SMEMBERS", entity.KeyGroup)

	if err != nil {
		return nil, err
	}

	retarr, err := redis.Values(ret, nil)

	if err != nil {
		return nil, err
	}

	grouplist := []string{}
	for _, group := range retarr {
		grouplist = append(grouplist, String(group))
	}

	return grouplist, err
}

func (rdm *RedisDataManager) IsGroupExists(entity *EntityKey, group string) (bool, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("SISMEMBER", entity.KeyGroup, group)
	return redis.Bool(ret, err)
}

func (rdm *RedisDataManager) IsFriendInGroup(entity *EntityKey, otheruid uint64, group string) (bool, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("SISMEMBER", entity.KeyGroup+":"+group, otheruid)
	return redis.Bool(ret, err)
}

func (rdm *RedisDataManager) MoveFriendToGroup(entity *EntityKey, otheruid uint64, srcgroup, destgroup string) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()

	conn.Send("MULTI")
	conn.Send("SREM", entity.KeyGroup+":"+srcgroup, otheruid)
	conn.Send("SADD", entity.KeyGroup+":"+destgroup, otheruid)
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
