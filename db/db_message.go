package gtdb

import (
	"github.com/garyburd/redigo/redis"
	. "github.com/gtechx/base/common"
)

//服务器记录所有消息？
//服务器记录7天之内的消息？如果超出，则删除更早的。msg:id:timestamp 设置定时器
//未读消息数量怎样存储

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

func (db *DBManager) AddPresence(from, to uint64, msg []byte) error {
	conn := db.rd.Get()
	defer conn.Close()
	_, err := conn.Do("HSET", "presence:"+String(to), from, msg) //记录到目的地用户presence列表
	return err
}

func (db *DBManager) RemovePresence(id, from uint64) error {
	conn := db.rd.Get()
	defer conn.Close()
	_, err := conn.Do("HDEL", "presence:"+String(id), from)
	return err
}

func (db *DBManager) GetAllPresence(id uint64) (map[string]string, error) {
	conn := db.rd.Get()
	defer conn.Close()
	ret, err := conn.Do("HGETALL", "presence:"+String(id))
	return redis.StringMap(ret, err) //.ByteSlices(ret, err)
}

func (db *DBManager) AddRoomPresence(rid, appdataid uint64, msg []byte) error {
	conn := db.rd.Get()
	defer conn.Close()
	_, err := conn.Do("HSET", "roompresence:"+String(rid), appdataid, msg)
	//_, err := conn.Do("SET", "roompresence:"+String(rid)+":"+String(appdataid), "", "ex 259200") //记录到目的地用户presence列表
	return err
}

func (db *DBManager) RemoveRoomPresence(rid, appdataid uint64) error {
	conn := db.rd.Get()
	defer conn.Close()
	_, err := conn.Do("HDEL", "roompresence:"+String(rid), appdataid)
	//_, err := conn.Do("DEL", "roompresence:"+String(rid)+":"+String(appdataid))
	return err
}

func (db *DBManager) GetAllRoomPresence(rid uint64) (map[string]string, error) {
	conn := db.rd.Get()
	defer conn.Close()
	ret, err := conn.Do("HGETALL", "roompresence:"+String(rid))
	return redis.StringMap(ret, err) //.ByteSlices(ret, err)
}

func (db *DBManager) IsRoomPresenceExists(rid, appdataid uint64) (bool, error) {
	conn := db.rd.Get()
	defer conn.Close()
	ret, err := conn.Do("HEXISTS", "roompresence:"+String(rid), appdataid)
	//ret, err := conn.Do("EXISTS", "roompresence:"+String(rid)+":"+String(appdataid))
	return redis.Bool(ret, err)
}

func (db *DBManager) PullOnlineMessage(serveraddr string) ([]byte, error) {
	conn := db.rd.Get()
	defer conn.Close()
	ret, err := conn.Do("LPOP", "message:"+serveraddr)
	return redis.Bytes(ret, err)
}

func (db *DBManager) GetOfflineMessage(id uint64) ([][]byte, error) {
	conn := db.rd.Get()
	defer conn.Close()

	ret, err := conn.Do("LRANGE", "message:offline:"+String(id), 0, -1)
	datalist, err := redis.ByteSlices(ret, err)
	conn.Do("LTRIM", "message:offline:"+String(id), len(datalist), -1)

	return datalist, err
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

func (db *DBManager) SendMsgToUserOnline(msg []byte, serveraddr string) error {
	conn := db.rd.Get()
	defer conn.Close()
	_, err := conn.Do("RPUSH", "message:"+serveraddr, msg)
	return err
}

func (db *DBManager) SendMsgToUserOffline(to uint64, data []byte) error {
	conn := db.rd.Get()
	defer conn.Close()
	_, err := conn.Do("RPUSH", "message:offline:"+String(to), data)
	return err
}
