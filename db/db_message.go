package gtdb

import (
	"github.com/garyburd/redigo/redis"
	. "github.com/gtechx/base/common"
)

// func (db *DBManager) IsPresenceRecordExists(from, to uint64) (bool, error) {
// 	conn := db.rd.Get()
// 	defer conn.Close()
// 	ret, err := conn.Do("HEXISTS", "presence:record:"+String(from), to)
// 	return redis.Bool(ret, err)
// }

// func (db *DBManager) AddPresenceRecord(from, to uint64, data []byte) error {
// 	conn := db.rd.Get()
// 	defer conn.Close()
// 	_, err := conn.Do("HSET", "presence:record:"+String(from), to, data) //记录到发送者用户记录列表，用于校验
// 	_, err := conn.Do("HSET", "presence:"+String(to), from, data)        //记录到目的地用户presence列表
// 	return err
// }

// func (db *DBManager) RemovePresenceRecord(from, to uint64) error {
// 	conn := db.rd.Get()
// 	defer conn.Close()
// 	_, err := conn.Do("HDEL", "presence:record:"+String(from), to)
// 	return err
// }

func (db *DBManager) IsPresenceExists(id, from uint64) (bool, error) {
	conn := db.rd.Get()
	defer conn.Close()
	ret, err := conn.Do("HEXISTS", "presence:"+String(id), from)
	return redis.Bool(ret, err)
}

func (db *DBManager) AddPresence(from, to uint64, data []byte) error {
	conn := db.rd.Get()
	defer conn.Close()
	_, err := conn.Do("HSET", "presence:"+String(to), from, data) //记录到目的地用户presence列表
	return err
}

func (db *DBManager) RemovePresence(id, from uint64) error {
	conn := db.rd.Get()
	defer conn.Close()
	_, err := conn.Do("HDEL", "presence:"+String(id), from)
	return err
}

func (db *DBManager) GetAllPresence(id uint64) ([][]byte, error) {
	conn := db.rd.Get()
	defer conn.Close()
	ret, err := conn.Do("HGETALL", "presence:"+String(id))
	return redis.ByteSlices(ret, err)
}

func (db *DBManager) PullOnlineMessage(serveraddr string, timeout int) ([]byte, error) {
	conn := db.rd.Get()
	defer conn.Close()
	ret, err := conn.Do("BLPOP", "message:"+serveraddr, timeout)
	return redis.Bytes(ret, err)
}

func (db *DBManager) GetOfflineMessage(id uint64) ([][]byte, error) {
	conn := db.rd.Get()
	defer conn.Close()

	ret, err := conn.Do("LRANGE", "message:offline:"+String(id), 0, -1)

	return redis.ByteSlices(ret, err)
	// if err != nil {
	// 	return nil, err
	// }

	// retarr, err := redis.Values(ret, nil)

	// if err != nil {
	// 	return nil, err
	// }

	// msglist := [][]byte{}
	// for i := 1; i < len(retarr); i++ {
	// 	msglist = append(msglist, Bytes(retarr[i]))
	// }

	// return msglist, err
}

func (db *DBManager) SendMsgToUserOnline(data []byte, serveraddr string) error {
	conn := db.rd.Get()
	defer conn.Close()
	_, err := conn.Do("RPUSH", "message:"+serveraddr, data)
	return err
}

func (db *DBManager) SendMsgToUserOffline(to uint64, data []byte) error {
	conn := db.rd.Get()
	defer conn.Close()
	_, err := conn.Do("RPUSH", "message:offline:"+String(to), data)
	return err
}
