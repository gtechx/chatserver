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

func (sm *SessManager) CreateSess(conn net.Conn, appname, zonename, account string, id uint64) *Sess {
	sess := &Sess{account, appname, zonename, id, conn}
	sm.sessMap.Store(id, sess)
	return sess
}

func (sm *SessManager) DelSess(id uint64) {
	sm.sessMap.Delete(id)
}

func (sm *SessManager) GetSess(id uint64) *Sess {
	sess, ok := sm.sessMap.Load(id)
	if ok {
		return sess.(*Sess)
	}
	return nil
}
