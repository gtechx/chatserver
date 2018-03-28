package main

import (
	//"errors"

	"github.com/garyburd/redigo/redis"
	. "github.com/gtechx/base/common"
)

//[sorted sets]serverlist pair(count,addr)
//[sets]ttl:addr
var serverListKeyName string = "serverlist"

//server op
func (rdm *RedisDataManager) RegisterServer(addr string) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("ZADD", serverListKeyName, 0, addr)

	return err
}

func (rdm *RedisDataManager) UnRegisterServer(addr string) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("ZREM", serverListKeyName, addr)

	return err
}

func (rdm *RedisDataManager) IncrByServerClientCount(addr string, count int) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("ZINCRBY", serverListKeyName, count, addr)

	return err
}

func (rdm *RedisDataManager) GetServerList() ([]string, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("ZRANGE", serverListKeyName, 0, -1)

	if err != nil {
		return nil, err
	}

	slist, _ := redis.Strings(ret, err)
	return slist, err
}

func (rdm *RedisDataManager) GetServerCount() (int, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("ZCARD", serverListKeyName)

	return Int(ret), err
}

func (rdm *RedisDataManager) SetServerTTL(addr string, seconds int) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", "ttl:"+addr, 0, "EX", seconds)

	return err
}

func (rdm *RedisDataManager) CheckServerTTL() error {
	return nil
}

func (rdm *RedisDataManager) VoteServerDie() error {
	return nil
}
