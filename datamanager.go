package main

import (
	"time"

	"github.com/garyburd/redigo/redis"
	. "github.com/gtechx/base/common"
	"github.com/gtechx/chatserver/config"
)

type RedisDataManager struct {
	redisPool *redis.Pool
}

var instanceDataManager *RedisDataManager

func DataManager() *RedisDataManager {
	if instanceDataManager == nil {
		instanceDataManager = &RedisDataManager{}
	}
	return instanceDataManager
}

func (rdm *RedisDataManager) Initialize() error {
	rdm.redisPool = &redis.Pool{
		MaxIdle:      3,
		IdleTimeout:  240 * time.Second,
		Dial:         redisDial,
		TestOnBorrow: redisOnBorrow,
	}

	conn := rdm.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("EXISTS", "UID")

	if err != nil {
		return err
	}

	if !Bool(ret) {
		_, err = conn.Do("SET", "UID", config.StartUID)

		if err != nil {
			return err
		}
	}

	ret, err = conn.Do("EXISTS", "APPID")

	if err != nil {
		return err
	}

	if !Bool(ret) {
		_, err = conn.Do("SET", "APPID", config.StartAPPID)

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

		_, err = conn.Do("HSET", 0, "password", Md5("ztgame@123"))

		if err != nil {
			return err
		}
	}

	return err
}

func redisDial() (redis.Conn, error) {
	c, err := redis.Dial("tcp", config.RedisAddr)
	if err != nil {
		return nil, err
	}
	if _, err := c.Do("AUTH", "ztgame@123"); err != nil {
		c.Close()
		return nil, err
	}
	// if _, err := c.Do("SELECT", db); err != nil {
	// 	c.Close()
	// 	return nil, err
	// }
	return c, nil
}

func redisOnBorrow(c redis.Conn, t time.Time) error {
	if time.Since(t) < time.Minute {
		return nil
	}
	_, err := c.Do("PING")
	return err
}
