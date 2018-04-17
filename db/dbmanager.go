package gtdb

import (
	"time"

	"github.com/garyburd/redigo/redis"
	//. "github.com/gtechx/base/common"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
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

	return nil
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
	prefix         string
}

func (mdm *Mysql) Initialize(saddr, user_pass, defaultdb, prefix string) error {
	mdm.serverAddr = saddr
	mdm.serverPassword = user_pass
	mdm.defaultDB = defaultdb
	mdm.prefix = prefix

	db, err := gorm.Open("mysql", user_pass+"@tcp("+saddr+")/"+defaultdb+"?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		return err
	}

	gorm.DefaultTableNameHandler = mdm.DefaultTableNameHandler

	db.DB().SetMaxIdleConns(10)
	db.LogMode(true)

	mdm.DB = db
	return err
}

func (mdm *Mysql) DefaultTableNameHandler(db *gorm.DB, defaultTableName string) string {
	return mdm.prefix + "_" + defaultTableName
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
	return db.rd.Initialize(saddr, spass, defaultdb)
}

func (db *DBManager) InitializeMysql(saddr, user_pass, defaultdb, prefix string) error {
	db.sql = &Mysql{}
	return db.sql.Initialize(saddr, user_pass, defaultdb, prefix)
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

func (db *DBManager) Install() error {
	conn := db.rd.Get()
	defer conn.Close()

	_, err := conn.Do("FLUSHDB")

	if err != nil {
		return err
	}

	tx := db.sql.Begin()
	for _, dbtable := range db_tables {
		if err = tx.DropTableIfExists(dbtable).Error; err != nil {
			tx.Rollback()
			return err
		}
		if err = tx.CreateTable(dbtable).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	//create test data
	tbl_account := &Account{Account: "wyq", Password: "edf06a849c9ec19ea725bd3c6c4ce225", Salt: "p99U86", Regip: "127.0.0.1"}
	if err = tx.Create(tbl_account).Error; err != nil {
		tx.Rollback()
		return err
	}

	tbl_app := &App{Name: "test1", Owner: "wyq", Desc: "ddddd", Share: ""}
	if err = tx.Create(tbl_app).Error; err != nil {
		tx.Rollback()
		return err
	}

	tbl_zone := &AppZone{Name: "aaa", Owner: "test1"}
	if err = tx.Create(tbl_zone).Error; err != nil {
		tx.Rollback()
		return err
	}
	//end

	tx.Commit()
	return err
}
