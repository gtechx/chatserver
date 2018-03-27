package main

import (
	"github.com/garyburd/redigo/redis"
	. "github.com/nature19862001/base/common"
)

func (rdm *RedisDataManager) PullOnlineMessage(serveraddr string, timeout int) ([]byte, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("BLPOP", "message:"+serveraddr, timeout)

	if err != nil {
		return nil, err
	}

	retarr, err := redis.Values(ret, nil)

	if err != nil {
		return nil, err
	}

	return Bytes(retarr[1]), err
}

func (rdm *RedisDataManager) GetOfflineMessage(entity *UserEntity) ([][]byte, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("LRANGE", entity.KeyMessageOffline, 0, -1)

	if err != nil {
		return nil, err
	}

	retarr, err := redis.Values(ret, nil)

	if err != nil {
		return nil, err
	}

	msglist := [][]byte{}
	for i := 1; i < len(retarr); i++ {
		msglist = append(msglist, Bytes(retarr[i]))
	}

	return msglist, err
}

func (rdm *RedisDataManager) SendMsgToUserOnline(uid, appid uint64, data []byte, serveraddr string) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("RPUSH", "message:"+serveraddr, data)
	return err
}

func (rdm *RedisDataManager) SendMsgToUserOffline(entity *UserEntity, data []byte) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("RPUSH", entity.KeyMessageOffline, data)
	return err
}

func (rdm *RedisDataManager) SendMsgToRoom() {

}
