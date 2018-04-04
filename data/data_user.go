package gtdata

import (
	"time"

	"github.com/garyburd/redigo/redis"
	. "github.com/gtechx/base/common"
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
	conn.Send("HMSET", "hset:user:account"+account, "uid", uid, "account", account, "password", password, "regip", regip, "regdate", time.Now().Unix())
	conn.Send("HSET", "hset:user:uid:account", uid, account)
	conn.Send("SADD", "set:user", account)

	_, err = conn.Do("EXEC")
	return err
}

func (rdm *RedisDataManager) SetPassword(account, password string) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("HSET", "hset:user:account"+account, "password", password)
	return err
}

func (rdm *RedisDataManager) CreateAppData(datakey *DataKey) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()

	conn.Send("MULTI")
	conn.Send("HSET", datakey.KeyAppDataHsetByAppidZonenameAccount, "regdate", time.Now().Unix())
	conn.Send("SADD", datakey.KeyAppDataSetGroupByAppidZonenameAccount, defaultGroupName)
	_, err := conn.Do("EXEC")
	return err
}

func (rdm *RedisDataManager) DeleteAppData(datakey *DataKey) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()

	conn.Send("MULTI")
	conn.Send("DEL", datakey.KeyAppDataHsetByAppidZonenameAccount)
	conn.Send("DEL", datakey.KeyAppDataSetGroupByAppidZonenameAccount)
	conn.Send("DEL", datakey.KeyAppDataHsetFriendByAppidZonenameAccount)
	conn.Send("DEL", datakey.KeyAppDataHsetFriendrequestGroupByAppidZonenameAccount)
	conn.Send("DEL", datakey.KeyAppDataSetBlackByAppidZonenameAccount)
	conn.Send("DEL", datakey.KeyAppDataListMsgByAppidZonenameAccount)
	_, err := conn.Do("EXEC")
	return err
}

func (rdm *RedisDataManager) IsAppDataExists(datakey *DataKey) (bool, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("EXISTS", datakey.KeyAppDataHsetByAppidZonenameAccount)
	return redis.Bool(ret, err)
}

func (rdm *RedisDataManager) SetAppDataField(datakey *DataKey, fieldname string, value interface{}) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("HSET", datakey.KeyAppDataHsetByAppidZonenameAccount, fieldname, value)
	return err
}

func (rdm *RedisDataManager) GetAppDataField(datakey *DataKey, fieldname string) (interface{}, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("HGET", datakey.KeyAppDataHsetByAppidZonenameAccount, fieldname)
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
	ret, err := conn.Do("SISMEMBER", "set:user", account)
	return redis.Bool(ret, err)
}

// func (rdm *RedisDataManager) IsUIDExists(uid uint64) (bool, error) {
// 	conn := rdm.redisPool.Get()
// 	defer conn.Close()
// 	ret, err := conn.Do("EXISTS", "uid:"+String(uid))
// 	return redis.Bool(ret, err)
// }

// func (rdm *RedisDataManager) GetUIDByAccount(account string) (uint64, error) {
// 	conn := rdm.redisPool.Get()
// 	defer conn.Close()
// 	ret, err := conn.Do("HGET", "account:uid", account)
// 	return redis.Uint64(ret, err)
// }

// func (rdm *RedisDataManager) GetAccountByUID(uid uint64) (string, error) {
// 	conn := rdm.redisPool.Get()
// 	defer conn.Close()
// 	ret, err := conn.Do("HGET", "uid:"+String(uid), "account")
// 	return redis.String(ret, err)
// }

func (rdm *RedisDataManager) GetPassword(account string) (string, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("HGET", "hset:user:account"+account, "password")
	return redis.String(ret, err)
}

func (rdm *RedisDataManager) SetUserOnline(datakey *DataKey) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	conn.Send("MULTI")
	conn.Send("HSET", datakey.KeyAppDataHsetByAppidZonenameAccount, "online", rdm.serverAddr)
	conn.Send("SADD", "online:"+datakey.Appname+":"+datakey.Zonename, datakey.Account)
	_, err := conn.Do("EXEC")
	return err
}

func (rdm *RedisDataManager) SetUserOffline(datakey *DataKey) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	conn.Send("MULTI")
	conn.Send("HDEL", datakey.KeyAppDataHsetByAppidZonenameAccount, "online", rdm.serverAddr)
	conn.Send("SREM", "online:"+datakey.Appname+":"+datakey.Zonename, datakey.Account)
	_, err := conn.Do("EXEC")
	return err
}

func (rdm *RedisDataManager) IsUserOnline(datakey *DataKey, account string) (bool, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("HEXISTS", "hset:app:data:"+datakey.Appname+":"+datakey.Zonename+":"+datakey.Account, "online")
	return redis.Bool(ret, err)
}

func (rdm *RedisDataManager) SetUserState(datakey *DataKey, state uint8) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("HSET", datakey.KeyAppDataHsetByAppidZonenameAccount, "state", state)
	return err
}

func (rdm *RedisDataManager) AddUserToBlack(datakey *DataKey, otheraccount string) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("SADD", datakey.KeyAppDataSetBlackByAppidZonenameAccount, otheraccount)
	return err
}

func (rdm *RedisDataManager) RemoveUserFromBlack(datakey *DataKey, otheraccount string) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("SREM", datakey.KeyAppDataSetBlackByAppidZonenameAccount, otheraccount)
	return err
}

func (rdm *RedisDataManager) IsUserInBlack(datakey *DataKey, otheraccount string) (bool, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("SISMEMBER", datakey.KeyAppDataSetBlackByAppidZonenameAccount, otheraccount)
	return redis.Bool(ret, err)
}
