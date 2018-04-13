package gtdb

import (
	"time"

	"github.com/garyburd/redigo/redis"
	. "github.com/gtechx/base/common"
	"github.com/jinzhu/gorm"
)

type Redis struct {
	*redis.Pool

	serverAddr     string
	serverPassword string
	defaultDB      uint64
}

func (rdm *Redis) Initialize(saddr, spass string, defaultdb uint64) error {
	rdm.serverAddr = saddr
	rdm.serverPassword = spass
	rdm.defaultDB = defaultdb

	rdm.Pool = &redis.Pool{
		MaxIdle:      3,
		IdleTimeout:  240 * time.Second,
		Dial:         rdm.redisDial,
		TestOnBorrow: rdm.redisOnBorrow,
	}

	return err
}

func (rdm *Redis) UnInitialize() error {
	var err error
	if rdm.Pool != nil {
		err = rdm.Pool.Close()
	}
	return err
}

func (rdm *Redis) redisDial() (redis.Conn, error) {
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

func (rdm *Redis) redisOnBorrow(c redis.Conn, t time.Time) error {
	if time.Since(t) < time.Minute {
		return nil
	}
	_, err := c.Do("PING")
	return err
}

type Mysql struct {
	*gorm.DB

	serverAddr     string
	serverPassword string
	defaultDB      string
}

func (mdm *Mysql) Initialize(saddr, user, spass, defaultdb string) error {
	mdm.serverAddr = saddr
	mdm.serverPassword = spass
	mdm.defaultDB = defaultdb

	db, err := gorm.Open("mysql", user+":"+spass+"@tcp("+saddr+")/"+defaultdb+"?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		return err
	}

	db.DB().SetMaxIdleConns(10)
	db.LogMode(true)

	mdm.DB = db
	return err
}

func (mdm *Mysql) UnInitialize() error {
	var err error
	if mdm.DB != nil {
		err = mdm.DB.Close()
	}
	return err
}

type DBManager struct {
	rd  *Redis
	sql *Mysql
}

var instance *DBManager

func Manager() *DBManager {
	if instance == nil {
		instance = &DBManager{}
	}
	return instance
}

func (db *DBManager) InitializeRedis(saddr, spass string, defaultdb uint64) error {
	db.rd = &Redis{}
	db.rd.Initialize(saddr, spass, defaultdb)
}

func (db *DBManager) InitializeMysql(saddr, spass, defaultdb string) error {
	db.sql = &Mysql{}
	db.sql.Initialize(saddr, spass, defaultdb)
}

func (db *DBManager) UnInitialize() error {
	var err error
	if db.rd != nil {
		err = db.rd.UnInitialize()
		db.rd = nil
	}
	if db.sql != nil {
		err = db.sql.UnInitialize()
		db.sql = nil
	}
	return err
}
