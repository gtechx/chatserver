package gtdata

import (
	"time"

	"github.com/garyburd/redigo/redis"
	. "github.com/gtechx/base/common"
)

//[set]app aid set
//[hset]app:aid:uid aid owner desc regdate
//[hset]app:aid:uid:config

//app op
func (rdm *RedisDataManager) CreateApp(uid uint64, name string) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("INCR", "UID")

	if err != nil {
		return err
	}

	appid := Uint64(ret)

	conn.Send("MULTI")
	conn.Send("SADD", "app", appid)
	conn.Send("SADD", "app:"+String(uid), appid)
	conn.Send("HMSET", "app:"+String(appid), "name", name, "owner", uid, "desc", "", "iconurl", "", "regdate", time.Now().Unix())

	_, err = conn.Do("EXEC")

	return err
}

func (rdm *RedisDataManager) DeleteApp(uid, appid uint64) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()

	conn.Send("MULTI")
	conn.Send("SREM", "app", appid)
	conn.Send("SREM", "app:"+String(uid), appid)
	conn.Send("DEL", "app:"+String(appid))

	_, err := conn.Do("EXEC")

	return err
}

func (rdm *RedisDataManager) SetAppType(appid uint64, typestr string) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("HSET", "app:"+String(appid), "type", typestr)
	return err
}

func (rdm *RedisDataManager) GetAppType(appid uint64) (string, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("HGET", "app:"+String(appid), "type")
	return redis.String(ret, err)
}

func (rdm *RedisDataManager) IsAppExists(appid uint64) (bool, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("SISMEMBER", "app", appid)
	return redis.Bool(ret, err)
}

func (rdm *RedisDataManager) AddAppZone(appid uint64, zone uint32, zonename string) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("HSET", "app:zone:"+String(appid), zone, zonename)
	return err
}

func (rdm *RedisDataManager) GetAppZones(appid uint64) (map[uint32]string, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("HGETALL", "app:zone:"+String(appid))

	if err != nil {
		return nil, err
	}

	retarr, err := redis.Values(ret, err)

	if err != nil {
		return nil, err
	}

	zonemap := map[uint32]string{}
	for i := 0; i < len(retarr); {
		zoneid := Uint32(retarr[i])
		zonename := String(retarr[i+1])
		zonemap[zoneid] = zonename
		i = i + 2
	}

	return zonemap, err
}

func (rdm *RedisDataManager) GetAppOwner(appid uint64) (uint64, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("HGET", "app:"+String(appid), "owner")
	return redis.Uint64(ret, err)
}

func (rdm *RedisDataManager) IsAppZone(appid uint64, zone uint32) (bool, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("HEXISTS", "app:zone:"+String(appid), zone)
	return redis.Bool(ret, err)
}

func (rdm *RedisDataManager) GetAppZoneName(appid uint64, zone uint32) (string, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("HGET", "app:"+String(appid), zone)
	return redis.String(ret, err)
}

func (rdm *RedisDataManager) AddShareApp(uid, appid, otherappid uint64) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	conn.Send("MULTI")
	conn.Send("HMSET", "app:"+String(uid)+":"+String(appid), "share", otherappid)
	conn.Send("SADD", "app:share:"+String(otherappid), appid)
	_, err := conn.Do("EXEC")
	return err
}

func (rdm *RedisDataManager) IsShareWithOtherApp(uid, appid uint64) (bool, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("HEXISTS", "app:"+String(uid)+":"+String(appid), "share")
	return redis.Bool(ret, err)
}

func (rdm *RedisDataManager) GetShareApp(uid, appid uint64) (uint64, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("HGET", "app:"+String(uid)+":"+String(appid), "share")
	return redis.Uint64(ret, err)
}

func (rdm *RedisDataManager) GetMyShareAppList(appid uint64) ([]uint64, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("SMEMBERS", "app:share:"+String(appid))

	if err != nil {
		return nil, err
	}

	retarr, err := redis.Values(ret, nil)

	if err != nil {
		return nil, err
	}

	applist := []uint64{}
	for _, otherappid := range retarr {
		applist = append(applist, Uint64(otherappid))
	}

	return applist, err
}
