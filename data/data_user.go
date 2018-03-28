package gtdata

import (
	"time"

	. "github.com/gtechx/base/common"
	"github.com/gtechx/chatserver/config"
)

//每个app之间可以是独立的数据，也可以共享数据，根据你的设置
func (rdm *RedisDataManager) CreateAccount(account, password, regip string) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("INCR", "UID")

	if err != nil {
		return err
	}

	uid := Uint64(ret)

	conn.Send("MULTI")
	conn.Send("HMSET", uid, "account", account, "password", password, "regip", regip, "regdate", time.Now().Unix())
	conn.Send("HSET", "account:uid", account, uid)
	conn.Send("SADD", "user", uid)

	_, err = conn.Do("EXEC")
	return err
}

func (rdm *RedisDataManager) CreateAppData(entity *EntityKey) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()

	conn.Send("MULTI")
	conn.Send("HSET", entity.KeyAppData, "createdate", time.Now().Unix())
	conn.Send("SADD", entity.KeyGroup, defaultGroupName)
	_, err := conn.Do("EXEC")
	return err
}

func (rdm *RedisDataManager) DeleteAppData(entity *EntityKey) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()

	conn.Send("MULTI")
	conn.Send("HDEL", entity.KeyAppData)
	conn.Send("DEL", entity.KeyGroup)
	conn.Send("HDEL", entity.KeyFriend)
	conn.Send("HDEL", entity.KeyFriendRequest)
	conn.Send("HDEL", entity.KeyBlack)
	_, err := conn.Do("EXEC")
	return err
}

func (rdm *RedisDataManager) IsAppDataExists(entity *EntityKey) (bool, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("HEXISTS", entity.KeyAppData)
	return Bool(ret), err
}

func (rdm *RedisDataManager) SetAppDataConfig(entity *EntityKey, configname string, data interface{}) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("HSET", entity.KeyAppData, configname)
	return err
}

func (rdm *RedisDataManager) GetAppDataConfig(entity *EntityKey, configname string) (interface{}, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("HGET", entity.KeyAppData, configname)
	return ret, err
}

// func (rdm *RedisDataManager) SetMaxFriends(uid uint64, count int) error {
// 	conn := rdm.redisPool.Get()
// 	defer conn.Close()
// 	_, err := conn.Do("HSET", uid, "maxfriends", count)
// 	return err
// }

// func (rdm *RedisDataManager) SetDesc(uid uint64, desc string) error {
// 	conn := rdm.redisPool.Get()
// 	defer conn.Close()
// 	_, err := conn.Do("HSET", uid, "desc", desc)
// 	return err
// }

func (rdm *RedisDataManager) IsAccountExists(account string) (bool, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("HEXISTS", "account:uid", account)
	return Bool(ret), err
}

func (rdm *RedisDataManager) IsUIDExists(uid uint64) (bool, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("EXISTS", uid)
	return Bool(ret), err
}

func (rdm *RedisDataManager) GetUIDByAccount(account string) (uint64, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("HGET", "account:uid", account)
	return Uint64(ret), err
}

func (rdm *RedisDataManager) GetAccountByUID(uid uint64) (string, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("HGET", uid, "account")
	return String(ret), err
}

func (rdm *RedisDataManager) GetPassword(uid uint64) (string, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("HGET", uid, "password")
	return String(ret), err
}

func (rdm *RedisDataManager) SetUserOnline(entity *EntityKey) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	conn.Send("MULTI")
	conn.Send("HSET", entity.KeyAppData, "online", config.ServerAddr)
	conn.Send("SADD", "online", entity.KeyUID)
	_, err := conn.Do("EXEC")
	return err
}

func (rdm *RedisDataManager) SetUserOffline(entity *EntityKey) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	conn.Send("MULTI")
	conn.Send("HDEL", entity.KeyAppData, "online", config.ServerAddr)
	conn.Send("SREM", "online", entity.KeyUID)
	_, err := conn.Do("EXEC")
	return err
}

func (rdm *RedisDataManager) IsUserOnline(uid uint64) (bool, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("SISMEMBER", "online", uid)
	return Bool(ret), err
}

func (rdm *RedisDataManager) SetUserState(entity *EntityKey, state uint8) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("HSET", entity.KeyAppData, "state", state)
	return err
}

func (rdm *RedisDataManager) AddUserToBlack(entity *EntityKey, otheruid uint64) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("SADD", entity.KeyBlack, otheruid)
	return err
}

func (rdm *RedisDataManager) RemoveUserFromBlack(entity *EntityKey, otheruid uint64) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("SREM", entity.KeyBlack, otheruid)
	return err
}

func (rdm *RedisDataManager) IsUserInBlack(entity *EntityKey, otheruid uint64) (bool, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("SISMEMBER", entity.KeyBlack, otheruid)
	return Bool(ret), err
}
