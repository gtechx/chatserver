package main

import (
	"net"
	"sync"
)

type SessManager struct {
	sessMap *sync.Map
}

var sessmgr *SessManager

func SessMgr() *SessManager {
	if sessmgr == nil {
		sessmgr = &SessManager{sessMap: &sync.Map{}}
	}
	return sessmgr
}

func (sm *SessManager) CreateSess(conn net.Conn, appname, zonename, account string, id uint64) ISession {
	sess := &Sess{account: account, appname: appname, zonename: zonename, id: id, conn: conn}
	sm.sessMap.Store(id, sess)
	return sess
}

func (sm *SessManager) DelSess(id uint64) {
	sm.sessMap.Delete(id)
}

func (sm *SessManager) GetSess(id uint64) ISession {
	sess, ok := sm.sessMap.Load(id)
	if ok {
		return sess.(*Sess)
	}
	return nil
}

func (sm *SessManager) SendMsgToId(id uint64, msg []byte) bool {
	sess, ok := sm.sessMap.Load(id)
	if ok {
		sess.(*Sess).Send(msg)
		return true
	}
	return false
}
