package main

import (
	"fmt"

	"github.com/gtechx/base/gtnet"
	"github.com/gtechx/chatserver/entity"
)

var BigMsgIDCounter uint8 = 0
var SmallMsgIDCounter uint8 = 0

const (
	TYPE_NULL int = 1
	TYPE_USER int = 2
)

const (
	BIG_MSG_ID_USER uint8 = iota
	BIG_MSG_ID_ERR

	BIG_MSG_ID_COUNT
)

const (
	SMALL_MSG_ID_ERR_UNKNOWN = iota
	SMALL_MSG_ID_ERR_REDIS
	SMALL_MSG_ID_ERR_CODE
)

const BIG_MSG_ID_START = BIG_MSG_ID_USER

var msgProcesser [][]func(*UserEntity, []byte)

// var bigMsgProcesser []func(IEntity, []byte) = []func(IEntity, []byte){}
// var smallMsgProcesser []func(IEntity, []byte) = []func(IEntity, []byte){}

func init() {
	if msgProcesser == nil {
		msgProcesser = make([][]func(*UserEntity, []byte), BIG_MSG_ID_COUNT)
	}
}

type CEntityManager struct {
	nullEntityMap             map[uint64]gtentity.IEntity
	userIDEntityMap           map[uint64]gtentity.IEntity
	userAPPIDZONEUIDEntityMap map[uint64]map[uint32]map[uint64]gtentity.IEntity
	curID                     uint64

	delChan chan gtentity.IEntity
	addChan chan gtentity.IEntity
}

var instanceEntityManager *CEntityManager

func EntityManager() *CEntityManager {
	if instanceEntityManager == nil {
		instanceEntityManager = &CEntityManager{nullEntityMap: make(map[uint64]gtentity.IEntity), userIDEntityMap: make(map[uint64]gtentity.IEntity), userAPPIDZONEUIDEntityMap: make(map[uint64]map[uint32]map[uint64]gtentity.IEntity)}
	}
	return instanceEntityManager
}

func (this *CEntityManager) Initialize() {
	this.delChan = make(chan gtentity.IEntity, 1024)
	this.addChan = make(chan gtentity.IEntity, 1024)

	go this.userEntityProcess()
}

func (this *CEntityManager) CreateNullEntity(conn gtnet.IConn) gtentity.IEntity {
	this.curID++
	entity := newNullEntity(this.curID, conn)
	entity.start()
	return entity
}

// func (this *CEntityManager) RemoveNullEntity(id uint64) {
// 	delete(this.nullEntityMap, id)
// }

func (this *CEntityManager) CreateEntity(etype int, entity gtentity.IEntity) gtentity.IEntity {
	switch etype {
	case TYPE_USER:
		newentity := newUserEntity(entity)
		this.addChan <- newentity
		return newentity
	}

	return nil
}

func (this *CEntityManager) RemoveEntity(entity gtentity.IEntity) {
	this.delChan <- entity
}

func (this *CEntityManager) doAddEntity(entity gtentity.IEntity) {
	eid := entity.ID()
	uid := entity.UID()
	zone := entity.ZONE()
	appid := entity.APPID()

	oldappmap, ok := this.userAPPIDZONEUIDEntityMap[appid]

	if !ok {
		this.userAPPIDZONEUIDEntityMap[appid] = make(map[uint32]map[uint64]gtentity.IEntity)
		this.userAPPIDZONEUIDEntityMap[appid][zone] = make(map[uint64]gtentity.IEntity)
	} else {
		oldzonemap, ok := oldappmap[zone]

		if !ok {
			this.userAPPIDZONEUIDEntityMap[appid][zone] = make(map[uint64]gtentity.IEntity)
		} else {
			oldentity, ok := oldzonemap[uid]

			if ok {
				oldeid := oldentity.ID()
				oldentity.ForceOffline()
				delete(oldzonemap, uid)
				delete(this.userIDEntityMap, oldeid)
			}
		}
	}

	this.userIDEntityMap[eid] = entity
	this.userAPPIDZONEUIDEntityMap[appid][zone][uid] = entity
}

func (this *CEntityManager) userEntityProcess() {
	select {
	case entity := <-this.addChan:
		if entity.Conn() != nil {
			this.doAddEntity(entity)
		}
	case entity := <-this.delChan:
		if entity.Conn() != nil {
			this.doRemoveEntity(entity)
		}
	}
}

func (this *CEntityManager) doRemoveEntity(entity gtentity.IEntity) {
	eid := entity.ID()
	uid := entity.UID()
	zone := entity.ZONE()
	appid := entity.APPID()

	entity, ok := this.userIDEntityMap[eid]

	if ok {
		delete(this.userIDEntityMap, eid)
		delete(this.userAPPIDZONEUIDEntityMap[appid][zone], uid)
	}
}

func (this *CEntityManager) CleanOnlineUsers() {
	for _, entity := range this.userIDEntityMap {
		DataManager().SetUserOffline(entity.(*UserEntity))
	}

	fmt.Println("cleanOnlineUsers end")
}
