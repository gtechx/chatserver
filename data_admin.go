package main

import (
	"github.com/garyburd/redigo/redis"
	. "github.com/gtechx/base/common"
)

//[hashes]admin pair(uid, privilege) --管理员权限
//[sets]online uid --在线用户uid

func (rdm *RedisDataManager) IsAdmin(uid uint64) (bool, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("HEXISTS", "admin", uid)
	return Bool(ret), err
}

func (rdm *RedisDataManager) AddAdmin(uid uint64, privilege uint32) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("HSET", "admin", uid, privilege)
	return err
}

func (rdm *RedisDataManager) RemoveAdmin(uid, uuid uint64) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("HDEL", "admin", uuid)
	return err
}

func (rdm *RedisDataManager) GetPrivilege(uid uint64) (uint32, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("HGET", "admin", uid)
	return Uint32(ret), err
}

func (rdm *RedisDataManager) SetAdminPrivilege(uid uint64, privilege uint32) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("HSET", "admin", uid, privilege)
	return err
}

func (rdm *RedisDataManager) GetAdminList(uid uint64) ([]uint64, error) {
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

	adminlist := []uint64{}
	for _, uid := range retarr {
		adminlist = append(adminlist, Uint64(uid))
	}

	return adminlist, err
}

func (rdm *RedisDataManager) GetUserOnline() ([]uint64, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("SMEMBERS", "online")

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
