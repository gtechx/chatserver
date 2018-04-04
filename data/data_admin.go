package gtdata

import (
	"github.com/garyburd/redigo/redis"
	. "github.com/gtechx/base/common"
)

//[hashes]admin pair(uid, privilege) --管理员权限
//[sets]online uid --在线用户uid

func (rdm *RedisDataManager) IsAdmin(account string) (bool, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("HEXISTS", "admin", account)
	return redis.Bool(ret, err)
}

func (rdm *RedisDataManager) AddAdmin(account string, privilege uint64) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("HSET", "admin", account, privilege)
	return err
}

func (rdm *RedisDataManager) RemoveAdmin(account string) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("HDEL", "admin", account)
	return err
}

func (rdm *RedisDataManager) GetPrivilege(account string) (uint64, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("HGET", "admin", account)
	return redis.Uint64(ret, err)
}

func (rdm *RedisDataManager) SetAdminPrivilege(account string, privilege uint64) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("HSET", "admin", account, privilege)
	return err
}

func (rdm *RedisDataManager) GetAdminList() ([]string, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("HKEYS", "admin")

	if err != nil {
		return nil, err
	}

	retarr, err := redis.Values(ret, nil)

	if err != nil {
		return nil, err
	}

	adminlist := []string{}
	for _, account := range retarr {
		adminlist = append(adminlist, String(account))
	}

	return adminlist, err
}

func (rdm *RedisDataManager) GetUserOnline(appname, zonename string) ([]string, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("SMEMBERS", "online:"+appname+":"+zonename)

	if err != nil {
		return nil, err
	}

	retarr, err := redis.Values(ret, nil)

	if err != nil {
		return nil, err
	}

	userlist := []string{}
	for _, account := range retarr {
		userlist = append(userlist, String(account))
	}

	return userlist, err
}
