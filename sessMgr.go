package main

import (
	"net"
	"sync"
)

type SessMgr struct {
	sessMap *sync.Map
}

func (sm *SessMgr) CreateSess(conn net.Conn, appname, zonename string, id uint64) *Sess {
	sess := &Sess{appname, zonename, id, conn}
	sm.sessMap.Store(id, sess)
	return sess
}

func (sm *SessMgr) DelSess(id uint64) {
	sm.sessMap.Delete(id)
}

func (sm *SessMgr) GetSess(id uint64) *Sess {
	sess, ok := sm.sessMap.Load(id)
	if ok {
		return sess.(*Sess)
	}
	return nil
}
