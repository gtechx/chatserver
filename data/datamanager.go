package gtdata

import (
	"time"

	"github.com/garyburd/redigo/redis"
	. "github.com/gtechx/base/common"
)

type RedisDataManager struct {
	redisPool *redis.Pool

	serverAddr     string
	serverPassword string
	defaultDB      uint64
}

var instanceDataManager *RedisDataManager

func Manager() *RedisDataManager {
	if instanceDataManager == nil {
		instanceDataManager = &RedisDataManager{}
	}
	return instanceDataManager
}

func (rdm *RedisDataManager) Initialize(saddr, spass string, defaultdb, startuid, startappid uint64) error {
	rdm.serverAddr = saddr
	rdm.serverPassword = spass
	rdm.defaultDB = defaultdb

	rdm.redisPool = &redis.Pool{
		MaxIdle:      3,
		IdleTimeout:  240 * time.Second,
		Dial:         rdm.redisDial,
		TestOnBorrow: rdm.redisOnBorrow,
	}

	conn := rdm.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("EXISTS", "UID")

	if err != nil {
		return err
	}

	if !Bool(ret) {
		_, err = conn.Do("SET", "UID", startuid)

		if err != nil {
			return err
		}
	}

	ret, err = conn.Do("EXISTS", "APPID")

	if err != nil {
		return err
	}

	if !Bool(ret) {
		_, err = conn.Do("SET", "APPID", startappid)

		if err != nil {
			return err
		}
	}

	ret, err = conn.Do("HEXISTS", "admin", 0)

	if err != nil {
		return err
	}

	if !Bool(ret) {
		_, err = conn.Do("HSET", "admin", 0, 0xffffffff)

		if err != nil {
			return err
		}

		_, err = conn.Do("HSET", "hset:user:account:admin", "password", Md5("ztgame@123"))

		if err != nil {
			return err
		}
	}

	return err
}

func (rdm *RedisDataManager) UnInitialize() error {
	var err error
	if rdm.redisPool != nil {
		err = rdm.redisPool.Close()
	}
	return err
}

func (rdm *RedisDataManager) redisDial() (redis.Conn, error) {
	c, err := redis.Dial("tcp", rdm.serverAddr)
	if err != nil {
		return nil, err
	}
	if rdm.serverPassword != "" {
		if _, err := c.Do("AUTH", rdm.serverPassword); err != nil {
			c.Close()
			return nil, err
		}
	}
	if _, err := c.Do("SELECT", rdm.defaultDB); err != nil {
		c.Close()
		return nil, err
	}
	return c, nil
}

func (rdm *RedisDataManager) redisOnBorrow(c redis.Conn, t time.Time) error {
	if time.Since(t) < time.Minute {
		return nil
	}
	_, err := c.Do("PING")
	return err
}
