package gtdb

import (
	//"errors"

	"github.com/garyburd/redigo/redis"
	. "github.com/gtechx/base/common"
)

//[sorted sets]serverlist pair(count,addr)
//[sets]ttl:addr
var serverListKeyName string = "serverlist"

//server op
func (db *DBManager) RegisterServer(addr string) error {
	conn := db.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("ZADD", serverListKeyName, 0, addr)

	return err
}

func (db *DBManager) UnRegisterServer(addr string) error {
	conn := db.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("ZREM", serverListKeyName, addr)

	return err
}

func (db *DBManager) IncrByServerClientCount(addr string, count int) error {
	conn := db.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("ZINCRBY", serverListKeyName, count, addr)

	return err
}

func (db *DBManager) GetServerList() ([]string, error) {
	conn := db.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("ZRANGE", serverListKeyName, 0, -1)

	if err != nil {
		return nil, err
	}

	slist, _ := redis.Strings(ret, err)
	return slist, err
}

func (db *DBManager) GetServerCount() (int, error) {
	conn := db.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("ZCARD", serverListKeyName)

	count, err := redis.Uint64(ret, err)

	return Int(count), err
}

func (db *DBManager) SetServerTTL(addr string, seconds int) error {
	conn := db.redisPool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", "ttl:"+addr, 0, "EX", seconds)

	return err
}

func (db *DBManager) CheckServerTTL() error {
	return nil
}

func (db *DBManager) VoteServerDie() error {
	return nil
}
