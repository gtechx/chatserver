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
func (rdm *RedisDataManager) CreateApp(account, appname string) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("INCR", "APPID")

	if err != nil {
		return err
	}

	appid := Uint64(ret)
	regdate := time.Now().Unix()

	conn.Send("MULTI")
	conn.Send("SADD", "set:app", appname)
	conn.Send("SADD", "set:app:account:"+account, appname) //添加uid防止app:appid和app:uid重复
	conn.Send("HMSET", "hset:app:appname:"+appname, "appid", appid, "appname", appname, "owner", account, "desc", "", "iconurl", "", "regdate", regdate)
	conn.Send("ZADD", "zset:app:regdate:account:"+account, regdate, appname) //create index of app regdate

	_, err = conn.Do("EXEC")

	return err
}

func (rdm *RedisDataManager) DeleteApp(account, appname string) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()

	conn.Send("MULTI")
	conn.Send("SREM", "set:app", appname)
	conn.Send("SREM", "set:app:account:"+account, appname)
	conn.Send("DEL", "hset:app:appname:"+appname)
	conn.Send("DEL", "set:app:share:"+appname)
	conn.Send("DEL", "set:app:zone:"+appname)
	conn.Send("ZREM", "zset:app:regdate:account:"+account, appname)

	_, err := conn.Do("EXEC")

	return err
}

func (rdm *RedisDataManager) IsAppExists(appname string) (bool, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("SISMEMBER", "set:app", appname)
	return redis.Bool(ret, err)
}

func (rdm *RedisDataManager) GetApp(datakey *DataKey) (*App, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("HGETALL", datakey.KeyAppHsetByAppname)

	retarr, err := redis.Values(ret, err)

	if err != nil {
		return nil, err
	}

	app := new(App)
	err = redis.ScanStruct(retarr, app)

	if err != nil {
		return nil, err
	}

	tm := time.Unix(app.Regdate, 0)
	//var stamp = fmt.Sprintf("\"%s\"", time.Time(this).Format("2006-01-02 15:04:05"))
	app.Sregdate = tm.Format("2006-01-02 15:04:05")

	return app, err
}

func (rdm *RedisDataManager) SetAppField(datakey *DataKey, fieldname string, value interface{}) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("HSET", datakey.KeyAppHsetByAppname, fieldname, value)
	return err
}

func (rdm *RedisDataManager) GetAppField(datakey *DataKey, fieldname string, value interface{}) (interface{}, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("HGET", datakey.KeyAppHsetByAppname, fieldname)
	return ret, err
}

func (rdm *RedisDataManager) GetAppCount(datakey *DataKey) (uint64, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("ZCARD", datakey.KeyAppZsetRegdateAppnameByAccount)
	return redis.Uint64(ret, err)
}

func (rdm *RedisDataManager) GetAppnameByPage(datakey *DataKey, start, end int) ([]string, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("ZRANGE", datakey.KeyAppZsetRegdateAppnameByAccount, start, end)

	retarr, err := redis.Values(ret, err)

	if err != nil {
		return nil, err
	}

	applist := []string{}
	for i := 0; i < len(retarr); i++ {
		appname, _ := redis.String(retarr[i], err)
		applist = append(applist, appname)
	}

	return applist, err
}

// func (rdm *RedisDataManager) SetAppType(appid uint64, typestr string) error {
// 	conn := rdm.redisPool.Get()
// 	defer conn.Close()
// 	_, err := conn.Do("HSET", "app:"+String(appid), "type", typestr)
// 	return err
// }

// func (rdm *RedisDataManager) GetAppType(appid uint64) (string, error) {
// 	conn := rdm.redisPool.Get()
// 	defer conn.Close()
// 	ret, err := conn.Do("HGET", "app:"+String(appid), "type")
// 	return redis.String(ret, err)
// }

// func (rdm *RedisDataManager) IsAppIDExists(datakey *DataKey) (bool, error) {
// 	conn := rdm.redisPool.Get()
// 	defer conn.Close()
// 	ret, err := conn.Do("SISMEMBER", datakey.KeyAppSet, datakey.Appid)
// 	return redis.Bool(ret, err)
// }

func (rdm *RedisDataManager) AddAppZone(datakey *DataKey) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("SADD", datakey.KeyAppSetZonenameByAppname, datakey.Zonename)
	return err
}

func (rdm *RedisDataManager) DelAppZone(datakey *DataKey) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("SREM", datakey.KeyAppSetZonenameByAppname, datakey.Zonename)
	return err
}

func (rdm *RedisDataManager) GetAppZones(datakey *DataKey) ([]string, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("SMEMBERS", datakey.KeyAppSetZonenameByAppname)

	retarr, err := redis.Values(ret, err)

	if err != nil {
		return nil, err
	}

	zonelist := []string{}
	for i := 0; i < len(retarr); i++ {
		zonename := String(retarr[i])
		zonelist = append(zonelist, zonename)
	}

	return zonelist, err
}

func (rdm *RedisDataManager) GetAppOwner(datakey *DataKey) (uint64, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("HGET", datakey.KeyAppHsetByAppname, "owner")
	return redis.Uint64(ret, err)
}

func (rdm *RedisDataManager) IsAppZone(datakey *DataKey, zonename string) (bool, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("SISMEMBER", datakey.KeyAppSetZonenameByAppname, zonename)
	return redis.Bool(ret, err)
}

// func (rdm *RedisDataManager) GetAppZoneName(datakey *DataKey) (string, error) {
// 	conn := rdm.redisPool.Get()
// 	defer conn.Close()
// 	ret, err := conn.Do("HGET", "app:"+String(appid), zone)
// 	return redis.String(ret, err)
// }

func (rdm *RedisDataManager) AddShareApp(datakey *DataKey, otherappname string) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	conn.Send("MULTI")
	conn.Send("HSET", datakey.KeyAppHsetByAppname, "share", otherappname)
	conn.Send("SADD", "set:app:share:"+otherappname, otherappname)
	_, err := conn.Do("EXEC")
	return err
}

func (rdm *RedisDataManager) DelShareApp(datakey *DataKey, otherappname string) error {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	conn.Send("MULTI")
	conn.Send("HDEL", datakey.KeyAppHsetByAppname, "share", otherappname)
	conn.Send("SREM", "set:app:share:"+otherappname, otherappname)
	conn.Send("DEL", datakey.KeyAppSetZonenameByAppname)
	_, err := conn.Do("EXEC")
	return err
}

func (rdm *RedisDataManager) IsShareWithOtherApp(datakey *DataKey) (bool, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("HEXISTS", datakey.KeyAppHsetByAppname, "share")
	return redis.Bool(ret, err)
}

func (rdm *RedisDataManager) GetShareApp(datakey *DataKey) (string, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()
	ret, err := conn.Do("HGET", datakey.KeyAppHsetByAppname, "share")
	return redis.String(ret, err)
}

func (rdm *RedisDataManager) GetShareAppList(datakey *DataKey) ([]string, error) {
	conn := rdm.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("SMEMBERS", datakey.KeyAppSetShareByAppname)

	retarr, err := redis.Values(ret, nil)

	if err != nil {
		return nil, err
	}

	applist := []string{}
	for _, otherapp := range retarr {
		applist = append(applist, String(otherapp))
	}

	return applist, err
}
