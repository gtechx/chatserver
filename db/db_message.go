package gtdb

import (
	"github.com/garyburd/redigo/redis"
	. "github.com/gtechx/base/common"
)

func (db *DBManager) PullOnlineMessage(serveraddr string, timeout int) ([]byte, error) {
	conn := db.redisPool.Get()
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

func (db *DBManager) GetOfflineMessage(datakey *DataKey) ([][]byte, error) {
	conn := db.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("LRANGE", datakey.KeyAppDataListMsgByAppidZonenameAccount, 0, -1)

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

func (db *DBManager) SendMsgToUserOnline(uid, appid uint64, data []byte, serveraddr string) error {
	conn := db.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("RPUSH", "message:"+serveraddr, data)
	return err
}

func (db *DBManager) SendMsgToUserOffline(datakey *DataKey, data []byte) error {
	conn := db.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("RPUSH", datakey.KeyAppDataListMsgByAppidZonenameAccount, data)
	return err
}

func (db *DBManager) SendMsgToRoom() {

}
